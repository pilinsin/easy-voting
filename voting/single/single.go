package singlevoting

import (
	"EasyVoting/util/ecies"
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

func (sv *SingleVoting) Vote(data voting.VoteInt) {
	if sv.WithinTime() && sv.IsValidData(data) {
		vd := sv.GenVotingData(data)
		mvd := vd.Marshal()
		sv.BaseVote(mvd)
	}
}

func (sv *SingleVoting) Count(votes *voting.VoteMap, manPriKey ecies.PriKey) map[string](voting.VoteInt) {
	var votingMap map[string](voting.VoteInt)
	for h, v := range votes.Votes {
		data := voting.UnmarshalVoteInt(manPriKey.Decrypt(v.Data))
		if sv.IsValidData(data) {
			votingMap[h] = data
		}
	}

	return votingMap
}
