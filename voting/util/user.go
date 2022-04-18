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
	mui := &pb.KeyPair{
		SignKey: ukp.signKey.Marshal(),
		VerfKey: ukp.verfKey.Marshal(),
	}
	m, _ := proto.Marshal(mui)
	return m
}
func (ukp *UserKeyPair) Unmarshal(m []byte) error{
	mui := &pb.KeyPair{}
	if err := proto.Unmarshal(m, mui); err != nil{return err}

	signKey, err := crypto.UnmarshalSignKey(mui.SignKey)
	if err != nil{return err}
	verfKey, err := crypto.UnmarshalVerfKey(mui.VerfKey)
	if err != nil{return err}

	ukp.signKey = signKey
	ukp.verfKey = verfKey
	return nil
}

