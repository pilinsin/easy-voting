package voting

import (
	"testing"
	"context"
	"time"

	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"

	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
	crdt "github.com/pilinsin/p2p-verse/crdt"
	rgst "github.com/pilinsin/easy-voting/registration"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	riface "github.com/pilinsin/easy-voting/registration/interface"
	vutil "github.com/pilinsin/easy-voting/voting/util"
	//viface "github.com/pilinsin/easy-voting/voting/interface"
)


func checkError(t *testing.T, err error, args ...interface{}) {
	if err != nil {
		args0 := make([]interface{}, len(args)+1)
		args0[0] = err
		copy(args0[1:], args)

		t.Fatal(args0...)
	}
}

func userDataset() (<-chan []string, []string){
	labels := []string{"name", "age", "sex"}
	users :=[][]string{
		{"alice",	"15",	"f"},
		{"bob",		"17",	"m"},
		{"c",		"19",	"f"},
		{"d",		"22",	"m"},
		{"e",		"58",	"m"},
		{"f",		"32",	"m"},
		{"g",		"19",	"f"},
		{"h",		"98",	"f"},
		{"i",		"39",	"f"},
	}
	ch := make(chan []string)
	go func(){
		defer close(ch)
		for _, user := range users{
			ch <- user
		}
	}()
	return ch, labels
}

func genKp() (crdt.IPrivKey, crdt.IPubKey, error){
	kp := crypto.NewSignKeyPair()
	return kp.Sign(), kp.Verify(), nil
}
func marshalPub(pub crdt.IPubKey) ([]byte, error){
	return crypto.MarshalVerfKey(pub.(crypto.IVerfKey))
}
func unmarshalPub(m []byte) (crdt.IPubKey, error){
	return crypto.UnmarshalVerfKey(m)
}

func registrate(t *testing.T, baiStr string) (riface.IRegistration, string, string){
	users, labels := userDataset()
	rCfgAddr, manIdStr, err := rutil.NewConfig("test_rTitle", users, labels, baiStr)
	checkError(t, err)

	man, err := rgst.NewRegistration(context.Background(), rCfgAddr, manIdStr)
	checkError(t, err)

	user, err := rgst.NewRegistration(context.Background(), rCfgAddr, "")
	checkError(t, err)
	defer user.Close()
	
	var uidStr string
	users, _ = userDataset()
	for ud := range users{
		name, age, sex := ud[0], ud[1], ud[2]
		uidStr, err = user.Registrate(name, age, sex)
		checkError(t, err)
	}
	time.Sleep(time.Minute*5)

	return man, rCfgAddr, uidStr
}



func makeTimeInfo(t *testing.T) *util.TimeInfo{
	now := time.Now()
	begin := now.Format(util.Layout)
	end := now.Add(time.Minute*30).Format(util.Layout)
	tInfo, err := util.NewTimeInfo(begin, end, now.Location().String())
	checkError(t, err)
	return tInfo
}
func makeCandidates() []*vutil.Candidate{
	nameGroups := [][]string{
		{"A", "gA"},
		{"B", "gB"},
		{"C", "gC"},
		{"D", "gD"},
		{"E", "gE"},
	}

	cands := make([]*vutil.Candidate, len(nameGroups))
	for idx, ng := range nameGroups{
		cands[idx] = &vutil.Candidate{ng[0], ng[1], "", nil, ""}
	}
	
	return cands
}
func makeVote(name string, cands []*vutil.Candidate) vutil.VoteInt{
	vi := make(vutil.VoteInt)
	for _, cand := range cands{
		k := cand.Name + ", " + cand.Group
		if cand.Name == name{
			vi[k] = 1
		}else{
			vi[k] = 0
		}
	}
	return vi
}
func vote(t *testing.T, baiStr string){
	rMan, rcAddr, uidStr := registrate(t, baiStr)

	ttl := "test_vtitle"
	nv := 1
	ti := makeTimeInfo(t)
	cands := makeCandidates()
	vp := &vutil.VoteParams{0,1,1}
	vt := vutil.Approval
	vCfgAddr, manIdStr, err := vutil.NewConfig(ttl, rcAddr, nv, ti, cands, vp, vt)
	checkError(t, err)
	t.Log("vCfg generated")
	rMan.Close()

	vMan, err := NewVoting(context.Background(), vCfgAddr, manIdStr)
	checkError(t, err)
	t.Log("vMan generated")

	user, err := NewVoting(context.Background(), vCfgAddr, uidStr)
	checkError(t, err)
	t.Log("vUser generated")

/*
	verifiers := make([]viface.IVoting, nv)
	for idx := 0; idx < nv; idx++{
		verifier, err := NewVoting(context.Background(), vCfgAddr, "")
		checkError(t, err)
		verifiers[idx] = verifier
		t.Log("verifier",idx," generated")
	}
	t.Log("verifiers generated")
*/
	checkError(t, user.Vote(makeVote("D", cands)))
	t.Log("user vote")
	time.Sleep(time.Minute*5)

/*
	for _, verifier := range verifiers{
		verifier.Close()
	}
	t.Log("verified")
*/

	//upload manPriKey
	checkError(t, vMan.Vote(makeVote("", cands)))
	t.Log("man upload manPriKey")
	time.Sleep(time.Minute*2)

	vi, err := user.GetMyVote()
	checkError(t, err)
	t.Log("my vote:", *vi)
	res, nVoters, nVoted, err := user.GetResult()
	checkError(t, err)
	t.Log("turnout:", nVoted, "/", nVoters, ", result:", *res)

	vMan.Close()
	user.Close()
}


//go test -test.v=true -timeout 1h .
func TestVoting(t *testing.T){
	crdt.InitCryptoFuncs(genKp, marshalPub, unmarshalPub)

	bstrp, err := pv.NewBootstrap(i2p.NewI2pHost)
	checkError(t, err)
	defer bstrp.Close()
	bAddrInfo := bstrp.AddrInfo()
	t.Log("bootstrap AddrInfo: ", bAddrInfo)
	baiStr := pv.AddrInfosToString(bAddrInfo)

	vote(t, baiStr)
}