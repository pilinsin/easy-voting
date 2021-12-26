package registrationutil

import (
	"EasyVoting/ipfs"
)

type ConstHashMap struct {
	rm *ipfs.ReccurentMap
}

func NewConstHashMap(hashes []UhHash, capacity int, is *ipfs.IPFS) *ConstHashMap {
	chm := &ConstHashMap{
		rm: ipfs.NewReccurentMap(capacity),
	}
	for _, hash := range hashes {
		chm.rm.Append(hash, nil, is)
	}
	return chm
}
func (chm *ConstHashMap) Append(hash UhHash, is *ipfs.IPFS) {
	chm.rm.Append(hash, nil, is)
}
func (chm ConstHashMap) ContainHash(hash UhHash, is *ipfs.IPFS) bool {
	if ok := chm.rm.Len(is) == 0; ok {
		return true
	}

	_, ok := chm.rm.ContainKey(hash, is)
	return ok
}
func (chm ConstHashMap) Marshal() []byte {
	return chm.rm.Marshal()
}
func (chm *ConstHashMap) Unmarshal(m []byte) error {
	return chm.rm.Unmarshal(m)
}
func (chm *ConstHashMap) FromCid(chmCid string, is *ipfs.IPFS) error {
	mchm, err := ipfs.FromCid(chmCid, is)
	if err != nil {
		return err
	}
	return chm.Unmarshal(mchm)
}
