package votingutil

import(
	"EasyVoting/util"
	"EasyVoting/util/crypto"
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

