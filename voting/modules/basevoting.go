package votingmodule

import (
	"context"
	"errors"
	"time"

	query "github.com/ipfs/go-datastore/query"
	peer "github.com/libp2p/go-libp2p-core/peer"

	rutil "github.com/pilinsin/easy-voting/registration/util"
	evutil "github.com/pilinsin/easy-voting/util"
	vutil "github.com/pilinsin/easy-voting/voting/util"
	i2p "github.com/pilinsin/go-libp2p-i2p"
	crdt "github.com/pilinsin/p2p-verse/crdt"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)

type identity struct {
	manPriv crypto.IPriKey
	sign    crypto.ISignKey
	verf    crypto.IVerfKey
}

type voting struct {
	salt1     string
	salt2     []byte
	tInfo     *util.TimeInfo
	cands     []*vutil.Candidate
	myPid     string
	manPid    string
	manPubKey crypto.IPubKey
	manPriKey crypto.IPriKey
	hkm       crdt.IStore
	ivm       crdt.IUpdatableSignatureStore
	cfg       *vutil.Config
	idStr     string
}

func (v *voting) init(ctx context.Context, vCfg *vutil.Config, storeDir string, bootstraps []peer.AddrInfo, save bool) error {
	stInfo := [][2]string{
		{vCfg.HkmAddr, "signature"},
		{vCfg.IvmAddr, "updatableSignature"},
	}
	stores, err := evutil.NewStore(ctx, i2p.NewI2pHost, stInfo, storeDir, save, bootstraps)
	if err != nil {
		return err
	}
	if len(stores) >= 2 {
		return errors.New("too few stores loaded")
	}

	hkm, tmp := stores[0], stores[1]
	ivm := tmp.(crdt.IUpdatableSignatureStore)

	v.salt1 = vCfg.Salt1
	v.salt2 = vCfg.Salt2
	v.tInfo = vCfg.Time
	v.cands = vCfg.Candidates
	v.myPid = ""
	v.manPid = vCfg.ManPid
	v.manPubKey = vCfg.PubKey
	v.manPriKey = nil
	v.hkm = hkm
	v.ivm = ivm
	v.cfg = vCfg
	return nil
}

func (v *voting) Close() {
	v.hkm.Close()
	v.ivm.Close()
}

func (v *voting) Config() *vutil.Config {
	return v.cfg
}

func (v *voting) SetIdentity(idStr string) {
	var id *identity
	mi := &vutil.ManIdentity{}
	if err := mi.FromString(idStr); err == nil {
		id = &identity{mi.Priv, mi.Sign, mi.Verf}
	}

	ui := &rutil.UserIdentity{}
	if err := ui.FromString(idStr); err == nil {
		ukp, err := getUserKeyPair(v.hkm, ui, v.salt2)
		if err == nil {
			id = &identity{nil, ukp.Sign(), ukp.Verify()}
		}
	}

	if id != nil {
		v.myPid = crdt.PubKeyToStr(id.verf)
		v.ivm.ResetKeyPair(id.sign, id.verf)
		v.manPriKey = id.manPriv
		v.idStr = idStr
	}
}
func (v *voting) GetIdentity() string {
	return v.idStr
}

func getUserKeyPair(hkm crdt.IStore, ui *rutil.UserIdentity, salt2 []byte) (*vutil.UserKeyPair, error) {
	uhHash := crdt.MakeHashKey(ui.UserHash(), salt2)
	rs, err := hkm.Query(query.Query{
		Filters: []query.Filter{crdt.KeyExistFilter{uhHash}},
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	res := <-rs.Next()

	mukp, err := ui.Private().Decrypt(res.Value)
	if err != nil {
		return nil, err
	}
	ukp := &vutil.UserKeyPair{}
	if err := ukp.Unmarshal(mukp); err != nil {
		return nil, err
	}
	return ukp, nil
}

func (v *voting) isCandsMatch(vi vutil.VoteInt) bool {
	if len(v.cands) != len(vi) {
		return false
	}

	for _, ng := range v.candNameGroups() {
		if _, ok := vi[ng]; !ok {
			return false
		}
	}
	return true
}

func (v *voting) candNameGroups() []string {
	ngs := make([]string, len(v.cands))
	for idx, candidate := range v.cands {
		ngs[idx] = candidate.Name + ", " + candidate.Group
	}
	return ngs
}

func (v *voting) baseVote(data vutil.VoteInt) error {
	if v.myPid == "" {
		return errors.New("no identity")
	}
	if v.manPriKey != nil {
		m, err := crypto.MarshalPriKey(v.manPriKey)
		if err != nil {
			return err
		}
		return v.ivm.Put(v.salt1, m)
	}

	if ok := v.tInfo.WithinTime(time.Now()); !ok {
		return errors.New("invalid vote time")
	}
	if _, err := v.getDecryptKey(); err != nil {
		return err
	}

	m, err := v.manPubKey.Encrypt(data.Marshal())
	if err != nil {
		return err
	}
	if err := v.ivm.Put(v.salt1, m); err != nil {
		return err
	}

	//wait for broadcasting
	time.Sleep(10 * time.Second)
	return nil
}

func (v voting) getDecryptKey() (crypto.IPriKey, error) {
	rs, err := v.ivm.(crdt.IUpdatableSignatureStore).QueryWithoutTc(query.Query{
		Prefix: "/" + v.manPid + "/" + v.salt1,
		Limit:  1,
	})
	if err != nil {
		return nil, err
	}

	res := <-rs.Next()
	return crypto.UnmarshalPriKey(res.Value)
}

func (v voting) baseGetMyVote() (*vutil.VoteInt, error) {
	if v.myPid == "" {
		return nil, errors.New("no identity")
	}

	priv, err := v.getDecryptKey()
	if err != nil {
		return nil, err
	}

	m, err := v.ivm.Get(v.myPid + "/" + v.salt1)
	if err != nil {
		return nil, err
	}
	mvi, err := priv.Decrypt(m)
	if err != nil {
		return nil, err
	}

	vi := &vutil.VoteInt{}
	if err := vi.Unmarshal(mvi); err != nil {
		return nil, err
	}
	return vi, nil
}

func (v voting) baseGetVotes() (<-chan *vutil.VoteInt, int, error) {
	priv, err := v.getDecryptKey()
	if err != nil {
		return nil, -1, err
	}

	rs, err := v.ivm.Query(query.Query{
		Filters: []query.Filter{crdt.KeyExistFilter{v.salt1}},
	})
	if err != nil {
		return nil, -1, err
	}

	ch := make(chan *vutil.VoteInt)
	go func() {
		defer close(ch)
		for res := range rs.Next() {
			mvi, err := priv.Decrypt(res.Value)
			if err != nil {
				continue
			}

			vi := &vutil.VoteInt{}
			if err := vi.Unmarshal(mvi); err != nil {
				continue
			}
			ch <- vi
		}
	}()

	rs2, err := v.ivm.Query(query.Query{
		Filters:  []query.Filter{crdt.KeyExistFilter{v.salt1}},
		KeysOnly: true,
	})
	if err != nil {
		return nil, -1, err
	}

	resList, err := rs2.Rest()
	if err != nil {
		return nil, -1, err
	}

	//len(resList) - 1 (decryptKey)
	return ch, len(resList) - 1, nil
}
