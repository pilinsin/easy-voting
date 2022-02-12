package votingutil

import (
	"github.com/pilinsin/ipfs-util"
	scmap "github.com/pilinsin/ipfs-util/scalablemap"
	"github.com/pilinsin/util"
)

type namedVBox struct{
	uvhHash UvhHash
	vb *votingBox
}
func (nvb namedVBox) Key() UvhHash {
	return nvb.uvhHash
}
func (nvb namedVBox) Value() *votingBox {
	return nvb.vb
}
func (nvb *namedVBox) marshal() []byte{
	mnvb := &struct{
		H UvhHash
		M []byte
	}{nvb.uvhHash, nvb.vb.Marshal()}
	m, _ := util.Marshal(mnvb)
	return m
}
func (nvb *namedVBox) unmarshal(m []byte) error{
	mnvb := &struct{
		H UvhHash
		M []byte
	}{}
	if err := util.Unmarshal(m, mnvb); err != nil{return err}
	vb, err := UnmarshalVotingBox(mnvb.M)
	if err != nil{return err}

	nvb.uvhHash = mnvb.H
	nvb.vb = vb
	return nil
}


type HashVoteMap struct {
	sm    scmap.IScalableMap
	tInfo *util.TimeInfo
	vid string
}
func NewHashVoteMap(capacity int, tInfo *util.TimeInfo, vid string) *HashVoteMap {
	return &HashVoteMap{
		sm:    scmap.NewScalableMap("var", capacity),
		tInfo: tInfo,
		vid: vid,
	}
}
func (hvtm HashVoteMap) Len() int{
	return hvtm.sm.Len()
}
func (hvtm HashVoteMap) Next(is *ipfs.IPFS) <-chan *votingBox {
	ch := make(chan *votingBox)
	go func() {
		defer close(ch)
		for m := range hvtm.sm.Next(is) {
			nvb := &namedVBox{}
			if err := nvb.unmarshal(m); err == nil{
				ch <- nvb.vb
			}
		}
	}()
	return ch
}
func (hvtm HashVoteMap) NextKeyValue(is *ipfs.IPFS) <-chan *namedVBox {
	ch := make(chan *namedVBox)
	go func() {
		defer close(ch)
		for kv := range hvtm.sm.NextKeyValue(is) {
			nvb := &namedVBox{}
			if err := nvb.unmarshal(kv.Value()); err == nil{
				ch <- nvb
			}
		}
	}()
	return ch
}
func (hvtm *HashVoteMap) Append(vInfo *VoteInfo, hvkm *hashVerfMap, is *ipfs.IPFS) error{
	uvhHash := NewUvhHash(vInfo.UvHash(), hvtm.vid)
	verfKey, ok := hvkm.ContainHash(uvhHash, is)
	if !ok{return util.NewError("invalid vote")}
	vb := vInfo.VotingBox()
	if ok, err := vb.Verify(verfKey); !ok || err != nil{return util.NewError("invalid vote")}
	if ok := vb.withinTime(hvtm.tInfo); !ok{return util.NewError("invalid vote")}
	
	nvb := &namedVBox{uvhHash, vb}
	return hvtm.sm.Append(uvhHash, nvb.marshal(), is)
}
func (hvtm HashVoteMap) ContainHash(hash UvhHash, is *ipfs.IPFS) (*votingBox, bool) {
	if mnvb, ok := hvtm.sm.ContainKey(hash, is); !ok {
		return nil, false
	} else {
		nvb := &namedVBox{}
		if err := nvb.unmarshal(mnvb); err != nil{return nil, false}
		return nvb.vb, true
	}
}
func (hvtm HashVoteMap) VerifyMap(hvkm *hashVerfMap, is *ipfs.IPFS) bool{
	for kv := range hvtm.NextKeyValue(is){
		verfKey, exist := hvkm.ContainHash(kv.Key(), is)
		if ok := exist && kv.Value().Verify(verfKey); !ok{
			return false
		}
	}
	return true
}
func (hvtm HashVoteMap) Marshal() []byte {
	mMap := &struct {
		M      []byte
		T *util.TimeInfo
		V string
	}{hvtm.sm.Marshal(), hvtm.tInfo, hvtm.vid}
	m, _ := util.Marshal(mMap)
	return m
}
func UnmarshalHashVoteMap(m []byte) (*HashVoteMap, error) {
	mMap := &struct {
		M      []byte
		T *util.TimeInfo
		V string
	}{}
	err := util.Unmarshal(m, mMap)
	if err != nil {
		return nil, err
	}

	sm, err := scmap.UnmarshalScalableMap(mMap.M)
	return &HashVoteMap{sm, mMap.T, mMap.V}, err
}
func HashVoteMapFromCid(hvmCid string, is *ipfs.IPFS) (*HashVoteMap, error) {
	m, err := ipfs.File.Get(hvmCid, is)
	if err != nil {
		return nil, err
	}
	return UnmarshalHashVoteMap(m)
}
