package registrationutil

import (
	"errors"
	"encoding/base64"
	pb "github.com/pilinsin/easy-voting/registration/util/pb"
	proto "google.golang.org/protobuf/proto"
)

type ManIdentity struct {
	IpfsDir string
	StoreDir string
}

func (mi ManIdentity) Marshal() []byte {
	mManId := &pb.ManIdentity {
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
	mi.IpfsDir = mManId.GetIpfsDir()
	mi.StoreDir = mManId.GetStoreDir()
	return nil
}

func (mi ManIdentity) toString() string{
	return base64.URLEncoding.EncodeToString(mi.Marshal())
}
func (mi *ManIdentity) FromString(addr string) error{
	if addr == ""{return errors.New("invalid addr")}
	m, err := base64.URLEncoding.DecodeString(addr)
	if err != nil{return err}
	return mi.Unmarshal(m)
}