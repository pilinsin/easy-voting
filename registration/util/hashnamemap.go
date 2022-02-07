package registrationutil

import (
	"fmt"

	"github.com/pilinsin/ipfs-util"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)

type keyValue struct {
	keyHash   string
	rIpnsName string
}
func (kv keyValue) Key() string {
	return kv.keyHash
}
func (kv keyValue) Value() string {
	return kv.rIpnsName
}

type hashNameMap struct {
	sm ipfs.IScalableMap
}
func NewHashNameMap(capacity int) *hashNameMap {
	return &hashNameMap{
		sm: ipfs.NewScalableMap("const", capacity),
	}
}
func (hnm hashNameMap) Next(is *ipfs.IPFS) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for m := range hnm.sm.Next(is) {
			ch <- util.Bytes64ToAnyStr(m)
		}
	}()
	return ch
}
func (hnm hashNameMap) NextKeyValue(is *ipfs.IPFS) <-chan *keyValue {
	ch := make(chan *keyValue)
	go func() {
		defer close(ch)
		for kv := range hnm.sm.NextKeyValue(is) {
			name := util.Bytes64ToAnyStr(kv.Value())
			ch <- &keyValue{kv.Key(), name}
		}
	}()
	return ch
}
func (hnm hashNameMap) ContainHash(hash UhHash, is *ipfs.IPFS) (string, bool) {
	if m, ok := hnm.sm.ContainKey(hash, is); !ok {
		return "", false
	} else {
		return util.Bytes64ToAnyStr(m), true
	}
}
func (hnm *hashNameMap) Append(uInfo *UserInfo, salt, uhmCid string, is *ipfs.IPFS) error{
	if _, err := RBoxFromName(uInfo.Name(), is); err != nil{
		fmt.Println("hnm.Append err: ", err)
		return err
	}else{
		uhm, err := UhHashMapFromCid(uhmCid, is)
		if err != nil{return err}

		hash := NewUhHash(salt, uInfo.UserHash())
		if ok := uhm.ContainHash(hash, is); !ok{
			return util.NewError("hnm.Append err: not contain hash")
		}

		data := util.AnyStrToBytes64(uInfo.Name())
		return hnm.sm.Append(hash, data, is)
	}
}
//Verify no falsification
func (hnm hashNameMap) VerifyCid(cid string, is *ipfs.IPFS) bool {
	mhnm, err := ipfs.File.Get(cid, is)
	if err != nil {
		return false
	}
	hnm2, err := UnmarshalHashNameMap(mhnm)
	if err != nil {
		return false
	}
	return hnm.sm.ContainCid(hnm2.sm, is)
}
func (hnm hashNameMap) VerifyHashes(uhmCid string, is *ipfs.IPFS) bool {
	uhm, err := UhHashMapFromCid(uhmCid, is)
	if err != nil{return false}

	if uhm.Len(is) == 0 {
		return true
	}
	for hn := range hnm.NextKeyValue(is) {
		if ok := uhm.ContainHash(hn.Key(), is); !ok {
			return false
		}
	}
	return true
}
func (hnm hashNameMap) VerifyUserInfo(uInfo *UserInfo, salt string, is *ipfs.IPFS) bool {
	uhHash := NewUhHash(salt, uInfo.userHash)
	if name, ok := hnm.ContainHash(uhHash, is); !ok {
		fmt.Println("verifyUserInfo: not contain uhHash")
		return false
	} else {
		rn := uInfo.rIpnsName == name
		uInfoCid, err1 := ipfs.Name.GetCid(uInfo.rIpnsName, is)
		hndCid, err2 := ipfs.Name.GetCid(name, is)
		cid := (uInfoCid == hndCid)
		return rn && cid && (err1 == nil) && (err2 == nil)
	}
}
func (hnm hashNameMap) VerifyUserIdentity(identity *UserIdentity, salt string, is *ipfs.IPFS) bool {
	uhHash := NewUhHash(salt, identity.userHash)
	if name, ok := hnm.ContainHash(uhHash, is); !ok {
		fmt.Println("verifyUserIdentity: not contain uhHash")
		return false
	} else {
		kfName, err := identity.KeyFile().Name()
		nm := kfName == name
		return (err == nil) && nm
	}
}
func (hnm hashNameMap) Marshal() []byte {
	return hnm.sm.Marshal()
}
func UnmarshalHashNameMap(m []byte) (*hashNameMap, error) {
	sm, err := ipfs.UnmarshalScalableMap("const", m)
	return &hashNameMap{sm}, err
}
func HashNameMapFromName(hnmName string, is *ipfs.IPFS) (*hashNameMap, error) {
	mhnm, err := ipfs.FromName(hnmName, is)
	if err != nil {
		return nil, err
	}
	return UnmarshalHashNameMap(mhnm), nil
}
