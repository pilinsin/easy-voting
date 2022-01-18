package registrationutil

import (
	"fmt"

	"EasyVoting/ipfs"
	"EasyVoting/util"
	"EasyVoting/util/crypto"
)

type hashNameData struct {
	userHash  UserHash
	rIpnsName string
}

func NewHashNameData(uInfo *UserInfo) *hashNameData {
	return &hashNameData{
		userHash:  uInfo.userHash,
		rIpnsName: uInfo.rIpnsName,
	}
}
func (hnd hashNameData) Public(is *ipfs.IPFS) crypto.IPubKey {
	rb, _ := RBoxFromName(hnd.rIpnsName, is)
	return rb.Public()
}
func (hnd hashNameData) Verify(is *ipfs.IPFS) crypto.IVerfKey {
	rb, _ := RBoxFromName(hnd.rIpnsName, is)
	return rb.Verify()
}
func (hnd hashNameData) Name() string {
	return hnd.rIpnsName
}
func (hnd hashNameData) Marshal() []byte {
	mhnd := &struct {
		UserHash  UserHash
		RIpnsName string
	}{hnd.userHash, hnd.rIpnsName}
	m, _ := util.Marshal(mhnd)
	return m
}
func (hnd *hashNameData) Unmarshal(m []byte) error {
	mhnd := &struct {
		UserHash  UserHash
		RIpnsName string
	}{}
	err := util.Unmarshal(m, mhnd)
	if err != nil {
		return err
	}

	hnd.userHash = mhnd.UserHash
	hnd.rIpnsName = mhnd.RIpnsName
	return nil
}

type keyValue struct {
	key   UhHash
	value *hashNameData
}

func (kv keyValue) Key() UhHash {
	return kv.key
}
func (kv keyValue) Value() *hashNameData {
	return kv.value
}

type HashNameMap struct {
	sm *ipfs.ScalableMap
}

func NewHashNameMap(capacity int) *HashNameMap {
	return &HashNameMap{
		sm: ipfs.NewScalableMap(capacity),
	}
}
func (hnm HashNameMap) Next(is *ipfs.IPFS) <-chan *hashNameData {
	ch := make(chan *hashNameData)
	go func() {
		defer close(ch)
		for m := range hnm.sm.Next(is) {
			hnd := &hashNameData{}
			err := hnd.Unmarshal(m)
			if err == nil {
				ch <- hnd
			}
		}
	}()
	return ch
}
func (hnm HashNameMap) NextKeyValue(is *ipfs.IPFS) <-chan *keyValue {
	ch := make(chan *keyValue)
	go func() {
		defer close(ch)
		for mkv := range hnm.sm.NextKeyValue(is) {
			hnd := &hashNameData{}
			err := hnd.Unmarshal(mkv.Value())
			if err == nil {
				ch <- &keyValue{UhHash(mkv.Key()), hnd}
			}
		}
	}()
	return ch
}
func (hnm HashNameMap) ContainHash(hash UhHash, is *ipfs.IPFS) (*hashNameData, bool) {
	if m, ok := hnm.sm.ContainKey(hash, is); !ok {
		return nil, false
	} else {
		hnd := &hashNameData{}
		err := hnd.Unmarshal(m)
		if err != nil {
			return nil, false
		} else {
			return hnd, true
		}
	}
}
func (hnm *HashNameMap) Append(uInfo *UserInfo, salt string, is *ipfs.IPFS) error{
	if _, err := RBoxFromName(uInfo.Name(), is); err != nil{
		fmt.Println("hnm.Append err: ", err)
		return err
	}else{
		hash := NewUhHash(is, salt, uInfo.userHash)
		data := NewHashNameData(uInfo)
		return hnm.sm.Append(hash, data.Marshal(), is)
	}
}

//Verify no falsification
func (hnm HashNameMap) VerifyCid(cid string, is *ipfs.IPFS) bool {
	return hnm.sm.ContainCid(cid, is)
}
func (hnm HashNameMap) VerifyHashes(chm *ConstHashMap, is *ipfs.IPFS) bool {
	if chm.Len(is) == 0 {
		return true
	}
	for hd := range hnm.NextKeyValue(is) {
		uhHash := hd.Key()
		if ok := chm.ContainHash(uhHash, is); !ok {
			return false
		}
	}
	return true
}
func (hnm HashNameMap) VerifyUserInfo(uInfo *UserInfo, salt string, is *ipfs.IPFS) bool {
	uhHash := NewUhHash(is, salt, uInfo.userHash)
	if hnd, ok := hnm.ContainHash(uhHash, is); !ok {
		fmt.Println("verifyUserInfo: not contain uhHash")
		return false
	} else {
		uh := uInfo.userHash == hnd.userHash
		rn := uInfo.rIpnsName == hnd.rIpnsName
		uInfoCid, err1 := ipfs.CidFromName(uInfo.rIpnsName, is)
		hndCid, err2 := ipfs.CidFromName(hnd.rIpnsName, is)
		cid := (uInfoCid == hndCid)
		return uh && rn && cid && (err1 == nil) && (err2 == nil)
	}
}
func (hnm HashNameMap) VerifyUserIdentity(identity *UserIdentity, salt string, is *ipfs.IPFS) bool {
	uhHash := NewUhHash(is, salt, identity.userHash)
	if hnd, ok := hnm.ContainHash(uhHash, is); !ok {
		fmt.Println("verifyUserIdentity: not contain uhHash")
		return false
	} else {
		name, err := identity.KeyFile().Name()
		nm := name == hnd.rIpnsName
		pub := identity.Private().Public().Equals(hnd.Public(is))
		verf := identity.Sign().Verify().Equals(hnd.Verify(is))
		return (err == nil) && nm && pub && verf
	}
}
func (hnm HashNameMap) Marshal() []byte {
	return hnm.sm.Marshal()
}
func (hnm *HashNameMap) Unmarshal(m []byte) error {
	sm := &ipfs.ScalableMap{}
	if err := sm.Unmarshal(m); err != nil {
		return err
	}
	hnm.sm = sm
	return nil
}
func (hnm *HashNameMap) FromName(hnmName string, is *ipfs.IPFS) error {
	mhnm, err := ipfs.FromName(hnmName, is)
	if err != nil {
		return err
	}
	return hnm.Unmarshal(mhnm)
}
