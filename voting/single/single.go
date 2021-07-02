package singlevoting

import (
	crsa "crypto/rsa"
	"errors"

	"EasyVoting/util"
	"EasyVoting/voting"
)

type SingleVoting struct {
	voting.Voting
}

func New(cfg *voting.InitConfig, huidListAddrs []string) *SingleVoting {
	//if ok := voting.VerifyUserID(cfg.Is, cfg.UserID, huidListAddrs); !ok{
	//	util.CheckError(errors.New("Invalid userID"))
	//	return nil
	//}

	sv := &SingleVoting{}
	sv.Init(cfg)
	return sv
}

func (sv *SingleVoting) IsValidData(vb voting.VoteBool) bool {
	if !sv.NumCandsMatch(len(vb)) {
		return false
	}

	numTrue := 0
	for _, vote := range vb {
		if vote {
			numTrue++
		}
	}
	return numTrue >= 0 && numTrue < 2
}

func (sv *SingleVoting) Type() string {
	return "singlevoting"
}

func (sv *SingleVoting) Vote(data voting.VoteInt, pubKey crsa.PublicKey) string {
	vb := data.Cast2Bool()
	if sv.WithinTime() && sv.IsValidData(vb) {
		mvb := vb.Marshal()
		return sv.BaseVote(mvb, pubKey)
	}

	util.CheckError(errors.New("Invalid Data"))
	return ""
}

func (sv *SingleVoting) Get(ipnsName string, priKey crsa.PrivateKey) voting.VoteBool {
	mvb := sv.BaseGet(ipnsName, priKey)
	return voting.UnmarshalVoteBool(mvb)
}

func (sv *SingleVoting) Count(nameList map[string]struct{}, proKey crsa.PrivateKey) {

}
