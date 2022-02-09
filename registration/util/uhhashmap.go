package registrationutil

import (
	"github.com/pilinsin/ipfs-util"
	scmap "github.com/pilinsin/ipfs-util/scalablemap"
)

type uhHashMap struct {
	sm scmap.IScalableMap
}
func NewUhHashMap(capacity int, is *ipfs.IPFS) *uhHashMap {
	uhm := &uhHashMap{
		sm: scmap.NewScalableMap("const", capacity),
	}
	return uhm
}
func (uhm *uhHashMap) Len() int {
	return uhm.sm.Len()
}
func (uhm *uhHashMap) append(hash UhHash, is *ipfs.IPFS) error{
	return uhm.sm.Append(hash, nil, is)
}
func (uhm uhHashMap) ContainHash(hash UhHash, is *ipfs.IPFS) bool {
	if ok := uhm.sm.Len() == 0; ok {
		return true
	}
	_, ok := uhm.sm.ContainKey(hash, is)
	return ok
}
func (uhm uhHashMap) Marshal() []byte {
	return uhm.sm.Marshal()
}
func UnmarshalUhHashMap(m []byte) (*uhHashMap, error) {
	sm, err := scmap.UnmarshalScalableMap("const", m)
	return &uhHashMap{sm}, err
}
func UhHashMapFromCid(uhmCid string, is *ipfs.IPFS) (*uhHashMap, error) {
	muhm, err := ipfs.File.Get(uhmCid, is)
	if err != nil {
		return nil, err
	}
	return UnmarshalUhHashMap(muhm)
}
