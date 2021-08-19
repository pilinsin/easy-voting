package main

import (
	"fmt"
	"time"

	"EasyVoting/util"
	"EasyVoting/util/ecies"
	"EasyVoting/util/ed25519"
	"EasyVoting/voting"
	sv "EasyVoting/voting/single"
)

//sudo sysctl -w net.core.rmem_max=2500000
func main() {
	topic := util.GenUniqueID(50, 50)

	eciesKeyPair := ecies.GenKeyPair()
	cfg0 := voting.InitConfig{
		RepoStr:  "ipfs0",
		Topic:    topic,
		VotingID: "nil",
		UserID:   "nil",
		Begin:    "nil",
		End:      "nil",
		NCands:   3,
	}
	sVoting0 := sv.New(&cfg0, nil)

	usrVrfKeyMap := make(map[string](ed25519.VerfKey))
	for itr := 0; itr < 1000; itr++ {
		<-time.After(10 * time.Second)
		signKeyPair := ed25519.GenKeyPair()
		cfg := voting.InitConfig{
			RepoStr:  "ipfs" + string(itr+1),
			Topic:    topic,
			VotingID: util.GenUniqueID(30, 30),
			UserID:   util.GenUniqueID(30, 6),
			Begin:    "2021-6-1  00:00am",
			End:      "2021-8-30 11:59pm",
			NCands:   3,
			PubKey:   eciesKeyPair.Public(),
			SignKey:  signKeyPair.Sign(),
		}
		sVoting := sv.New(&cfg, nil)

		v := voting.VoteInt{"A": 0, "B": 1, "C": 0}
		sVoting.Vote(v)
		v2 := voting.VoteInt{"A": 0, "B": 0, "C": 1}
		sVoting.Vote(v2)

		hash := voting.GenerateKeyHash(cfg.UserID, cfg.VotingID)
		usrVrfKeyMap[hash] = *signKeyPair.Verf()

		sVoting.Close()
	}

	fmt.Println("a")
	vm := sVoting0.Get(nil, usrVrfKeyMap)
	fmt.Println(vm)
	vi := sVoting0.Count(vm, *eciesKeyPair.Private())
	fmt.Println(vi)

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
