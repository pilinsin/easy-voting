package votingutil

import (
	pb "github.com/pilinsin/easy-voting/voting/util/pb"
	proto "google.golang.org/protobuf/proto"

	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)

type ManIdentity struct {
	Priv	crypto.IPriKey
	Sign	crypto.ISignKey
	Verf	crypto.IVerfKey
	IpfsDir string
	StoreDir string
}

func (mi ManIdentity) Marshal() []byte {
	mManId := &pb.ManIdentity {
		Priv: mi.Priv.Marshal(),
		Sign: mi.Sign.Marshal(),
		Verf: mi.Verf.Marshal(),
		IpfsDir: mi.IpfsDir,
		StoreDir: mi.StoreDir,
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
	mi.IpfsDir = mManId.GetIpfsDir()
	mi.StoreDir = mManId.GetStoreDir()
	return nil
}

func (mi ManIdentity) toString() string{
	return base64.URLEncoding.EncodeToString(mi.Marshal())
}
func (mi *ManIdentity) FromString(addr string) error{
	m, err := base64.URLEncoding.DecodeString(addr)
	if err != nil{return err}
	return mi.Unmarshal(m)
}