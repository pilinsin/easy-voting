package votingutil

import (
	"github.com/pilinsin/easy-voting/ipfs"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	"github.com/pilinsin/easy-voting/util"
)

type NameIdMap struct {
	sm  *ipfs.ScalableMap
	vid string
}

func NewNameIdMap(capacity int, vid string) *NameIdMap {
	return &NameIdMap{
		sm:  ipfs.NewScalableMap(capacity),
		vid: vid,
	}
}
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
	nim.sm.Append(nvHash, encUid, is)
}
func (nim NameIdMap) ContainIdentity(identity *rutil.UserIdentity, is *ipfs.IPFS) (string, bool) {
	name, err := identity.KeyFile().Name()
	if err != nil {
		return "", false
	}
	nvHash := NewNameVidHash(name, nim.vid)
	if encUid, ok := nim.sm.ContainKey(nvHash, is); ok {
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
	_, ok := nim.sm.ContainKey(nvHash, is)
	return ok
}
func (nim NameIdMap) Marshal() []byte {
	mnim := &struct {
		N []byte
		V string
	}{nim.sm.Marshal(), nim.vid}
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

	sm := &ipfs.ScalableMap{}
	if err := sm.Unmarshal(mnim.N); err != nil {
		return err
	}
	nim.sm = sm
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
