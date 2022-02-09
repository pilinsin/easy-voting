package registrationutil

import (
	"fmt"

	"github.com/pilinsin/ipfs-util"
	scmap "github.com/pilinsin/ipfs-util/scalablemap"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)

type keyValue struct {
	keyHash   string
	rBox *registrationBox
}
func (kv keyValue) Key() string {
	return kv.keyHash
}
func (kv keyValue) Value() string {
	return kv.rBox
}

type hashBoxMap struct {
	sm scmap.IScalableMap
}
func NewHashBoxMap(capacity int) *hashBoxMap {
	return &hashBoxMap{
		sm: scmap.NewScalableMap("const", capacity),
	}
}
func (hbm hashBoxMap) Next(is *ipfs.IPFS) <-chan *registrationBox {
	ch := make(chan *registrationBox)
	go func() {
		defer close(ch)
		for m := range hbm.sm.Next(is) {
			rBox, err := UnmarshalRegistrationBox(m)
			if err == nil{
				ch <- rBox
			}
		}
	}()
	return ch
}
func (hbm hashBoxMap) NextKeyValue(is *ipfs.IPFS) <-chan *keyValue {
	ch := make(chan *keyValue)
	go func() {
		defer close(ch)
		for kv := range hbm.sm.NextKeyValue(is) {
			rBox, err := UnmarshalRegistrationBox(kv.Value())
			if err == nil{
				ch <- &keyValue{kv.Key(), rBox}
			}			
		}
	}()
	return ch
}
func (hbm hashBoxMap) ContainHash(hash UhHash, is *ipfs.IPFS) (*registrationBox, bool) {
	if m, ok := hbm.sm.ContainKey(hash, is); !ok {
		return nil, false
	} else {
		rBox, err := UnmarshalRegistrationBox(m)
		return rBox, err == nil
	}
}
func (hbm hashBoxMap) ContainPubKey(pub crypto.IPubKey, is *ipfs.IPFS) bool{
	for rBox := range hbm.Next(r.is){
		if rBox.Public().Equals(pub){return true}
	}
	return false
}
func (hbm *hashBoxMap) Append(uInfo *UserInfo, salt, uhmCid string, is *ipfs.IPFS) error{
	uhm, err := UhHashMapFromCid(uhmCid, is)
	if err != nil{return err}
	hash := NewUhHash(salt, uInfo.UserHash())
	if ok := uhm.ContainHash(hash, is); !ok{
		return util.NewError("Append: not contain hash")
	}

	data := uInfo.RegistrationBox().Marshal()
	return hbm.sm.Append(hash, data, is)
}
//Verify no falsification
func (hbm hashBoxMap) VerifyCid(cid string, is *ipfs.IPFS) bool {
	mhbm, err := ipfs.File.Get(cid, is)
	if err != nil {
		return false
	}
	hnm2, err := UnmarshalHashNameMap(mhbm)
	if err != nil {
		return false
	}
	return hbm.sm.ContainCid(hnm2.sm, is)
}
func (hbm hashBoxMap) VerifyHashes(uhmCid string, is *ipfs.IPFS) bool {
	uhm, err := UhHashMapFromCid(uhmCid, is)
	if err != nil{return false}

	if uhm.Len(is) == 0 {
		return true
	}
	for hn := range hbm.NextKeyValue(is) {
		if ok := uhm.ContainHash(hn.Key(), is); !ok {
			return false
		}
	}
	return true
}
func (hbm hashBoxMap) VerifyUserInfo(uInfo *UserInfo, salt string, is *ipfs.IPFS) bool {
	uhHash := NewUhHash(salt, uInfo.userHash)
	if rBox, ok := hbm.ContainHash(uhHash, is); !ok {
		fmt.Println("verifyUserInfo: not contain uhHash")
		return false
	} else {
		return uInfo.RegistrationBox().Public().Equals(rBox.Public())
	}
}
func (hbm hashBoxMap) VerifyUserIdentity(identity *UserIdentity, salt string, is *ipfs.IPFS) bool {
	uhHash := NewUhHash(salt, identity.userHash)
	if rBox, ok := hbm.ContainHash(uhHash, is); !ok {
		fmt.Println("verifyUserIdentity: not contain uhHash")
		return false
	} else {
		return rBox.Public().Equals(identity.Public())
	}
}
func (hbm hashBoxMap) Marshal() []byte {
	return hbm.sm.Marshal()
}
func UnmarshalHashBoxMap(m []byte) (*hashBoxMap, error) {
	sm, err := scmap.UnmarshalScalableMap("const", m)
	return &hashBoxMap{sm}, err
}
func HashBoxMapFromName(hnmName string, is *ipfs.IPFS) (*hashBoxMap, error) {
	mhbm, err := ipfs.Name.Get(hnmName, is)
	if err != nil {
		return nil, err
	}
	return UnmarshalHashNameMap(mhbm), nil
}
