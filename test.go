package main

import (
	"fmt"
	"time"

	"EasyVoting/ipfs"
	"EasyVoting/util"
	"EasyVoting/util/ecies"
	"EasyVoting/util/ed25519"
	"EasyVoting/voting"
	sv "EasyVoting/voting/single"
)

//sudo sysctl -w net.core.rmem_max=2500000
func main() {
	topic := util.GenUniqueID(50, 50)
	//user
	ctx := util.NewContext()
	Ipfs := ipfs.New(ctx, ".ipfs", topic)
	//manager
	ctx2 := util.NewContext()
	Ipfs2 := ipfs.New(ctx2, ".ipfs2", topic)
	//other
	ctx3 := util.NewContext()
	Ipfs3 := ipfs.New(ctx3, ".ipfs3", topic)
	fmt.Println("ipfs initialized.")
	fmt.Println(time.Now().String())

	Ipfs.PubsubConnect()
	Ipfs.PubsubPublish([]byte("moavesilbn;omva"))
	Ipfs.PubsubPublish([]byte("moavesilbn;omva2"))
	//Ipfs.PubsubSubTest()
	fmt.Println("a")
	Ipfs2.PubsubConnect()
	Ipfs2.PubsubSubTest()
	fmt.Println("u")
	Ipfs3.PubsubConnect()
	Ipfs3.PubsubSubTest()
	fmt.Println("test pubsub finished")

	eciesKeyPair := ecies.GenKeyPair()
	signKeyPair := ed25519.GenKeyPair()
	signKeyPair2 := ed25519.GenKeyPair()

	cfg := voting.InitConfig{
		Is:       Ipfs,
		Topic:    "test_sample",
		VotingID: util.GenUniqueID(30, 30),
		UserID:   util.GenUniqueID(30, 6),
		Begin:    "2021-6-1  00:00am",
		End:      "2021-8-13 11:59pm",
		NCands:   3,
		PubKey:   eciesKeyPair.Public(),
		SignKey:  signKeyPair.Sign(),
	}
	sVoting := sv.New(&cfg, nil)
	v := voting.VoteInt{"A": 0, "B": 1, "C": 0}
	sVoting.Vote(v)
	v2 := voting.VoteInt{"A": 0, "B": 0, "C": 1}
	sVoting.Vote(v2)

	cfg2 := voting.InitConfig{
		Is:       Ipfs2,
		Topic:    "test_sample",
		VotingID: util.GenUniqueID(30, 30),
		UserID:   util.GenUniqueID(30, 6),
		Begin:    "2021-6-1  00:00am",
		End:      "2021-8-13 11:59pm",
		NCands:   3,
		SignKey:  signKeyPair2.Sign(),
	}
	sVoting2 := sv.New(&cfg2, nil)
	sVoting2.BaseVote(sVoting2.MarshalVoteEnd())

	cfg3 := voting.InitConfig{
		Is:       Ipfs3,
		Topic:    "test_sample",
		VotingID: "nil",
		UserID:   "nil",
		Begin:    "nil",
		End:      "nil",
		NCands:   0,
	}
	sVoting3 := sv.New(&cfg3, nil)
	hash := voting.GenerateKeyHash(cfg.UserID, cfg.VotingID)
	usrVrfKeyMap := map[string](ed25519.VerfKey){hash: *signKeyPair.Verf()}
	fmt.Println("a")
	vm := sVoting3.Get(nil, usrVrfKeyMap, *signKeyPair2.Verf())
	vi := sVoting3.Count(vm, *eciesKeyPair.Private())
	fmt.Println(vi)
	fmt.Println(time.Now().String())

	bpuk := cfg.PubKey.Marshal()
	puk2 := ecies.UnmarshalPublic(bpuk)
	fmt.Println(puk2.Equals(cfg.PubKey))
	bprk := eciesKeyPair.Private().Marshal()
	prk2 := ecies.UnmarshalPrivate(bprk)
	fmt.Println(prk2.Equals(eciesKeyPair.Private()))

	bsk := cfg.SignKey.Marshal()
	sk2 := ed25519.UnmarshalSign(bsk)
	fmt.Println(sk2.Equals(cfg.SignKey))
	bvk := signKeyPair.Verf().Marshal()
	vk2 := ed25519.UnmarshalVerf(bvk)
	fmt.Println(vk2.Equals(signKeyPair.Verf()))

	fmt.Println(time.Now().String())

	eciesKeyPair0 := ecies.GenKeyPair()
	fmt.Println(eciesKeyPair0.Public().Encrypt(cfg.SignKey.Marshal()))
	fmt.Println(eciesKeyPair0.Public().Encrypt(eciesKeyPair.Private().Marshal()))
	fmt.Println(time.Now().String())

	/*
		sk := ipfs.GenKeyFile()
		fmt.Println(sk)
		resolved := Ipfs.FileAdd([]byte("=^.^= meow meow"), true)
		fmt.Println(resolved.String())
		ipnsEntry := Ipfs.NamePublishWithKeyFile(resolved, "240h", sk, "test-test-test")
		fmt.Println(ipnsEntry.Name())
		ipnsEntry = Ipfs.NamePublishWithKeyFile(resolved, "240h", sk, "test-test-test")
		fmt.Println(ipnsEntry.Name())
		fmt.Println(Ipfs.NameGet(sk))

		pth := Ipfs2.NameResolve(ipnsEntry.Name())
		fmt.Println(pth.String())

		f := Ipfs2.FileGet(pth)
		fmt.Println(string(f))
	*/
}
