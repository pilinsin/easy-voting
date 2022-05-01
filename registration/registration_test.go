package registration

import (
	"context"
	"testing"

	rutil "github.com/pilinsin/easy-voting/registration/util"
	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
	crdt "github.com/pilinsin/p2p-verse/crdt"
	"github.com/pilinsin/util/crypto"
)

func checkError(t *testing.T, err error, args ...interface{}) {
	if err != nil {
		args0 := make([]interface{}, len(args)+1)
		args0[0] = err
		copy(args0[1:], args)

		t.Fatal(args0...)
	}
}

func userDataset() (<-chan []string, []string) {
	labels := []string{"name", "age", "sex"}
	users := [][]string{
		{"alice", "15", "f"},
		{"bob", "17", "m"},
		{"c", "19", "f"},
		{"d", "22", "m"},
		{"e", "58", "m"},
		{"f", "32", "m"},
		{"g", "19", "f"},
		{"h", "98", "f"},
		{"i", "39", "f"},
	}
	ch := make(chan []string)
	go func() {
		defer close(ch)
		for _, user := range users {
			ch <- user
		}
	}()
	return ch, labels
}

func genKp() (crdt.IPrivKey, crdt.IPubKey, error) {
	kp := crypto.NewSignKeyPair()
	return kp.Sign(), kp.Verify(), nil
}
func marshalPub(pub crdt.IPubKey) ([]byte, error) {
	return crypto.MarshalVerfKey(pub.(crypto.IVerfKey))
}
func unmarshalPub(m []byte) (crdt.IPubKey, error) {
	return crypto.UnmarshalVerfKey(m)
}

//go test -test.v=true -timeout 1h .
func TestRegistration(t *testing.T) {
	crdt.InitCryptoFuncs(genKp, marshalPub, unmarshalPub)

	bstrp, err := pv.NewBootstrap(i2p.NewI2pHost)
	checkError(t, err)
	defer bstrp.Close()
	bAddrInfo := bstrp.AddrInfo()
	t.Log("bootstrap AddrInfo: ", bAddrInfo)
	baiStr := pv.AddrInfosToString(bAddrInfo)

	users, labels := userDataset()
	rCfgAddr, manIdStr, err := rutil.NewConfig("test_title", users, labels, baiStr)
	checkError(t, err)
	t.Log("config generated")

	man, err := NewRegistration(context.Background(), rCfgAddr, manIdStr)
	checkError(t, err)
	t.Log("man registration")

	user, err := NewRegistration(context.Background(), rCfgAddr, "")
	checkError(t, err)
	t.Log("user registration")

	uidStr, err := user.Registrate("alice", "15", "f")
	checkError(t, err)
	uid := &rutil.UserIdentity{}
	checkError(t, uid.FromString(uidStr))
	t.Log(*uid)

	man.Close()
	user.Close()

}
