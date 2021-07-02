package blockvoting

import (
	crsa "crypto/rsa"
	"errors"

	"EasyVoting/util"
	"EasyVoting/voting"
)

type BlockVoting struct {
	voting.Voting
	total int
}

func New(cfg *voting.InitConfig, huidListAddrs []string, total int) *BlockVoting {
	if ok := voting.VerifyUserID(cfg.Is, cfg.UserID, huidListAddrs); !ok {
		util.CheckError(errors.New("Invalid userID"))
		return nil
	}

	bv := &BlockVoting{total: total}
	bv.Init(cfg)
	return bv
}

func (bv *BlockVoting) IsValidData(vb voting.VoteBool) bool {
	if !bv.NumCandsMatch(len(vb)) {
		return false
	}

	numTrue := 0
	for _, vote := range vb {
		if vote {
			numTrue++
		}
	}
	return numTrue >= 0 && numTrue < bv.total
}

func (bv *BlockVoting) Type() string {
	return "blockvoting"
}

func (bv *BlockVoting) Vote(data voting.VoteInt, pubKey crsa.PublicKey) string {
	vb := data.Cast2Bool()
	if bv.WithinTime() && bv.IsValidData(vb) {
		mvb := vb.Marshal()
		return bv.BaseVote(mvb, pubKey)
	}

	util.CheckError(errors.New("Invalid Data"))
	return ""
}

func (bv *BlockVoting) Get(ipnsName string, priKey crsa.PrivateKey) voting.VoteBool {
	mvb := bv.BaseGet(ipnsName, priKey)
	return voting.UnmarshalVoteBool(mvb)
}

func (bv *BlockVoting) Count(nameList map[string]struct{}, proKey crsa.PrivateKey) {

}
