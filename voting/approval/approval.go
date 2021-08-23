package approvalvoting

import (
	"EasyVoting/util/ecies"
	"EasyVoting/util/ed25519"
	"EasyVoting/voting"
)

type ApprovalVoting struct {
	voting.Voting
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

func (av *ApprovalVoting) GenDefaultVoteInt() VoteInt {
	vi := make(VoteInt)
	for _, name := range av.CandNames() {
		vi[name] = 0
	}
	return vi
}

func (av *ApprovalVoting) IsValidData(vi voting.VoteInt) bool {
	return av.IsCandsMatch(vi)
}

func (av *ApprovalVoting) Type() string {
	return "approvalvoting"
}

func (av *ApprovalVoting) Vote(userID string, signKey ed25519.SignKey, data voting.VoteInt) {
	if av.WithinTime() && av.IsValidData(data) {
		av.BaseVote(userID, signKey, data)
	}

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
