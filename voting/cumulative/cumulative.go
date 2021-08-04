package cumulativevoting

import (
	"EasyVoting/util/ecies"
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

	cv := &CumulativeVoting{min: min, total: total}
	cv.Init(cfg)
	return cv
}

func (cv *CumulativeVoting) IsValidData(vi voting.VoteInt) bool {
	if !cv.NumCandsMatch(len(vi)) {
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

func (cv *CumulativeVoting) Vote(data voting.VoteInt) string {
	if cv.WithinTime() && cv.IsValidData(data) {
		vd := cv.GenVotingData(data)
		mvd := vd.Marshal()
		return cv.BaseVote(mvd)
	}

	util.CheckError(errors.New("Invalid Data"))
	return ""

}

func (cv *CumulativeVoting) Count(votes *voting.VoteMap, manPriKey ecies.PriKey) map[string](voting.VoteInt) {
	var votingMap map[string](voting.VoteInt)
	for h, v := range votes.Votes {
		data := voting.UnmarshalVoteInt(manPriKey.Decrypt(v.Data))
		if cv.IsValidData(data) {
			votingMap[h] = data
		}
	}

	return votingMap
}
