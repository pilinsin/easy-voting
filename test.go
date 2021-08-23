package main

import (
	"fmt"
	"strconv"
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
	vid := util.GenUniqueID(30, 30)
	begin := "2021-6-1  00:00am"
	end := "2021-8-30 11:59pm"
	cNames := []string{"A", "B", "C"}

	eciesKeyPair := ecies.GenKeyPair()
	cfg0 := voting.InitConfig{
		RepoStr:   "ipfs_man",
		Topic:     topic,
		VotingID:  vid,
		Begin:     begin,
		End:       end,
		CandNames: cNames,
		PubKey:    eciesKeyPair.Public(),
	}
	sVoting0 := sv.New(&cfg0, nil)

	uids := make([]string, 0)
	usrSignKeys := make([]ed25519.SignKey, 0)
	usrVrfKeyMap := make(map[string](ed25519.VerfKey))
	users := make([]*sv.SingleVoting, 0)
	for itr := 0; itr < 3; itr++ {
		signKeyPair := ed25519.GenKeyPair()
		uid := util.GenUniqueID(30, 6)
		cfg := voting.InitConfig{
			RepoStr:   "ipfs_usr" + strconv.Itoa(itr),
			Topic:     topic,
			VotingID:  vid,
			Begin:     begin,
			End:       end,
			CandNames: cNames,
			PubKey:    eciesKeyPair.Public(),
		}
		users = append(users, sv.New(&cfg, nil))

		uids = append(uids, uid)
		usrSignKeys = append(usrSignKeys, *signKeyPair.Sign())
		hash := voting.GenerateKeyHash(uid, cfg.VotingID)
		usrVrfKeyMap[hash] = *signKeyPair.Verf()
	}

	defaultVote := sVoting0.GenDefaultVoteInt()
	for idx, _ := range uids {
		uid := uids[idx]
		signKey := usrSignKeys[idx]
		sVoting0.Vote(uid, signKey, defaultVote)
	}
	vm := sVoting0.Get(nil, usrVrfKeyMap)
	vi := sVoting0.Count(vm, *eciesKeyPair.Private())
	fmt.Println(vi)

	time.After(5 * time.Second)

	v := voting.VoteInt{"A": 0, "B": 1, "C": 0}
	for idx, user := range users {
		uid := uids[idx]
		signKey := usrSignKeys[idx]
		user.Vote(uid, signKey, v)
	}

	fmt.Println("a")
	vm2 := sVoting0.Get(nil, usrVrfKeyMap)
	vi2 := sVoting0.Count(vm2, *eciesKeyPair.Private())
	fmt.Println(vi2)

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
