package votingutil

import (
	"github.com/pilinsin/ipfs-util"
	scmap "github.com/pilinsin/ipfs-util/scalablemap"
	"github.com/pilinsin/util"
	rutil "github.com/pilinsin/easy-voting/registration/util"
)

type boxIdMap struct {
	sm  ipfs.IScalableMap
	vid string
}
func NewBoxIdMap(capacity int, vid string) *boxIdMap {
	return &boxIdMap{
		sm:  ipfs.NewScalableMap("const", capacity),
		vid: vid,
	}
}
func (bim *boxIdMap) Append(mRBox []byte, uid string, is *ipfs.IPFS) error{
	rb, err := rutil.UnmarshalRegistrationBox(mRBox)
	if err != nil {
		return
	}
	encUid, err := rb.Public().Encrypt(util.AnyStrToBytes64(uid))
	if err != nil {
		return
	}
	bvHash := NewBoxVidHash(rb.Public().Marshal(), bim.vid)
	bim.sm.Append(bvHash, encUid, is)
}
func (bim boxIdMap) ContainIdentity(identity *rutil.UserIdentity, is *ipfs.IPFS) (string, bool) {
	bvHash := NewBoxVidHash(identity.Public().Marshal(), bim.vid)
	if encUid, ok := bim.sm.ContainKey(bvHash, is); ok {
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
func (bim boxIdMap) VerifyIdentity(identity *rutil.UserIdentity, is *ipfs.IPFS) bool {
	bvHash := NewBoxVidHash(identity.Public().Marshal(), bim.vid)
	_, ok := bim.sm.ContainKey(bvHash, is)
	return ok
}
func (bim boxIdMap) Marshal() []byte {
	mbim := &struct {
		M []byte
		V string
	}{bim.sm.Marshal(), bim.vid}
	m, _ := util.Marshal(mbim)
	return m
}
func UnmarshalBoxIdMap(m []byte) (*boxIdMap, error) {
	mbim := &struct {
		M []byte
		V string
	}{}
	if err := util.Unmarshal(m, mbim); err != nil {
		return nil, err
	}

	sm, err := scmap.UnmarshalScalableMap(mbim.M)
	return &boxIdMap{sm, mbim.V}, err
}
func (bim *boxIdMap) BoxIdMapFromCid(nimCid string, is *ipfs.IPFS) (*boxIdMap, error) {
	m, err := ipfs.File.Get(nimCid, is)
	if err != nil {
		return nil, err
	}
	return UnmarshalBoxIdMap(m)
}
