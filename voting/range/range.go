package rangevoting

import (
	"EasyVoting/util/ecies"
	"EasyVoting/util/ed25519"
	"EasyVoting/voting"
)

type RangeVoting struct {
	voting.Voting
	min int
	max int
}

func New(cfg *voting.InitConfig, ipnsAddrs []string, min int, max int) *RangeVoting {
	if ok := voting.VerifyUserID(cfg.KeyFile, ipnsAddrs); !ok {
		util.CheckError(errors.New("Invalid KeyFile"))
		return nil
	}

	rv := &RangeVoting{
		min,
		max,
	}
	rv.Init(cfg)
	return rv
}

func (rv *RangeVoting) GenDefaultVoteInt() VoteInt {
	vi := make(VoteInt)
	for _, name := range rv.CandNames() {
		vi[name] = rv.min
	}
	return vi
}

func (rv *RangeVoting) IsValidData(vi voting.VoteInt) bool {
	if !rv.IsCandsMatch(vi) {
		return false
	}

	for _, vote := range vi {
		if vote < rv.min || vote > rv.max {
			return false
		}
	}

	return true
}

func (rv *RangeVoting) Type() string {
	return "rangevoting"
}

func (rv *RangeVoting) Vote(userID string, signKey ed25519.SignKey, data voting.VoteInt) {
	if rv.WithinTime() && rv.IsValidData(data) {
		rv.BaseVote(userID, signKey, data)
	}

}

func (rv *RangeVoting) Count(votes map[string](voting.Vote), manPriKey ecies.PriKey) map[string](voting.VoteInt) {
	votingMap := make(map[string](voting.VoteInt))
	for h, v := range votes {
		data := voting.UnmarshalVoteInt(manPriKey.Decrypt(v.Data))
		if rv.IsValidData(data) {
			votingMap[h] = data
		}
	}

	return votingMap
}
