package cumulativevoting

import (
	"EasyVoting/util/ecies"
	"EasyVoting/util/ed25519"
	"EasyVoting/voting"
)

type CumulativeVoting struct {
	voting.Voting
	min   int
	total int
}

func New(cfg *voting.InitConfig, ipnsAddrs []string, min int, total int) *CumulativeVoting {
	if ok := voting.VerifyUserID(cfg.KeyFile, ipnsAddrs); !ok {
		util.CheckError(errors.New("Invalid KeyFile"))
		return nil
	}

	cv := &CumulativeVoting{
		min,
		total,
	}
	cv.Init(cfg)
	return cv
}

func (cv *CumulativeVoting) GenDefaultVoteInt() VoteInt {
	vi := make(VoteInt)
	for _, name := range cv.CandNames() {
		vi[name] = cv.min
	}
	return vi
}

func (cv *CumulativeVoting) IsValidData(vi voting.VoteInt) bool {
	if !cv.IsCandsMatch(vi) {
		return false
	}

	tl := 0
	for _, vote := range vi {
		if vote < cv.min {
			return false
		}
		tl += vote
	}
	return tl <= cv.total
}

func (cv *CumulativeVoting) Type() string {
	return "cumulativevoting"
}

func (cv *CumulativeVoting) Vote(userID string, signKey ed25519.SignKey, data voting.VoteInt) {
	if cv.WithinTime() && cv.IsValidData(data) {
		cv.BaseVote(userID, signKey, data)
	}

}

func (cv *CumulativeVoting) Count(votes map[string](voting.Vote), manPriKey ecies.PriKey) map[string](voting.VoteInt) {
	votingMap := make(map[string](voting.VoteInt))
	for h, v := range votes {
		data := voting.UnmarshalVoteInt(manPriKey.Decrypt(v.Data))
		if cv.IsValidData(data) {
			votingMap[h] = data
		}
	}

	return votingMap
}
