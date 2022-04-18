package registrationutil

import (
	"errors"
	"encoding/base64"
	"github.com/pilinsin/util/crypto"
	pb "github.com/pilinsin/easy-voting/registration/util/pb"
	proto "google.golang.org/protobuf/proto"
)


type UserIdentity struct {
	hash    string
	pubKey crypto.IPubKey
	privKey crypto.IPriKey
}

func NewUserIdentity(userHash string, pub crypto.IPubKey, priv crypto.IPriKey) *UserIdentity {
	return &UserIdentity{userHash, pub, priv}
}
func (ui UserIdentity) UserHash() string {
	return ui.hash
}
func (ui UserIdentity) Public() crypto.IPubKey {
	return ui.pubKey
}
func (ui UserIdentity) Private() crypto.IPriKey {
	return ui.privKey
}
func (ui UserIdentity) Marshal() []byte {
	mpub, _ := crypto.MarshalPubKey(ui.pubKey)
	mpri, _ := crypto.MarshalPriKey(ui.privKey)
	mui := &pb.Identity{
		Hash:	ui.hash,
		Pub:	mpub,
		Priv:  	mpri,
	}
	m, _ := proto.Marshal(mui)
	return m
}
func (ui *UserIdentity) Unmarshal(m []byte) error {
	mui := &pb.Identity{}
	if err := proto.Unmarshal(m, mui); err != nil {
		return err
	}

	pubKey, err := crypto.UnmarshalPubKey(mui.GetPub())
	if err != nil {
		return err
	}
	privKey, err := crypto.UnmarshalPriKey(mui.GetPriv())
	if err != nil {
		return err
	}

	ui.hash = mui.Hash
	ui.pubKey = pubKey
	ui.privKey = privKey
	return nil
}

func (ui UserIdentity) ToString() string{
	return base64.URLEncoding.EncodeToString(ui.Marshal())
}
func (ui *UserIdentity) FromString(addr string) error{
	if addr == ""{return errors.New("invalid addr")}
	m, err := base64.URLEncoding.DecodeString(addr)
	if err != nil{return err}
	return ui.Unmarshal(m)
}