package votingutil

import (
	"fmt"

	"github.com/pilinsin/ipfs-util"
	scmap "github.com/pilinsin/ipfs-util/scalablemap"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)

type idVerfKeyMap struct {
	sm ipfs.IScalableMap
	tInfo *util.TimeInfo
}
func NewIdVerfKeyMap(capacity int, tInfo *util.TimeInfo) *idVerfKeyMap {
	return &idVerfKeyMap{
		sm: ipfs.NewScalableMap("const", capacity),
		tInfo: tInfo,
	}
}
func (ivm idVerfKeyMap) ContainHash(hash UidVidHash, is *ipfs.IPFS) (crypto.IVerfKey, bool) {
	if m, ok := ivm.sm.ContainKey(hash, is); !ok {
		return nil, false
	} else {
		verfKey, err := crypto.UnmarshalVerfKey(m)
		return verfKey, err == nil
	}
}
func (ivm *idVerfKeyMap) Append(hash UidVidHash, verfKey crypto.IVerfKey, mvb []byte, is *ipfs.IPFS) error{
	vb, err := UnmarshalVotingBox(mvb)
	if err != nil{return err}
	if ok := vb.WithinTime(ivm.tInfo); !ok{
		return util.NewError("votingBox.t is invalid")
	}
	if ok, err := vb.Verify(verfKey); ok && err == nil{
		ivm.sm.Append(hash, verfKey.Marshal(), is)
		return nil
	}else{
		err := util.NewError("idVerfKeyMap.Append err: invalid verfKey")
		fmt.Println(err)
		return err
	}
}
//Verify no falsification
func (ivm idVerfKeyMap) VerifyCid(cid string, is *ipfs.IPFS) bool {
	return ivm.sm.ContainCid(cid, is)
}
func (ivm idVerfKeyMap) VerifyUserInfo(uInfo *UserInfo, is *ipfs.IPFS) bool {
	if verfKey, ok := ivm.ContainHash(uInfo.UvHash(), is); !ok {
		fmt.Println("verifyUserIdentity: not contain uhHash")
		return false
	} else {
		return uInfo.Verify().Equals(verfKey)
	}
}
func (ivm idVerfKeyMap) Marshal() []byte {
	return ivm.sm.Marshal()
}
func UnmarshalIdVerfKeyMap(m []byte) (*idVerfKeyMap, error) {
	sm, err := scmap.UnmarshalScalableMap("const", m)
	return &idVerfKeyMap{sm}, err
}
func IdVerfKeyMapFromName(ivmName string, is *ipfs.IPFS) (*idVerfKeyMap, error){
	mivm, err := ipfs.Name.Get(ivmName, is)
	if err != nil {
		return nil, err
	}
	return UnmarshalIdVerfKeyMap(mivm)
}
