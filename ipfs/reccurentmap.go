package ipfs
/*
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

type ReccurentMap struct {
	curMap   map[string][]byte
	cidLog   string
	capacity int
}

func NewReccurentMap(capacity int) *ReccurentMap {
	return &ReccurentMap{
		curMap:   make(map[string][]byte, capacity),
		cidLog:   "",
		capacity: capacity,
	}
}
func (rm ReccurentMap) Len(is *IPFS) int {
	length := 0
	var err error
	for {
		length += len(rm.curMap)
		if rm.cidLog != "" {
			rm, err = rm.FromCidLog(is)
			if err != nil {
				return -1
			}
		} else {
			return length
		}
	}
}
func (rm ReccurentMap) Next(is *IPFS) <-chan []byte {
	ch := make(chan []byte)
	go func() {
		defer close(ch)
		var err error
		for {
			for _, v := range rm.curMap {
				ch <- v
			}
			if rm.cidLog != "" {
				rm, err = rm.FromCidLog(is)
				if err != nil {
					return
				}
			} else {
				return
			}
		}
	}()
	return ch
}
func (rm ReccurentMap) NextKeyValue(is *IPFS) <-chan *keyValue {
	ch := make(chan *keyValue)
	go func() {
		defer close(ch)
		var err error
		for {
			for k, v := range rm.curMap {
				ch <- &keyValue{k, v}
			}
			if rm.cidLog != "" {
				rm, err = rm.FromCidLog(is)
				if err != nil {
					return
				}
			} else {
				return
			}
		}
	}()
	return ch
}
func (rm *ReccurentMap) Append(key interface{}, value []byte, is *IPFS) {
	keyStr := fmt.Sprintln(key)
	if _, ok := rm.ContainKey(keyStr, is); ok {
		fmt.Println("rm.Append already contain key")
		return
	}

	rm.curMap[keyStr] = value
	if len(rm.curMap) >= rm.capacity {
		rm = rm.ToCidLog(is)
	}
}
func (rm ReccurentMap) ContainKey(key interface{}, is *IPFS) ([]byte, bool) {
	keyStr := fmt.Sprintln(key)

	var err error
	for {
		if val, ok := rm.curMap[keyStr]; ok {
			return val, true
		}
		if rm.cidLog == "" {
			return nil, false
		}
		rm, err = rm.FromCidLog(is)
		if err != nil {
			return nil, false
		}
	}
}
func (rm ReccurentMap) ContainCid(cid string, is *IPFS) bool {
	if ToCid(rm.Marshal(), is) == cid {
		return true
	}

	rm0 := ReccurentMap{}
	m, err := FromCid(cid, is)
	if err != nil {
		return false
	}
	err = rm0.Unmarshal(m)
	if err != nil {
		return false
	}

	for {
		if rm.cidLog == cid {
			return true
		}
		if rm.cidLog == rm0.cidLog && util.MapContainMap(rm.curMap, rm0.curMap) {
			return true
		}
		if rm.cidLog == "" {
			return false
		}
		rm, err = rm.FromCidLog(is)
		if err != nil {
			return false
		}
	}
}
func (rm ReccurentMap) ToCidLog(is *IPFS) *ReccurentMap {
	m := rm.Marshal()
	cid := ToCidWithAdd(m, is)
	capacity := rm.capacity

	return &ReccurentMap{
		curMap:   make(map[string][]byte, capacity),
		cidLog:   cid,
		capacity: capacity,
	}
}
func (rm ReccurentMap) FromCidLog(is *IPFS) (ReccurentMap, error) {
	rm0 := ReccurentMap{}
	m, err := FromCid(rm.cidLog, is)
	if err != nil {
		return rm0, err
	}
	err = rm0.Unmarshal(m)
	return rm0, err
}
func (rm ReccurentMap) Marshal() []byte {
	mRecMap := &struct {
		CurMap   map[string][]byte
		CidLog   string
		Capacity int
	}{rm.curMap, rm.cidLog, rm.capacity}
	m, _ := util.Marshal(mRecMap)
	return m
}
func (rm *ReccurentMap) Unmarshal(m []byte) error {
	mRecMap := &struct {
		CurMap   map[string][]byte
		CidLog   string
		Capacity int
	}{}
	err := util.Unmarshal(m, mRecMap)
	if err == nil {
		rm.curMap = mRecMap.CurMap
		rm.cidLog = mRecMap.CidLog
		rm.capacity = mRecMap.Capacity
		return nil
	} else {
		return err
	}
}
*/