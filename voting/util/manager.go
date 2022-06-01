package votingutil

import (
	"encoding/base64"
	"errors"
	pb "github.com/pilinsin/easy-voting/voting/util/pb"
	proto "google.golang.org/protobuf/proto"

	"github.com/pilinsin/util/crypto"
)

type ManIdentity struct {
	Priv crypto.IPriKey
	Sign crypto.ISignKey
	Verf crypto.IVerfKey
}

func (mi ManIdentity) Marshal() []byte {
	mpri, _ := crypto.MarshalPriKey(mi.Priv)
	msig, _ := crypto.MarshalSignKey(mi.Sign)
	mver, _ := crypto.MarshalVerfKey(mi.Verf)
	mManId := &pb.ManIdentity{
		Priv: mpri,
		Sign: msig,
		Verf: mver,
	}
	m, _ := proto.Marshal(mManId)
	return m
}
func (mi *ManIdentity) Unmarshal(m []byte) error {
	mManId := &pb.ManIdentity{}
	if err := proto.Unmarshal(m, mManId); err != nil {
		return err
	}

	priv, err := crypto.UnmarshalPriKey(mManId.GetPriv())
	if err != nil {
		return err
	}
	sign, err := crypto.UnmarshalSignKey(mManId.GetSign())
	if err != nil {
		return err
	}
	verf, err := crypto.UnmarshalVerfKey(mManId.GetVerf())
	if err != nil {
		return err
	}

	mi.Priv = priv
	mi.Sign = sign
	mi.Verf = verf
	return nil
}

func (mi ManIdentity) toString() string {
	return base64.URLEncoding.EncodeToString(mi.Marshal())
}
func (mi *ManIdentity) FromString(addr string) error {
	if addr == "" {
		return errors.New("invalid addr")
	}
	m, err := base64.URLEncoding.DecodeString(addr)
	if err != nil {
		return err
	}
	return mi.Unmarshal(m)
}
