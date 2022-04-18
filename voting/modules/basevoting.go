package votingmodule

import (
	"fmt"
	"time"

	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
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
	cands      	[]vutil.Candidate
	manPid		string
	manPubKey  	crypto.IPubKey
	manPriKey	crypto.IPriKey
	hkm 		crdt.IStore
	ivm			crdt.IStore
}

func (v *voting) init(vCfg *vutil.Config, idStr, bAddr string) {
	storeDir, save := parseIdStr(idStr)
	bootstraps := pv.AddrInfosFromString(bAddr)
	cv := crdt.NewVerse(i2p.NewI2pHost, storeDir, save, false, bootstraps...)
	hkm, err := cv.LoadStore(vCfg.HkmAddr, "signature")
	if err != nil{return}

	id := getIdentity(idStr, vCfg, hkm)
	opt := &crdt.StoreOpts{Priv: id.sign, Pub: id.verf}
	ivm, err := cv.LoadStore(vCfg.IvmAddr, "updatableSignature", opt)
	if err != nil{return}

	v.salt1 = vCfg.Salt1
	v.tInfo = vCfg.Time
	v.cands = vCfg.Candidates
	v.myPid = crdt.PubKeyToStr(id.verf)
	v.manPid = vCfg.ManPid
	v.manPubKey = vCfg.PubKey
	v.manPriKey = id.manPriv
	v.hkm = hkm
	v.ivm = ivm
}

func parseIdStr(idStr string) (string, bool){
	mi := &vutil.ManIdentity{}
	if err := mi.FromString(idStr); err == nil{
		return mi.StoreDir, true
	}
	
	return pv.RandString(8), false
}

func getIdentity(idStr string, vCfg *vutil.Config, hkm crdt.IStore) *identity{
	mi := &vutil.ManIdentity{}
	if err := mi.FromString(idStr); err == nil{
		return &identity{mi.Priv, mi.Sign, mi.Verf}
	}
	
	ui := &rutil.UserIdentity{}
	if err := ui.FromString(idStr); err == nil{
		ukp, err := getUserKeyPair(hkm, ui, vCfg.Salt2)
		if err == nil{
			return &identity{nil, ukp.Sign(), ukp.Verify()}
		}
	}
	return nil
}

func getUserKeyPair(hkm crdt.IStore, ui *rutil.UserIdentity, salt2 []byte) (*vutil.UserKeyPair, error){
	uhHash := crdt.MakeHashKey(ui.UserHash(), salt2)
	rs, err := hkm.Query(query.Query{
		Filters: []query.Filter{crdt.KeyExistFilter{uhHash}},
		Limit: 1,
	})
	if err != nil{return nil, err}
	m := <-rs.Next()

	mukp, err := ui.Private().Decrypt(m)
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
	if ok := v.tInfo.WithinTime(time.Now()); !ok{
		return errors.New("invalid vote time")
	}

	if v.manPriKey != nil{
		m, err := crypto.MarshalPriKey(v.manPriKey)
		if err != nil{return err}
		return v.ivm.Put(v.salt1, m)
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
	
	mpri := <-rs.Next()
	return crypto.UnmarshalPriKey(mpri)
}

func (v voting) baseGetMyVote() (*vutil.VoteInt, error) {
	priv, err := v.getDecryptKey()
	if err != nil{return nil, err}

	m, err := v.ivm.Get(v.myPid+"/"+v.salt1)
	if err != nil{return nil, err}
	vi := &vutil.VoteInt{}
	if err := vi.Unmarshal(m); err != nil{return nil, err}
	return vi, nil
}

func (v voting) baseGetVotes() (<-chan *vutil.VoteInt, int, error) {
	priv, err := v.getDecryptKey()
	if err != nil{return nil, err}

	rs, err := v.ivm.Query(query.Query{
		Filters: []query.Filter{crdt.KeyExistFilter{v.salt1}},
	})
	if err != nil{return nil, -1, err}
	defer rs.Close()

	ch := make(chan *vutil.VoteInt)
	go func(){
		defer close(ch)
		for m := range rs.Next(){
			mvi, err := priv.Decrypt(m)
			if err != nil{continue}

			vi := &vutil.VoteInt{}
			if err := vi.Unmarshal(mvi); err != nil{continue}
			ch <- vi
		}
	}

	rs2, err := v.ivm.Query(query.Query{
		Filters: []query.Filter{crdt.KeyExistFilter{v.salt1}},
		KeysOnly: true,
	})
	if err != nil{return nil, -1, err}
	defer rs2.Close()

	resList, err := rs2.Rest()
	if err != nil{return nil, -1, err}

	return ch, len(resList), nil
}

