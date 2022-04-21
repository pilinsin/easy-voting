package votingmodule

import (
	"time"
	"errors"
	"context"

	query "github.com/ipfs/go-datastore/query"

	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
	crdt "github.com/pilinsin/p2p-verse/crdt"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	vutil "github.com/pilinsin/easy-voting/voting/util"
)

type identity struct{
	manPriv crypto.IPriKey
	sign	crypto.ISignKey
	verf	crypto.IVerfKey
}

type voting struct {
	salt1		string
	tInfo      	*util.TimeInfo
	cands      	[]*vutil.Candidate
	myPid		string
	manPid		string
	manPubKey  	crypto.IPubKey
	manPriKey	crypto.IPriKey
	hkm 		crdt.IStore
	ivm			crdt.IStore
}

func (v *voting) init(ctx context.Context, vCfg *vutil.Config, idStr, storeDir, bAddr string, save bool) error {
	bootstraps := pv.AddrInfosFromString(bAddr)
	cv := crdt.NewVerse(i2p.NewI2pHost, storeDir, save, false, bootstraps...)

	hkm, err := cv.LoadStore(ctx, vCfg.HkmAddr, "signature")
	if err != nil{return err}

	id := getIdentity(idStr, vCfg.Salt2, hkm)
	opt := &crdt.StoreOpts{Priv: id.sign, Pub: id.verf}
	ivm, err := cv.LoadStore(ctx, vCfg.IvmAddr, "updatableSignature", opt)
	if err != nil{return err}

	myPid := ""
	if id.verf != nil{
		myPid = crdt.PubKeyToStr(id.verf)
	}

	v.salt1 = vCfg.Salt1
	v.tInfo = vCfg.Time
	v.cands = vCfg.Candidates
	v.myPid = myPid
	v.manPid = vCfg.ManPid
	v.manPubKey = vCfg.PubKey
	v.manPriKey = id.manPriv
	v.hkm = hkm
	v.ivm = ivm
	return nil
}

func getIdentity(idStr string, salt2 []byte, hkm crdt.IStore) *identity{
	mi := &vutil.ManIdentity{}
	if err := mi.FromString(idStr); err == nil{
		return &identity{mi.Priv, mi.Sign, mi.Verf}
	}
	
	ui := &rutil.UserIdentity{}
	if err := ui.FromString(idStr); err == nil{
		ukp, err := getUserKeyPair(hkm, ui, salt2)
		if err == nil{
			return &identity{nil, ukp.Sign(), ukp.Verify()}
		}
	}
	return &identity{}
}

func getUserKeyPair(hkm crdt.IStore, ui *rutil.UserIdentity, salt2 []byte) (*vutil.UserKeyPair, error){
	uhHash := crdt.MakeHashKey(ui.UserHash(), salt2)
	rs, err := hkm.Query(query.Query{
		Filters: []query.Filter{crdt.KeyExistFilter{uhHash}},
		Limit: 1,
	})
	if err != nil{return nil, err}
	res := <-rs.Next()

	mukp, err := ui.Private().Decrypt(res.Value)
	if err != nil{return nil, err}
	ukp := &vutil.UserKeyPair{}
	if err := ukp.Unmarshal(mukp); err != nil{return nil, err}
	return ukp, nil
}

func (v *voting) Close() {
	v.hkm.Close()
	v.ivm.Close()
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
	if v.myPid == ""{return errors.New("no identity")}
	if v.manPriKey != nil{
		m, err := crypto.MarshalPriKey(v.manPriKey)
		if err != nil{return err}
		return v.ivm.Put(v.salt1, m)
	}

	if ok := v.tInfo.WithinTime(time.Now()); !ok{
		return errors.New("invalid vote time")
	}
	m, err := v.manPubKey.Encrypt(data.Marshal())
	if err != nil{return err}
	return v.ivm.Put(v.salt1, m)
}


func (v voting) getDecryptKey() (crypto.IPriKey, error){
	rs, err := v.ivm.(crdt.IUpdatableSignatureStore).QueryWithoutTc(query.Query{
		Prefix: "/"+v.manPid+"/"+v.salt1,
		Limit: 1,
	})
	if err != nil{return nil, err}
	
	res := <-rs.Next()
	return crypto.UnmarshalPriKey(res.Value)
}

func (v voting) baseGetMyVote() (*vutil.VoteInt, error) {
	if v.myPid == ""{return nil, errors.New("no identity")}

	priv, err := v.getDecryptKey()
	if err != nil{return nil, err}

	m, err := v.ivm.Get(v.myPid+"/"+v.salt1)
	if err != nil{return nil, err}
	mvi, err := priv.Decrypt(m)
	if err != nil{return nil, err}

	vi := &vutil.VoteInt{}
	if err := vi.Unmarshal(mvi); err != nil{return nil, err}
	return vi, nil
}

func (v voting) baseGetVotes() (<-chan *vutil.VoteInt, int, error) {
	priv, err := v.getDecryptKey()
	if err != nil{return nil, -1, err}

	rs, err := v.ivm.Query(query.Query{
		Filters: []query.Filter{crdt.KeyExistFilter{v.salt1}},
	})
	if err != nil{return nil, -1, err}

	ch := make(chan *vutil.VoteInt)
	go func(){
		defer close(ch)
		for res := range rs.Next(){
			mvi, err := priv.Decrypt(res.Value)
			if err != nil{continue}

			vi := &vutil.VoteInt{}
			if err := vi.Unmarshal(mvi); err != nil{continue}
			ch <- vi
		}
	}()

	rs2, err := v.ivm.Query(query.Query{
		Filters: []query.Filter{crdt.KeyExistFilter{v.salt1}},
		KeysOnly: true,
	})
	if err != nil{return nil, -1, err}

	resList, err := rs2.Rest()
	if err != nil{return nil, -1, err}

	//len(resList) - 1 (decryptKey)
	return ch, len(resList)-1, nil
}

