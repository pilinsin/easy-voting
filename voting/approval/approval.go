package approvalvoting

import (
	"EasyVoting/util/ecies"
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

	av := &ApprovalVoting{nCands}
	av.Init(cfg)
	return av
}

func (av *ApprovalVoting) IsValidData(vi voting.VoteInt) bool {
	return av.NumCandsMatch(len(vi))
}

func (av *ApprovalVoting) Type() string {
	return "approvalvoting"
}

func (av *ApprovalVoting) Vote(data voting.VoteInt) string {
	if av.WithinTime() && av.IsValidData(data) {
		vd := av.GenVotingData(data)
		mvd := vd.Marshal()
		return av.BaseVote(mvd)
	}

	util.CheckError(errors.New("Invalid Data"))
	return ""

}

func (av *ApprovalVoting) Count(votes map[string](voting.Vote), manPriKey ecies.PriKey) map[string](voting.VoteInt) {
	votingMap := make(map[string](voting.VoteInt))
	for h, v := range votes {
		data := voting.UnmarshalVoteInt(manPriKey.Decrypt(v.Data))
		if av.IsValidData(data) {
			votingMap[h] = data
		}
	}

	return votingMap
}
