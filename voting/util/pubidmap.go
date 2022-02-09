package votingutil

import (
	"github.com/pilinsin/ipfs-util"
	scmap "github.com/pilinsin/ipfs-util/scalablemap"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
	rutil "github.com/pilinsin/easy-voting/registration/util"
)

type pubIdMap struct {
	sm  scmap.IScalableMap
	vid string
}
func NewPubIdMap(capacity int, vid string) *pubIdMap {
	return &pubIdMap{
		sm:  scmap.NewScalableMap("const", capacity),
		vid: vid,
	}
}
func (pim *pubIdMap) append(pubKey crypto.IPubKey, uid string, is *ipfs.IPFS) error{
	encUid, err := pubKey.Encrypt(util.AnyStrToBytes64(uid))
	if err != nil {
		return err
	}
	pvHash := NewPubVidHash(pubKey, pim.vid)
	return pim.sm.Append(pvHash, encUid, is)
}
func (pim pubIdMap) ContainIdentity(identity *rutil.UserIdentity, is *ipfs.IPFS) (string, bool) {
	pvHash := NewPubVidHash(identity.Public(), pim.vid)
	if encUid, ok := pim.sm.ContainKey(pvHash, is); ok {
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
func (pim pubIdMap) VerifyIdentity(identity *rutil.UserIdentity, is *ipfs.IPFS) bool {
	pvHash := NewPubVidHash(identity.Public(), pim.vid)
	_, ok := pim.sm.ContainKey(pvHash, is)
	return ok
}
func (pim pubIdMap) Marshal() []byte {
	mpim := &struct {
		M []byte
		V string
	}{pim.sm.Marshal(), pim.vid}
	m,_ := util.Marshal(mpim)
	return m
}
func UnmarshalPubIdMap(m []byte) (*pubIdMap, error) {
	mpim := &struct {
		M []byte
		V string
	}{}
	if err := util.Unmarshal(m, mpim); err != nil {
		return nil, err
	}

	sm, err := scmap.UnmarshalScalableMap(mpim.M)
	return &pubIdMap{sm, mpim.V}, err
}
func (pim *pubIdMap) PubIdMapFromCid(pimCid string, is *ipfs.IPFS) (*pubIdMap, error) {
	m, err := ipfs.File.Get(pimCid, is)
	if err != nil {
		return nil, err
	}
	return UnmarshalPubIdMap(m)
}
