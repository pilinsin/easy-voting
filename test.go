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
	ctx := util.NewContext()
	Ipfs := ipfs.New(ctx, ".ipfs")
	ctx2 := util.NewContext()
	Ipfs2 := ipfs.New(ctx2, ".ipfs2")
	fmt.Println("ipfs initialized.")
	fmt.Println(time.Now().String())

	Ipfs.PubsubPublish("topic", []byte("=^o^= meow meow"))
	fmt.Println(Ipfs.PubsubSubscribe("topic"))

	signKeyPair := ed25519.GenKeyPair()
	eccKeyPair := ecies.GenKeyPair()
	cfg := voting.InitConfig{
		Is:        Ipfs,
		ValidTime: "240h",
		VotingID:  util.GenUniqueID(30, 30),
		UserID:    util.GenUniqueID(30, 6),
		Begin:     "2021-6-1  00:00am",
		End:       "2021-8-13 11:59pm",
		NCands:    3,
		PubKey:    eccKeyPair.Public(),
		KeyFile:   ipfs.GenKeyFile(),
		SignKey:   signKeyPair.Sign(),
	}
	sVoting := sv.New(&cfg, nil)
	v := voting.VoteInt{"A": 0, "B": 1, "C": 0}
	res := sVoting.Vote(v)
	fmt.Println(res)
	fmt.Println(time.Now().String())

	cfg2 := voting.InitConfig{
		Is:        Ipfs2,
		ValidTime: "nil",
		VotingID:  "nil",
		UserID:    "nil",
		Begin:     "nil",
		End:       "nil",
		NCands:    0,
	}
	sVoting2 := sv.New(&cfg2, nil)
	vd := sVoting2.Get(res, eccKeyPair.Private())
	fmt.Println(vd)
	fmt.Println(vd.Verify(signKeyPair.Verf()))
	fmt.Println(time.Now().String())

	bk := cfg.KeyFile.Marshal()
	kf2 := ipfs.UnmarshalKeyFile(bk)
	fmt.Println(kf2.Equals(cfg.KeyFile))

	bpuk := cfg.PubKey.Marshal()
	puk2 := ecies.UnmarshalPublic(bpuk)
	fmt.Println(puk2.Equals(cfg.PubKey))
	bprk := eccKeyPair.Private().Marshal()
	prk2 := ecies.UnmarshalPrivate(bprk)
	fmt.Println(prk2.Equals(eccKeyPair.Private()))

	bsk := cfg.SignKey.Marshal()
	sk2 := ed25519.UnmarshalSign(bsk)
	fmt.Println(sk2.Equals(cfg.SignKey))
	bvk := signKeyPair.Verf().Marshal()
	vk2 := ed25519.UnmarshalVerf(bvk)
	fmt.Println(vk2.Equals(signKeyPair.Verf()))

	fmt.Println(time.Now().String())

	eccKeyPair0 := ecies.GenKeyPair()
	fmt.Println(eccKeyPair0.Public().Encrypt(cfg.KeyFile.Marshal()))
	fmt.Println(eccKeyPair0.Public().Encrypt(cfg.SignKey.Marshal()))
	fmt.Println(eccKeyPair0.Public().Encrypt(eccKeyPair.Private().Marshal()))
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
