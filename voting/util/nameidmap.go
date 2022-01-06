package votingutil

import (
	"EasyVoting/ipfs"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
)

type NameIdMap struct {
	rm  *ipfs.ReccurentMap
	vid string
}

func NewNameIdMap(capacity int, vid string) *NameIdMap {
	return &NameIdMap{
		rm:  ipfs.NewReccurentMap(capacity),
		vid: vid,
	}
}

/*
func (nim NameIdMap) Next(is *ipfs.IPFS) <-chan *rutil.RegistrationBox{
	ch := make(chan *rutil.RegistrationBox)
	go func(){
		defer close(ch)
		for kv := range nim.rm.NextKeyValue(is){
			rb := &rutil.RegistrationBox{}
			err := rb.FromName(kv.Key(), is)
			if err == nil{
				ch <- rb
			}
		}
	}()
	return ch
}
*/
func (nim *NameIdMap) Append(rIpnsName, uid string, is *ipfs.IPFS) {
	rb, err := rutil.RBoxFromName(rIpnsName, is)
	if err != nil {
		return
	}
	encUid, err := rb.Public().Encrypt(util.AnyStrToBytes64(uid))
	if err != nil {
		return
	}

	nvHash := NewNameVidHash(rIpnsName, nim.vid)
	nim.rm.Append(nvHash, encUid, is)
}
func (nim NameIdMap) ContainIdentity(identity *rutil.UserIdentity, is *ipfs.IPFS) (string, bool) {
	name, err := identity.KeyFile().Name()
	if err != nil {
		return "", false
	}
	nvHash := NewNameVidHash(name, nim.vid)
	if encUid, ok := nim.rm.ContainKey(nvHash, is); ok {
		bUid, err := identity.Private().Decrypt(encUid)
		if err != nil {
			return "", false
		}
		return util.Bytes64ToAnyStr(bUid), true
	} else {
		return "", false
	}
}

//login verification
func (nim NameIdMap) VerifyIdentity(identity *rutil.UserIdentity, is *ipfs.IPFS) bool {
	name, err := identity.KeyFile().Name()
	if err != nil {
		return false
	}
	nvHash := NewNameVidHash(name, nim.vid)
	_, ok := nim.rm.ContainKey(nvHash, is)
	return ok
}
func (nim NameIdMap) Marshal() []byte {
	mnim := &struct {
		N []byte
		V string
	}{nim.rm.Marshal(), nim.vid}
	m, _ := util.Marshal(mnim)
	return m
}
func (nim *NameIdMap) Unmarshal(m []byte) error {
	mnim := &struct {
		N []byte
		V string
	}{}
	if err := util.Unmarshal(m, mnim); err != nil {
		return err
	}

	rm := &ipfs.ReccurentMap{}
	if err := rm.Unmarshal(mnim.N); err != nil {
		return err
	}
	nim.rm = rm
	nim.vid = mnim.V
	return nil
}
func (nim *NameIdMap) FromCid(nimCid string, is *ipfs.IPFS) error {
	m, err := ipfs.FromCid(nimCid, is)
	if err != nil {
		return err
	}
	return nim.Unmarshal(m)
}
