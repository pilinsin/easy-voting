package blockvoting

import (
	"EasyVoting/util/ecies"
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

	bv := &BlockVoting{total: total}
	bv.Init(cfg)
	return bv
}

func (bv *BlockVoting) IsValidData(vi voting.VoteInt) bool {
	if !bv.NumCandsMatch(len(vi)) {
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

func (bv *BlockVoting) Vote(data voting.VoteInt) string {
	if bv.WithinTime() && bv.IsValidData(data) {
		vd := bv.GenVotingData(data)
		mvd := vd.Marshal()
		return bv.BaseVote(mvd)
	}

	util.CheckError(errors.New("Invalid Data"))
	return ""
}

func (bv *BlockVoting) Count(votes *voting.VoteMap, manPriKey ecies.PriKey) map[string](voting.VoteInt) {
	var votingMap map[string](voting.VoteInt)
	for h, v := range votes.Votes {
		data := voting.UnmarshalVoteInt(manPriKey.Decrypt(v.Data))
		if bv.IsValidData(data) {
			votingMap[h] = data
		}
	}

	return votingMap
}
