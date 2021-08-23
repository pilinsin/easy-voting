package blockvoting

import (
	"EasyVoting/util/ecies"
	"EasyVoting/util/ed25519"
	"EasyVoting/voting"
)

type BlockVoting struct {
	voting.Voting
	total int
}

func New(cfg *voting.InitConfig, ipnsAddrs []string, total int) *BlockVoting {
	if ok := voting.VerifyUserID(cfg.KeyFile, ipnsAddrs); !ok {
		util.CheckError(errors.New("Invalid KeyFile"))
		return nil
	}

	bv := &BlockVoting{total}
	bv.Init(cfg)
	return bv
}

func (bv *BlockVoting) GenDefaultVoteInt() VoteInt {
	vi := make(VoteInt)
	for _, name := range bv.CandNames() {
		vi[name] = 0
	}
	return vi
}

func (bv *BlockVoting) IsValidData(vi voting.VoteInt) bool {
	if !bv.IsCandsMatch(vi) {
		return false
	}

	numTrue := 0
	for _, vote := range vi {
		if vote > 0 {
			numTrue++
		}
	}
	return numTrue >= 0 && numTrue < bv.total
}

func (bv *BlockVoting) Type() string {
	return "blockvoting"
}

func (bv *BlockVoting) Vote(userID string, signKey ed25519.SignKey, data voting.VoteInt) {
	if bv.WithinTime() && bv.IsValidData(data) {
		bv.BaseVote(userID, signKey, data)
	}

}

func (bv *BlockVoting) Count(votes map[string](voting.Vote), manPriKey ecies.PriKey) map[string](voting.VoteInt) {
	votingMap := make(map[string](voting.VoteInt))
	for h, v := range votes {
		data := voting.UnmarshalVoteInt(manPriKey.Decrypt(v.Data))
		if bv.IsValidData(data) {
			votingMap[h] = data
		}
	}

	return votingMap
}
