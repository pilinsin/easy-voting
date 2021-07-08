package approvalvoting

import (
	crsa "crypto/rsa"
	"errors"

	"EasyVoting/util"
	"EasyVoting/voting"
)

type ApprovalVoting struct {
	voting.Voting
	nCands int
}

func New(cfg *voting.InitConfig, ipnsAddrs []string) *ApprovalVoting {
	if ok := voting.VerifyUserID(cfg.KeyFile, ipnsAddrs); !ok {
		util.CheckError(errors.New("Invalid KeyFile"))
		return nil
	}

	av := &ApprovalVoting{}
	av.Init(cfg)
	return av
}

func (av *ApprovalVoting) IsValidData(vi voting.VoteInt) bool {
	return av.NumCandsMatch(len(vi))
}

func (av *ApprovalVoting) Type() string {
	return "approvalvoting"
}

func (av *ApprovalVoting) Vote(data voting.VoteInt, pubKey crsa.PublicKey) string {
	if av.WithinTime() && av.IsValidData(data) {
		vd := av.GenVotingData(data)
		mvd := vd.Marshal()
		return av.BaseVote(mvd, pubKey)
	}

	util.CheckError(errors.New("Invalid Data"))
	return ""

}

func (av *ApprovalVoting) Get(ipnsName string, priKey crsa.PrivateKey) voting.VotingData {
	mvd := av.BaseGet(ipnsName, priKey)
	return voting.UnmarshalVotingData(mvd)
}

func (av *ApprovalVoting) Count(nameList map[string]struct{}, proKey crsa.PrivateKey) {

}
