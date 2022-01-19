package votingutil

import (
	"fmt"

	"EasyVoting/ipfs"
	"EasyVoting/util"
	"EasyVoting/util/crypto"
)

type idVerfKeyMap struct {
	sm *ipfs.ScalableMap
}

func NewIdVerfKeyMap(capacity int) *idVerfKeyMap {
	return &idVerfKeyMap{
		sm: ipfs.NewScalableMap(capacity),
	}
}
func (ivm idVerfKeyMap) ContainHash(hash UidVidHash, is *ipfs.IPFS) (crypto.IVerfKey, bool) {
	if m, ok := ivm.sm.ContainKey(hash, is); !ok {
		return nil, false
	} else {
		verfKey, err := crypto.UnmarshalVerfKey(m)
		if err != nil {
			return nil, false
		} else {
			return verfKey, true
		}
	}
}
func (ivm *idVerfKeyMap) Append(hash UidVidHash, verfKey crypto.IVerfKey, votingMap *idVotingMap, manPriKey crypto.IPriKey, is *ipfs.IPFS) error{
	vBox, ok := votingMap.ContainHash(hash, is)
	if !ok{
		err := util.NewError("idVerfKeyMap.Append err: not contain hash")
		fmt.Println(err)
		return err
	}
	sv, err := vBox.GetVote(manPriKey)
	if err != nil{
		fmt.Println("idVerfKeyMap.Append err:", err)
		return err
	}
	if ok := sv.Verify(verfKey); ok{
		ivm.sm.Append(hash, verfKey.Marshal(), is)
		return nil
	}else{
		err := util.NewError("idVerfKeyMap.Append err: invalid verfKey")
		fmt.Println(err)
		return err
	}
}
func (ivm idVerfKeyMap)VerifyIds(votingMap *idVotingMap, is *ipfs.IPFS) bool{
	for kv := range votingMap.NextKeyValue(is){
		if _, ok := ivm.ContainHash(kv.Key(), is); !ok{
			return false
		}
	}
	return true
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
	sm := &ipfs.ScalableMap{}
	if err := sm.Unmarshal(m); err != nil {
		return nil, err
	}
	return &idVerfKeyMap{sm}, nil
}
func IdVerfKeyMapFromName(ivmName string, is *ipfs.IPFS) (*idVerfKeyMap, error){
	mivm, err := ipfs.FromName(ivmName, is)
	if err != nil {
		return nil, err
	}
	return UnmarshalIdVerfKeyMap(mivm)
}
