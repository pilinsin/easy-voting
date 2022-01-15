package registrationutil

import (
	"EasyVoting/ipfs"
)

type ConstHashMap struct {
	sm *ipfs.ScalableMap
}

func NewConstHashMap(hashes []UhHash, capacity int, is *ipfs.IPFS) *ConstHashMap {
	chm := &ConstHashMap{
		sm: ipfs.NewScalableMap(capacity),
	}
	for _, hash := range hashes {
		chm.sm.Append(hash, nil, is)
	}
	return chm
}
func (chm *ConstHashMap) Len(is *ipfs.IPFS) int {
	return chm.sm.Len(is)
}
func (chm *ConstHashMap) Append(hash UhHash, is *ipfs.IPFS) {
	chm.sm.Append(hash, nil, is)
}
func (chm ConstHashMap) ContainHash(hash UhHash, is *ipfs.IPFS) bool {
	if ok := chm.sm.Len(is) == 0; ok {
		return true
	}

	_, ok := chm.sm.ContainKey(hash, is)
	return ok
}
func (chm ConstHashMap) Marshal() []byte {
	return chm.sm.Marshal()
}
func (chm *ConstHashMap) Unmarshal(m []byte) error {
	sm := &ipfs.ScalableMap{}
	if err := sm.Unmarshal(m); err != nil {
		return err
	}
	chm.sm = sm
	return nil
}
func (chm *ConstHashMap) FromCid(chmCid string, is *ipfs.IPFS) error {
	mchm, err := ipfs.FromCid(chmCid, is)
	if err != nil {
		return err
	}
	return chm.Unmarshal(mchm)
}
