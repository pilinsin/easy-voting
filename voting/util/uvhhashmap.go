package votingutil

import (
	"github.com/pilinsin/ipfs-util"
	scmap "github.com/pilinsin/ipfs-util/scalablemap"
)

type uvhHashMap struct {
	sm scmap.IScalableMap
}
func NewUvhHashMap(capacity int, is *ipfs.IPFS) *uvhHashMap {
	uhm := &uvhHashMap{
		sm: scmap.NewScalableMap("const", capacity),
	}
	return uhm
}
func (uhm uvhHashMap) Len() int {
	return uhm.sm.Len()
}
func (uhm *uvhHashMap) append(hash UvhHash, is *ipfs.IPFS) error{
	return uhm.sm.Append(hash, nil, is)
}
func (uhm uvhHashMap) ContainHash(hash UhHash, is *ipfs.IPFS) bool {
	if ok := uhm.sm.Len() == 0; ok {
		return true
	}
	_, ok := uhm.sm.ContainKey(hash, is)
	return ok
}
func (uhm uvhHashMap) Marshal() []byte {
	return uhm.sm.Marshal()
}
func UnmarshalUhHashMap(m []byte) (*uvhHashMap, error) {
	sm, err := scmap.UnmarshalScalableMap("const", m)
	return &uvhHashMap{sm}, err
}
func UhHashMapFromCid(uhmCid string, is *ipfs.IPFS) (*uvhHashMap, error) {
	muhm, err := ipfs.File.Get(uhmCid, is)
	if err != nil {
		return nil, err
	}
	return UnmarshalUhHashMap(muhm)
}
