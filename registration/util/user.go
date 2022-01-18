package registrationutil

import (
	"EasyVoting/ipfs"
	"EasyVoting/util"
	"EasyVoting/util/crypto"
)

type registrationBox struct {
	userPubKey  crypto.IPubKey
	userVerfKey crypto.IVerfKey
}

func NewRegistrationBox(pubKey crypto.IPubKey, verfKey crypto.IVerfKey) *registrationBox {
	return &registrationBox{pubKey, verfKey}
}
func (rb registrationBox) Public() crypto.IPubKey {
	return rb.userPubKey
}
func (rb registrationBox) Verify() crypto.IVerfKey {
	return rb.userVerfKey
}
func RBoxFromName(rIpnsName string, is *ipfs.IPFS) (*registrationBox, error) {
	m, err := ipfs.FromName(rIpnsName, is)
	if err != nil {
		return nil, err
	}
	rb := &registrationBox{}
	err = rb.Unmarshal(m)
	if err != nil {
		return nil, err
	}
	if rb.userPubKey == nil {
		return nil, util.NewError("userPubKey is nil")
	}
	return rb, nil
}
func (rb registrationBox) Marshal() []byte {
	mrb := &struct {
		PubKey, VerfKey []byte
	}{rb.userPubKey.Marshal(), rb.userVerfKey.Marshal()}
	m, _ := util.Marshal(mrb)
	return m
}
func (rb *registrationBox) Unmarshal(m []byte) error {
	mrb := &struct {
		PubKey, VerfKey []byte
	}{}
	err := util.Unmarshal(m, mrb)
	if err != nil {
		return err
	}

	pubKey, err := crypto.UnmarshalPubKey(mrb.PubKey)
	if  err != nil {
		return err
	}
	verfKey, err := crypto.UnmarshalVerfKey(mrb.VerfKey)
	if err != nil {
		return err
	}
	rb.userPubKey = pubKey
	rb.userVerfKey = verfKey
	return nil
}

type UserIdentity struct {
	userHash    UserHash
	rKeyFile    *ipfs.KeyFile
	userPriKey  crypto.IPriKey
	userSignKey crypto.ISignKey
}

func NewUserIdentity(uhHash UserHash, kf *ipfs.KeyFile, pri crypto.IPriKey, sign crypto.ISignKey) *UserIdentity {
	return &UserIdentity{uhHash, kf, pri, sign}
}
func (ui UserIdentity) UserHash() UserHash {
	return ui.userHash
}
func (ui UserIdentity) KeyFile() *ipfs.KeyFile {
	return ui.rKeyFile
}
func (ui UserIdentity) Private() crypto.IPriKey {
	return ui.userPriKey
}
func (ui UserIdentity) Sign() crypto.ISignKey {
	return ui.userSignKey
}
func (ui UserIdentity) Marshal() []byte {
	mui := &struct {
		UserHash    UserHash
		RKeyFile    []byte
		UserPriKey  []byte
		UserSignKey []byte
	}{ui.userHash, ui.rKeyFile.Marshal(), ui.userPriKey.Marshal(), ui.userSignKey.Marshal()}
	m, _ := util.Marshal(mui)
	return m
}
func (ui *UserIdentity) Unmarshal(m []byte) error {
	mui := &struct {
		UserHash    UserHash
		RKeyFile    []byte
		UserPriKey  []byte
		UserSignKey []byte
	}{}
	if err := util.Unmarshal(m, mui); err != nil {
		return err
	}

	kf := &ipfs.KeyFile{}
	if err := kf.Unmarshal(mui.RKeyFile); err != nil {
		return err
	}
	priKey, err := crypto.UnmarshalPriKey(mui.UserPriKey)
	if err != nil {
		return err
	}
	signKey, err := crypto.UnmarshalSignKey(mui.UserSignKey)
	if err != nil {
		return err
	}

	ui.userHash = mui.UserHash
	ui.rKeyFile = kf
	ui.userPriKey = priKey
	ui.userSignKey = signKey
	return nil
}

//sent for registration
type UserInfo struct {
	userHash  UserHash
	rIpnsName string
}

func NewUserInfo(uhHash UserHash, rIpnsName string) *UserInfo {
	return &UserInfo{uhHash, rIpnsName}
}
func (ui UserInfo) UserHash() UserHash {
	return ui.userHash
}
func (ui UserInfo) Name() string {
	return ui.rIpnsName
}
func (ui UserInfo) Marshal() []byte {
	mui := &struct {
		UserHash  UserHash
		RIpnsName string
	}{ui.userHash, ui.rIpnsName}
	m, _ := util.Marshal(mui)
	return m
}
func (ui *UserInfo) Unmarshal(m []byte) error {
	mui := &struct {
		UserHash  UserHash
		RIpnsName string
	}{}
	err := util.Unmarshal(m, mui)
	if err != nil {
		return err
	}

	ui.userHash = mui.UserHash
	ui.rIpnsName = mui.RIpnsName
	return nil
}
