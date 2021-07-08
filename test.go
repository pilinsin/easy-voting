package main

import (
	"fmt"

	"EasyVoting/ipfs"
	"EasyVoting/util"
	"EasyVoting/util/rsa"
	"EasyVoting/voting"
	sv "EasyVoting/voting/single"
)

//sudo sysctl -w net.core.rmem_max=2500000
func main() {

	ctx := util.NewContext()
	Ipfs := ipfs.New(ctx, ".ipfs")
	ctx2 := util.NewContext()
	Ipfs2 := ipfs.New(ctx2, ".ipfs2")

	signKeyPair := rsa.GenKeyPair(1024)
	cfg := voting.InitConfig{
		Is:        Ipfs,
		ValidTime: "240h",
		VotingID:  util.GenUniqueID(30, 30),
		UserID:    util.GenUniqueID(30, 6),
		Begin:     "2021-6-1  00:00am",
		End:       "2021-8-13 11:59pm",
		NCands:    3,
		KeyFile:   ipfs.GenKeyFile(),
		SignKey:   signKeyPair.Private(),
	}
	sVoting := sv.New(&cfg, nil)
	v := voting.VoteInt{"A": 0, "B": 1, "C": 0}
	rsaKeyPair := rsa.GenKeyPair(4096)
	res := sVoting.Vote(v, rsaKeyPair.Public())
	fmt.Println(res)

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
	vd := sVoting2.Get(res, rsaKeyPair.Private())
	fmt.Println(vd)

	bk := ipfs.MarshalKeyFile(cfg.KeyFile)
	kf2 := ipfs.UnmarshalKeyFile(bk)
	fmt.Println(kf2.Equals(cfg.KeyFile))
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
