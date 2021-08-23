package preferencevoting

import (
	"EasyVoting/util/ecies"
	"EasyVoting/util/ed25519"
	"EasyVoting/voting"
)

type PreferenceVoting struct {
	voting.Voting
}

func New(cfg *voting.InitConfig, ipnsAddrs []string) *PreferenceVoting {
	if ok := voting.VerifyUserID(cfg.KeyFile, ipnsAddrs); !ok {
		util.CheckError(errors.New("Invalid KeyFile"))
		return nil
	}

	pv := &PreferenceVoting{}
	pv.Init(cfg)
	return pv
}

func (pv *PreferenceVoting) GenDefaultVoteInt() VoteInt {
	vi := make(VoteInt)
	for _, name := range pv.CandNames() {
		vi[name] = 0
	}
	return vi
}

func (pv *PreferenceVoting) IsValidData(vi voting.VoteInt) bool {
	if !pv.IsCandsMatch(vi) {
		return false
	}

	var ps []int
	for _, vote := range vi {
		ps = append(ps, vote)
	}

	isDefaultVote := true
	for _, p := range ps {
		if p != 0 {
			isDefaultVote = false
			break
		}
	}
	if isDefaultVote {
		return true
	}

	sort.Ints(ps)
	for i, v := range ps {
		if i != v {
			return false
		}
	}

	return true
}

func (pv *PreferenceVoting) Type() string {
	return "preferencevoting"
}

func (pv *PreferenceVoting) Vote(userID string, signKey ed25519.SignKey, data voting.VoteInt) {
	if pv.WithinTime() && pv.IsValidData(data) {
		pv.BaseVote(userID, signKey, data)
	}

}

func (pv *PreferenceVoting) Count(votes map[string](voting.Vote), manPriKey ecies.PriKey) map[string](voting.VoteInt) {
	votingMap := make(map[string](voting.VoteInt))
	for h, v := range votes {
		data := voting.UnmarshalVoteInt(manPriKey.Decrypt(v.Data))
		if pv.IsValidData(data) {
			votingMap[h] = data
		}
	}

	return votingMap
}
