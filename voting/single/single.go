package singlevoting

import (
	"EasyVoting/util/ecies"
	"EasyVoting/util/ed25519"
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

func (sv *SingleVoting) GenDefaultVoteInt() voting.VoteInt {
	vi := make(voting.VoteInt)
	for _, name := range sv.CandNames() {
		vi[name] = 0
	}
	return vi
}

func (sv *SingleVoting) IsValidData(vi voting.VoteInt) bool {
	if !sv.IsCandsMatch(vi) {
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

func (sv *SingleVoting) Vote(userID string, signKey ed25519.SignKey, data voting.VoteInt) {
	if sv.WithinTime() && sv.IsValidData(data) {
		sv.BaseVote(userID, signKey, data)
	}
}

func (sv *SingleVoting) Count(votes map[string](voting.Vote), manPriKey ecies.PriKey) map[string](voting.VoteInt) {
	votingMap := make(map[string](voting.VoteInt))
	for h, v := range votes {
		data := voting.UnmarshalVoteInt(manPriKey.Decrypt(v.Data))
		if sv.IsValidData(data) {
			votingMap[h] = data
		}
	}

	return votingMap
}
