package ipfs

import (
	"EasyVoting/util"
	"fmt"
)

type keyValue struct {
	key   string
	value []byte
}

func (kv keyValue) Key() string {
	return kv.key
}
func (kv keyValue) Value() []byte {
	return kv.value
}

type baseMap map[string][]byte
func (bm baseMap) toMap() map[string][]byte{
	return map[string][]byte(bm)
}
func (bm baseMap) toCid(is *IPFS) string{
	m, _ := util.Marshal(bm)
	return ToCidWithAdd(m, is)
}
func (bm *baseMap) fromCid(cid string, is *IPFS) error{
	m, err := FromCid(cid, is)
	if err != nil{return err}
	return util.Unmarshal(m, bm)
}

type ScalableMap struct {
	bm baseMap
	cids   map[string]struct{}
	capacity int
}
func NewScalableMap(capacity int) *ScalableMap {
	return &ScalableMap{
		bm: make(baseMap, capacity),
		cids: make(map[string]struct{}, 0),
		capacity: capacity,
	}
}
func (sm ScalableMap) Len(is *IPFS) int {
	length := len(sm.bm)
	for cid, _ := range sm.cids {
		bm := baseMap{}
		if err := bm.fromCid(cid, is); err != nil{
			return -1
		}
		length += len(bm)
	}
	return length
}
func (sm ScalableMap) Next(is *IPFS) <-chan []byte {
	ch := make(chan []byte)
	go func() {
		defer close(ch)
		for cid, _ := range sm.cids{
			bm := baseMap{}
			if err := bm.fromCid(cid, is); err != nil{return}
			for _, v := range bm{
				ch <- v
			}
		}
		for _, v := range sm.bm{
			ch <- v
		}
	}()
	return ch
}
func (sm ScalableMap) NextKeyValue(is *IPFS) <-chan *keyValue {
	ch := make(chan *keyValue)
	go func() {
		defer close(ch)
		for cid, _ := range sm.cids{
			bm := baseMap{}
			if err := bm.fromCid(cid, is); err != nil{return}
			for k, v := range bm {
				ch <- &keyValue{k, v}
			}
		}
		for k, v := range sm.bm{
			ch <- &keyValue{k, v}
		}
	}()
	return ch
}
func (sm *ScalableMap) Append(key interface{}, value []byte, is *IPFS) {
	keyStr := fmt.Sprintln(key)
	if _, ok := sm.ContainKey(keyStr, is); ok {
		fmt.Println("rm.Append already contain key")
		return
	}

	sm.bm[keyStr] = value
	if len(sm.bm) >= sm.capacity {
		cid := sm.bm.toCid(is)
		sm.cids[cid] = struct{}{}
	}
}
func (sm ScalableMap) ContainKey(key interface{}, is *IPFS) ([]byte, bool) {
	keyStr := fmt.Sprintln(key)
	
	if v, ok := sm.bm[keyStr]; ok{
		return v, true
	}
	for cid, _ := range sm.cids{
		bm := baseMap{}
		if err := bm.fromCid(cid, is); err != nil{
			return nil, false
		}
		if v, ok := bm[keyStr]; ok{
			return v, true
		}
	}
	return nil, false
}
func (sm ScalableMap) ContainCid(cid string, is *IPFS) bool {
	m, err := FromCid(cid, is)
	if err != nil{return false}
	sm0 := &ScalableMap{}
	if err := sm0.Unmarshal(m); err != nil{
		return false
	}

	if sm0.capacity != sm.capacity{
		return false
	}
	for cid0, _ := range sm0.cids{
		if _, ok := sm.cids[cid0]; !ok{
			return false
		}
	}
	bMap0 := sm0.bm.toMap()
	if util.MapContainMap(sm.bm.toMap(), bMap0){
		return true
	}
	for cid, _ := range sm.cids{
		bm := &baseMap{}
		if err := bm.fromCid(cid, is); err != nil{
			return false
		}
		if util.MapContainMap(bm.toMap(), bMap0){
			return true
		}
	}
	return false
}
func (sm ScalableMap) Marshal() []byte {
	mScalableMap := &struct {
		Bm baseMap
		Cids   map[string]struct{}
		Cap int
	}{sm.bm, sm.cids, sm.capacity}
	m, _ := util.Marshal(mScalableMap)
	return m
}
func (sm *ScalableMap) Unmarshal(m []byte) error {
	mScalableMap := &struct {
		Bm baseMap
		Cids   map[string]struct{}
		Cap int
	}{}
	if err := util.Unmarshal(m, mScalableMap); err == nil {
		sm.bm = mScalableMap.Bm
		sm.cids = mScalableMap.Cids
		sm.capacity = mScalableMap.Cap
		return nil
	} else {
		return err
	}
}
