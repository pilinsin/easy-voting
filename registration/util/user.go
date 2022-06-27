package registrationutil

import (
	"encoding/base64"
	"errors"

	pb "github.com/pilinsin/easy-voting/registration/util/pb"
	evutil "github.com/pilinsin/easy-voting/util"
	proto "google.golang.org/protobuf/proto"
)

type UserIdentity struct {
	hash    string
	pubKey  evutil.IPubKey
	privKey evutil.IPriKey
}

func NewUserIdentity(userHash string, pub evutil.IPubKey, priv evutil.IPriKey) *UserIdentity {
	return &UserIdentity{userHash, pub, priv}
}
func (ui UserIdentity) UserHash() string {
	return ui.hash
}
func (ui UserIdentity) Public() evutil.IPubKey {
	return ui.pubKey
}
func (ui UserIdentity) Private() evutil.IPriKey {
	return ui.privKey
}
func (ui UserIdentity) Marshal() []byte {
	mpub, _ := ui.pubKey.Raw()
	mpri, _ := ui.privKey.Raw()
	mui := &pb.Identity{
		Hash: ui.hash,
		Pub:  mpub,
		Priv: mpri,
	}
	m, _ := proto.Marshal(mui)
	return m
}
func (ui *UserIdentity) Unmarshal(m []byte) error {
	mui := &pb.Identity{}
	if err := proto.Unmarshal(m, mui); err != nil {
		return err
	}
	pubKey, err := evutil.UnmarshalPub(mui.GetPub())
	if err != nil {
		return err
	}
	privKey, err := evutil.UnmarshalPri(mui.GetPriv())
	if err != nil {
		return err
	}

	ui.hash = mui.GetHash()
	ui.pubKey = pubKey
	ui.privKey = privKey
	return nil
}

func (ui UserIdentity) ToString() string {
	return base64.URLEncoding.EncodeToString(ui.Marshal())
}
func (ui *UserIdentity) FromString(addr string) error {
	if addr == "" {
		return errors.New("invalid addr")
	}
	m, err := base64.URLEncoding.DecodeString(addr)
	if err != nil {
		return err
	}
	return ui.Unmarshal(m)
}
