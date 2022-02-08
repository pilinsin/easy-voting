package votingutil

import (
	"time"

	"github.com/pilinsin/ipfs-util"
	scmap "github.com/pilinsin/ipfs-util/scalablemap"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
	rutil "github.com/pilinsin/easy-voting/registration/util"
)

type namedVBox struct{
	uvHash UidVidHash
	vb *votingBox
}
func (nvb namedVBox) Key() UidVidHash {
	return nvb.uvHash
}
func (nvb namedVBox) Value() *votingBox {
	return nvb.vb
}
func (nvb *namedVBox) marshal() []byte{
	mnvb := &struct{
		H UidVidHash
		M []byte
	}{nvb.uvHash, nvb.vb.Marshal()}
	m, _ := util.Marshal(mnvb)
	return m
}
func (nvb *namedVBox) unmarshal(m []byte) error{
	mnvb := &struct{
		H UidVidHash
		M []byte
	}{}
	if err := util.Unmarshal(m, mnvb); err != nil{return err}
	vb, err := UnmarshalVotingBox(mnvb.M)
	if err != nil{return err}

	nvb.uvHash = mnvb.H
	nvb.vb = vb
	return nil
}


type idVotingMap struct {
	sm    ipfs.IScalableMap
	tInfo *util.TimeInfo
}
func NewIdVotingMap(capacity int, tInfo *util.TimeInfo) *idVotingMap {
	return &idVotingMap{
		sm:    ipfs.NewScalableMap("var", capacity),
		tInfo: tInfo,
	}
}
func (ivm idVotingMap) Next(is *ipfs.IPFS) <-chan *votingBox {
	ch := make(chan *votingBox)
	go func() {
		defer close(ch)
		for m := range ivm.sm.Next(is) {
			nvb := &namedVBox{}
			if err := nvb.unmarshal(m); err == nil{
				ch <- nvb.vb
			}
		}
	}()
	return ch
}
func (ivm idVotingMap) NextKeyValue(is *ipfs.IPFS) <-chan *namedVBox {
	ch := make(chan *namedVBox)
	go func() {
		defer close(ch)
		for kv := range ivm.sm.NextKeyValue(is) {
			nvb := &namedVBox{}
			if err := nvb.unmarshal(kv.Value()); err == nil{
				ch <- nvb
			}
		}
	}()
	return ch
}
func (ivm *idVotingMap) Append(hash UidVidHash, mvb []byte, verfMap *idVerfKeyMap, is *ipfs.IPFS) {
	verfKey, ok := verfMap.ContainHash(hash, is)
	if !ok{return}
	vb, err := UnmarshalVotingBox(mvb)
	if err != nil{return}
	if ok, err := vb.Verify(verfKey); !ok || err != nil{return}
	if ok := vb.withinTime(ivm.tInfo); !ok{return}
	
	nvb := &namedVBox{hash, vb}
	ivm.sm.Append(hash, nvb.marshal(), is)
}
func (ivm idVotingMap) ContainHash(hash UidVidHash, is *ipfs.IPFS) (*votingBox, bool) {
	if mnvb, ok := ivm.sm.ContainKey(hash, is); !ok {
		return nil, false
	} else {
		nvb := &namedVBox{}
		if err := nvb.unmarshal(mnvb); err != nil{return nil, false}
		return nvb.vb, true
	}
}
func (ivm idVotingMap) VerifyMap(verfMap *idVerfKeyMap, is *ipfs.IPFS) bool{
	if ivm.sm.Len() != verfMap.sm.Len(){return false}
	for kv := range votingMap.NextKeyValue(is){
		verfKey, exist := verfMap.ContainHash(kv.Key(), is)
		if ok := exist && kv.Value().Verify(verfKey); !ok{
			return false
		}
	}
	return true
}
func (ivm idVotingMap) Marshal() []byte {
	mivm := &struct {
		M      []byte
		T *util.TimeInfo
	}{ivm.sm.Marshal(), ivm.tInfo}
	m, _ := util.Marshal(mivm)
	return m
}
func UnmarshalIdVotingMap(m []byte) (*idVotingMap, error) {
	mivm := &struct {
		M      []byte
		T *util.TimeInfo
	}{}
	err := util.Unmarshal(m, mivm)
	if err != nil {
		return nil, err
	}

	sm, err := scmap.UnmarshalScalableMap(mivm.M)
	return &idVotingMap{sm, mivm.T}, err
}
func IdVotingMapFromCid(ivmCid string, is *ipfs.IPFS) (*idVotingMap, error) {
	m, err := ipfs.File.Get(ivmCid, is)
	if err != nil {
		return nil, err
	}
	return UnmarshalIdVotingMap(m)
}
