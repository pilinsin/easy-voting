package votingutil

import (
	"time"

	"github.com/pilinsin/ipfs-util"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
	rutil "github.com/pilinsin/easy-voting/registration/util"
)

type votingBox struct {
	manEncVote  []byte
	t time.Time
	sign []byte
}
func NewVotingBox() *votingBox {
	return &votingBox{}
}
func (vb *votingBox) Vote(vi VoteInt, identity *rutil.UserIdentity, manPubKey crypto.IPubKey) {
	mev, err := manPubKey.Encrypt(vi.Marshal())
	if err != nil {
		return
	}
	tm := time.Now()
	vt := &struct{
		V: []byte
		T: time.Time
	}{mev, tm}

	mvt, _ := util.Marshal(vt)
	sgn, err := identity.Sign().Sign(mvt)
	if err != nil {
		return
	}

	vb.manEncVote = mev
	vb.t = tm
	vb.sign = sgn
}
func (vb votingBox) GetVote(tInfo *util.TimeInfo, manPriKey crypto.IPriKey) (VoteInt, error) {
	vi := VoteInt{}
	if ok := tInfo.WithinTime(vb.t); !ok{
		return vi, util.NewError("the time of this vote is invalid")
	}
	mvi, err := manPriKey.Decrypt(vb.manEncVote)
	if err != nil{return vi, err}
	err := vi.Unmarshal(mvi)
	return vi, err
}
func (vb votingBox) Verify(verfKey crypto.IVerfKey) (bool, error) {
	vt := &struct{
		V: []byte
		T: time.Time
	}{vb.manEncVote, vb.t}
	mvt, _ := util.Marshal(vt)
	return verfKey.Verify(mvt, vb.sign)
}
func (vb votingBox) WithinTime(tInfo *util.TimeInfo) bool{
	return tInfo.WithinTime(vb.t)
}
func (vb votingBox) Marshal() []byte {
	mvb := &struct {
		M []byte
		T time.Time
		S []byte
	}{vb.manEncVote, vb.t, vb.sign}
	m, _ := util.Marshal(mvb)
	return m
}
func UnmarshalVotingBox(m []byte) (*votingBox, error) {
	mvb := &struct {
		M []byte
		T time.Time
		S []byte
	}{}
	err := util.Unmarshal(m, mvb)
	if err != nil {
		return err
	}

	return &votingBox{mvb.M, mvb.T, mvb.S}, nil
}

type VoteInt map[string]int
func (vi *VoteInt) Marshal() []byte{
	m, _ := util.Marshal(vi)
	return m
}
func (vi *VoteInt) Unmarshal(m []byte) error{
	return util.Unmarshal(m)
}