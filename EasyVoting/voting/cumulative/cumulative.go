package cumulativevoting

import(
	"errors"
	crsa "crypto/rsa"

	"EasyVoting/voting"
	"EasyVoting/util"

)



type CumulativeVoting struct{
	voting.Voting
	min int
	total int
}
func New(cfg *voting.InitConfig, huidListAddrs []string, min int, total int) *CumulativeVoting{
	if ok := voting.VerifyUserID(cfg.Is, cfg.UserID, huidListAddrs); !ok{
		util.CheckError(errors.New("Invalid userID"))
		return nil
	}

	cv := &CumulativeVoting{min: min, total: total}
	cv.Init(cfg)
	return cv
}

func (cv *CumulativeVoting) IsValidData(vi voting.VoteInt) bool{
	if !cv.NumCandsMatch(len(vi)){
		return false
	}

	tl := 0
	for _, vote := range vi{
		if vote < cv.min{
			return false
		}
		tl += vote
	}
	return tl <= cv.total
}

func (cv *CumulativeVoting) Type() string{
	return "cumulativevoting"
}

func (cv * CumulativeVoting) Vote(data voting.VoteInt, pubKey crsa.PublicKey) string{
	if cv.WithinTime() && cv.IsValidData(data){
		mvi := data.Marshal()
		return cv.BaseVote(mvi, pubKey)
	}

	util.CheckError(errors.New("Invalid Data"))
	return ""

}

func (cv *CumulativeVoting)Get(ipnsName string, priKey crsa.PrivateKey) voting.VoteInt{
	mvi := cv.BaseGet(ipnsName, priKey)
	return voting.UnmarshalVoteInt(mvi)
}

func (cv *CumulativeVoting) Count(nameList map[string]struct{}, proKey crsa.PrivateKey){
	;
}
