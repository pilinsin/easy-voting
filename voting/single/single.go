package singlevoting

import (
	crsa "crypto/rsa"
	"errors"
	"fmt"

	"EasyVoting/util"
	"EasyVoting/voting"
)

type SingleVoting struct {
	voting.Voting
}

func New(cfg *voting.InitConfig, ipnsAddrs []string) *SingleVoting {
	//if ok := voting.VerifyUserID(cfg.KeyFile, ipnsAddrs); !ok{
	//	util.CheckError(errors.New("Invalid KeyFile"))
	//	return nil
	//}

	sv := &SingleVoting{}
	sv.Init(cfg)
	return sv
}

func (sv *SingleVoting) IsValidData(vi voting.VoteInt) bool {
	if !sv.NumCandsMatch(len(vi)) {
		return false
	}

	numTrue := 0
	for _, vote := range vi {
		if vote > 0 {
			numTrue++
		}
	}
	return numTrue >= 0 && numTrue < 2
}

func (sv *SingleVoting) Type() string {
	return "singlevoting"
}

func (sv *SingleVoting) Vote(data voting.VoteInt, pubKey crsa.PublicKey) string {
	if sv.WithinTime() && sv.IsValidData(data) {
		vd := sv.GenVotingData(data)
		fmt.Println(vd)
		mvd := vd.Marshal()
		return sv.BaseVote(mvd, pubKey)
	}

	util.CheckError(errors.New("Invalid Data"))
	return ""
}

func (sv *SingleVoting) Get(ipnsName string, priKey crsa.PrivateKey) voting.VotingData {
	mvd := sv.BaseGet(ipnsName, priKey)
	return voting.UnmarshalVotingData(mvd)
}

func (sv *SingleVoting) Count(nameList map[string]struct{}, proKey crsa.PrivateKey) {

}
