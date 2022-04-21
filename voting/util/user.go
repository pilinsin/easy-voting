package votingutil

import(
	pb "github.com/pilinsin/easy-voting/voting/util/pb"
	proto "google.golang.org/protobuf/proto"

	"github.com/pilinsin/util/crypto"
)

type UserKeyPair struct{
	signKey crypto.ISignKey
	verfKey crypto.IVerfKey
}
func NewUserKeyPair() *UserKeyPair{
	keyPair := crypto.NewSignKeyPair()
	return &UserKeyPair{keyPair.Sign(), keyPair.Verify()}
}
func (ukp UserKeyPair) Sign() crypto.ISignKey{
	return ukp.signKey
}
func (ukp UserKeyPair) Verify() crypto.IVerfKey{
	return ukp.verfKey
}
func (ukp UserKeyPair) Marshal() []byte{
	msign, _ := crypto.MarshalSignKey(ukp.signKey)
	mverf, _ := crypto.MarshalVerfKey(ukp.verfKey)
	muk := &pb.KeyPair{
		SignKey: msign,
		VerfKey: mverf,
	}
	m, _ := proto.Marshal(muk)
	return m
}
func (ukp *UserKeyPair) Unmarshal(m []byte) error{
	muk := &pb.KeyPair{}
	if err := proto.Unmarshal(m, muk); err != nil{return err}

	signKey, err := crypto.UnmarshalSignKey(muk.SignKey)
	if err != nil{return err}
	verfKey, err := crypto.UnmarshalVerfKey(muk.VerfKey)
	if err != nil{return err}

	ukp.signKey = signKey
	ukp.verfKey = verfKey
	return nil
}

