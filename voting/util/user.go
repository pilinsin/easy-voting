package votingutil

import(
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)

type UserInfo struct{
	uvHash UidVidHash
	verfKey crypto.IVerfKey
}
func NewUserInfo(hash UidVidHash, verfKey crypto.IVerfKey) *UserInfo{
	return &UserInfo{hash, verfKey}
}
func (ui UserInfo) UvHash() UidVidHash{
	return ui.uvHash
}
func (ui UserInfo) Verify() crypto.IVerfKey{
	return ui.verfKey
}
func (ui UserInfo) Marshal() []byte{
	mui := &struct{
		H UidVidHash
		K []byte
	}{ui.uvHash, ui.verfKey.Marshal()}
	m, _ := util.Marshal(mui)
	return m
}
func (ui *UserInfo) Unmarshal(m []byte) error{
	mui := &struct{
		H UidVidHash
		K []byte
	}{}
	if err := util.Unmarshal(m, mui); err != nil{return err}
	verfKey, err := crypto.UnmarshalVerfKey(mui.K)
	if err != nil{return err}

	ui.uvHash = mui.H
	ui.verfKey = verfKey
	return nil
}


type VoteInfo struct{
	uvHash UidVidHash
	vb *votingBox
}
func (vi VoteInfo) UvHash() UidVidHash {
	return vi.uvHash
}
func (vi VoteInfo) VotingBox() *votingBox {
	return vi.vb
}
func (vi *VoteInfo) marshal() []byte{
	mvi := &struct{
		H UidVidHash
		M []byte
	}{vi.uvhHash, vi.vb.Marshal()}
	m, _ := util.Marshal(mvi)
	return m
}
func (vi *VoteInfo) unmarshal(m []byte) error{
	mvi:= &struct{
		H UidVidHash
		M []byte
	}{}
	if err := util.Unmarshal(m, mvi); err != nil{return err}
	vb, err := UnmarshalVotingBox(mvi.M)
	if err != nil{return err}

	vi.uvHash = mvi.H
	vi.vb = vb
	return nil
}