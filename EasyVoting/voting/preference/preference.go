package preferencevoting

import(
	"sort"
	"errors"
	crsa "crypto/rsa"

	"EasyVoting/voting"
	"EasyVoting/util"

)



type PreferenceVoting struct{
	voting.Voting
}
func New(cfg *voting.InitConfig, huidListAddrs []string) *PreferenceVoting{
	if ok := voting.VerifyUserID(cfg.Is, cfg.UserID, huidListAddrs); !ok{
		util.CheckError(errors.New("Invalid userID"))
		return nil
	}

	pv := &PreferenceVoting{}
	pv.Init(cfg)
	return pv
}

func (pv *PreferenceVoting) IsValidData(vi voting.VoteInt) bool{
	if !pv.NumCandsMatch(len(vi)){
		return false
	}

	var ps []int
	for _, vote := range vi{
		ps = append(ps, vote)
	}
	sort.Ints(ps)

	for i, v := range ps{
		if i != v{
			return false
		}
	}
	
	return true
}

func (pv *PreferenceVoting) Type() string{
	return "preferencevoting"
}

func (pv * PreferenceVoting) Vote(data voting.VoteInt, pubKey crsa.PublicKey) string{
	if pv.WithinTime() && pv.IsValidData(data){
		mvi := data.Marshal()
		return pv.BaseVote(mvi, pubKey)
	}

	util.CheckError(errors.New("Invalid Data"))
	return ""

}

func (pv *PreferenceVoting)Get(ipnsName string, priKey crsa.PrivateKey) voting.VoteInt{
	mvi := pv.BaseGet(ipnsName, priKey)
	return voting.UnmarshalVoteInt(mvi)
}

func (pv *PreferenceVoting) Count(nameList map[string]struct{}, proKey crsa.PrivateKey){
	;
}