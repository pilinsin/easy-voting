package votingutil

import (
	"encoding/base64"
	"errors"

	pb "github.com/pilinsin/easy-voting/voting/util/pb"
	proto "google.golang.org/protobuf/proto"

	evutil "github.com/pilinsin/easy-voting/util"
)

type ManIdentity struct {
	Priv evutil.IPriKey
}

func (mi ManIdentity) Marshal() []byte {
	mpri, _ := mi.Priv.Raw()
	mManId := &pb.ManIdentity{
		Priv: mpri,
	}
	m, _ := proto.Marshal(mManId)
	return m
}
func (mi *ManIdentity) Unmarshal(m []byte) error {
	mManId := &pb.ManIdentity{}
	if err := proto.Unmarshal(m, mManId); err != nil {
		return err
	}
	mpri := mManId.GetPriv()
	//ecies ver.
	//len(mpri) > 32 -> go-ethereum/crypto/secp256k1.(*BitCurve).ScalarMult panics
	if len(mpri) > 32{return errors.New("invalid IPriKey")}

	priv, err := evutil.UnmarshalPri(mpri)
	if err != nil {
		return err
	}

	mi.Priv = priv
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
