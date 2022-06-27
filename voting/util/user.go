package votingutil

import (
	pb "github.com/pilinsin/easy-voting/voting/util/pb"
	proto "google.golang.org/protobuf/proto"

	evutil "github.com/pilinsin/easy-voting/util"
)

type UserKeyPair struct {
	signKey evutil.ISignKey
	verfKey evutil.IVerfKey
}

func NewUserKeyPair() *UserKeyPair {
	keyPair := evutil.NewSignKeyPair()
	return &UserKeyPair{keyPair.Sign(), keyPair.Verify()}
}
func (ukp UserKeyPair) Sign() evutil.ISignKey {
	return ukp.signKey
}
func (ukp UserKeyPair) Verify() evutil.IVerfKey {
	return ukp.verfKey
}
func (ukp UserKeyPair) Marshal() []byte {
	msign, _ := ukp.signKey.Raw()
	mverf, _ := ukp.verfKey.Raw()
	muk := &pb.KeyPair{
		SignKey: msign,
		VerfKey: mverf,
	}
	m, _ := proto.Marshal(muk)
	return m
}
func (ukp *UserKeyPair) Unmarshal(m []byte) error {
	muk := &pb.KeyPair{}
	if err := proto.Unmarshal(m, muk); err != nil {
		return err
	}

	signKey, err := evutil.UnmarshalSign(muk.SignKey)
	if err != nil {
		return err
	}
	verfKey, err := evutil.UnmarshalVerf(muk.VerfKey)
	if err != nil {
		return err
	}

	ukp.signKey = signKey
	ukp.verfKey = verfKey
	return nil
}
