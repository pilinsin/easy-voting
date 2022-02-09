package votingutil

import (
	"fmt"

	"github.com/pilinsin/ipfs-util"
	scmap "github.com/pilinsin/ipfs-util/scalablemap"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)

type hashVerfMap struct {
	sm scmap.IScalableMap
	tInfo *util.TimeInfo
	vid string
}
func NewHashVerfMap(capacity int, tInfo *util.TimeInfo, vid string) *hashVerfMap {
	return &hashVerfMap{
		sm: scmap.NewScalableMap("const", capacity),
		tInfo: tInfo,
		vid: vid,
	}
}
func (hvkm *hashVerfMap) Append(uInfo *UserInfo, uhm *uvhHashMap, is *ipfs.IPFS) error{
	uvhHash := NewUvhHash(uInfo.UvHash(), hvkm.vid)
	if exist := uhm.ContainHash(uvhHash, is); exist{return util.NewError("already appended")}
	return hvkm.sm.Append(uvhHash, uInfo.Verify().Marshal(), is)
}
func (hvkm hashVerfMap) ContainHash(hash UidVidHash, is *ipfs.IPFS) (crypto.IVerfKey, bool) {
	if m, ok := hvkm.sm.ContainKey(hash, is); !ok {
		return nil, false
	} else {
		verfKey, err := crypto.UnmarshalVerfKey(m)
		return verfKey, err == nil
	}
}
//Verify no falsification
func (hvkm hashVerfMap) VerifyCid(cid string, is *ipfs.IPFS) bool {
	return hvkm.sm.ContainCid(cid, is)
}
func (hvkm hashVerfMap) VerifyUserInfo(uInfo *UserInfo, is *ipfs.IPFS) bool {
	if verfKey, ok := hvkm.ContainHash(uInfo.UvHash(), is); !ok {
		fmt.Println("verifyUserIdentity: not contain uhHash")
		return false
	} else {
		return uInfo.Verify().Equals(verfKey)
	}
}
func (hvkm hashVerfMap) Marshal() []byte {
	mMap := &struct{
		M []byte
		T *util.TimeInfo
		V string
	}{hvkm.sm.Marshal(), hvkm.tInfo, hvkm.vid}
	m, _ := util.Marshal(mMap)
	return m
}
func UnmarshalHashVerfMap(m []byte) (*hashVerfMap, error) {
	mMap := &struct{
		M []byte
		T *util.TimeInfo
		V string
	}{}
	if err := util.Unmarshal(m, mMap); err != nil{return nil, err}
	sm, err := scmap.UnmarshalScalableMap("const", mMap.M)
	return &hashVerfMap{sm, mMap.T, mMap.V}, err
}
func HashVerfMapFromName(hvmName string, is *ipfs.IPFS) (*hashVerfMap, error){
	mivm, err := ipfs.Name.Get(hvmName, is)
	if err != nil {
		return nil, err
	}
	return UnmarshalHashVerfMap(mivm)
}
