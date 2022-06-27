package registration

import (
	"os"
	"testing"

	evutil "github.com/pilinsin/easy-voting/util"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
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

//go test -test.v=true -timeout 1h .
func TestRegistration(t *testing.T) {
	bstrp, err := pv.NewBootstrap(i2p.NewI2pHost)
	checkError(t, err)
	defer bstrp.Close()
	bAddrInfo := bstrp.AddrInfo()
	t.Log("bootstrap AddrInfo: ", bAddrInfo)
	baiStr := pv.AddrInfosToString(bAddrInfo)

	users, labels := userDataset()
	rCfgCid, rs, err := rutil.NewConfig("test_title", users, labels, baiStr)
	checkError(t, err)
	rCfgAddr := baiStr + "/" + rCfgCid
	t.Log("config generated")

	man, err := NewRegistrationWithStores(rCfgAddr, rs.Is, rs.Uhm)
	checkError(t, err)
	t.Log("man registration")

	baseDir2 := "registration_test"
	user, err := NewRegistration(rCfgAddr, baseDir2)
	checkError(t, err)
	t.Log("user registration")

	uidStr, err := user.Registrate("alice", "15", "f")
	checkError(t, err)
	uid := &rutil.UserIdentity{}
	checkError(t, uid.FromString(uidStr))
	t.Log(*uid)

	man.Close()
	user.Close()

	baseDir := evutil.BaseDir("registration", "setup")
	os.RemoveAll(baseDir)
	os.RemoveAll(baseDir2)
}
