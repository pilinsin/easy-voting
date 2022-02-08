package registrationutil

import (
	"github.com/pilinsin/ipfs-util"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)

type registrationBox struct {
	userPubKey  crypto.IPubKey
}
func NewRegistrationBox(pubKey crypto.IPubKey) *registrationBox {
	return &registrationBox{pubKey}
}
func (rb registrationBox) Public() crypto.IPubKey {
	return rb.userPubKey
}
func (rb registrationBox) Marshal() []byte {
	mrb := &struct {
		PubKey []byte
	}{rb.userPubKey.Marshal()}
	m, _ := util.Marshal(mrb)
	return m
}
func UnmarshalRegistrationBox(m []byte) (*registrationBox, error) {
	mrb := &struct {
		PubKey []byte
	}{}
	err := util.Unmarshal(m, mrb)
	if err != nil {
		return nil, err
	}
	pubKey, err := crypto.UnmarshalPubKey(mrb.PubKey)
	if  err != nil {
		return nil, err
	}

	return &registrationBox{pubKey}, nil
}

type UserIdentity struct {
	userHash    UserHash
	userPubKey crypto.IPubKey
	userPriKey  crypto.IPriKey
	userSignKey crypto.ISignKey
	userVerfKey crypto.IVerfKey
}

func NewUserIdentity(uhHash UserHash, pub crypto.IPubKey, pri crypto.IPriKey, sign crypto.ISignKey, verf crypto.IVerfKey) *UserIdentity {
	return &UserIdentity{uhHash, pub, pri, sign, verf}
}
func (ui UserIdentity) UserHash() UserHash {
	return ui.userHash
}
func (ui UserIdentity) Public() crypto.IPubKey {
	return ui.userPubKey
}
func (ui UserIdentity) Private() crypto.IPriKey {
	return ui.userPriKey
}
func (ui UserIdentity) Sign() crypto.ISignKey {
	return ui.userSignKey
}
func (ui UserIdentity) Verf() crypto.IVerfKey {
	return ui.userVerfKey
}
func (ui UserIdentity) Marshal() []byte {
	mui := &struct {
		UserHash    UserHash
		UserPubKey  []byte
		UserPriKey  []byte
		UserSignKey []byte
		UserVerfKey []byte
	}{ui.userHash, ui.userPubKey.Marshal(), ui.userPriKey.Marshal(), ui.userSignKey.Marshal(), ui.userVerfKey.Marshal()}
	m, _ := util.Marshal(mui)
	return m
}
func (ui *UserIdentity) Unmarshal(m []byte) error {
	mui := &struct {
		UserHash    UserHash
		UserPubKey  []byte
		UserPriKey  []byte
		UserSignKey []byte
		UserVerfKey []byte
	}{}
	if err := util.Unmarshal(m, mui); err != nil {
		return err
	}

	pubKey, err := crypto.UnmarshalPubKey(mui.UserPubKey)
	if err != nil {
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
	verfKey, err := crypto.UnmarshalVerfKey(mui.UserVerfKey)
	if err != nil {
		return err
	}

	ui.userHash = mui.UserHash
	ui.userPubKey = pubKey
	ui.userPriKey = priKey
	ui.userSignKey = signKey
	ui.userVerfKey = verfKey
	return nil
}

//sent for registration
type UserInfo struct {
	userHash  UserHash
	rBox *registrationBox
}
func NewUserInfo(uhHash UserHash, rBox *registrationBox) *UserInfo {
	return &UserInfo{uhHash, rBox}
}
func (ui UserInfo) UserHash() UserHash {
	return ui.userHash
}
func (ui UserInfo) RegistrationBox() *registrationBox {
	return ui.rBox
}
func (ui UserInfo) Marshal() []byte {
	mui := &struct {
		UserHash  UserHash
		MRBox []byte
	}{ui.userHash, ui.rBox.Marshal()}
	m, _ := util.Marshal(mui)
	return m
}
func (ui *UserInfo) Unmarshal(m []byte) error {
	mui := &struct {
		UserHash  UserHash
		MRBox []byte
	}{}
	err := util.Unmarshal(m, mui)
	if err != nil {
		return err
	}
	rBox, err := UnmarshalRegistrationBox(mui.MRBox)
	if err != nil {
		return err
	}

	ui.userHash = mui.UserHash
	ui.rBox = mui.rBox
	return nil
}
