package approvalvoting

import(
	"errors"
	crsa "crypto/rsa"

	"EasyVoting/voting"
	"EasyVoting/util"



)



type ApprovalVoting struct{
	voting.Voting
	nCands int
}
func New(cfg *voting.InitConfig, huidListAddrs []string) *ApprovalVoting{
	if ok := voting.VerifyUserID(cfg.Is, cfg.UserID, huidListAddrs); !ok{
		util.CheckError(errors.New("Invalid userID"))
		return nil
	}

	av := &ApprovalVoting{}
	av.Init(cfg)
	return av
}

func (av *ApprovalVoting) IsValidData(vb voting.VoteBool) bool{
	return av.NumCandsMatch(len(vb))
}

func (av *ApprovalVoting) Type() string{
	return "approvalvoting"
}

func (av * ApprovalVoting) Vote(data voting.VoteInt, pubKey crsa.PublicKey) string{
	vb := data.Cast2Bool()
	if av.WithinTime() && av.IsValidData(vb){
		mvb := vb.Marshal()
		return av.BaseVote(mvb, pubKey)
	}

	util.CheckError(errors.New("Invalid Data"))
	return ""

}

func (av *ApprovalVoting)Get(ipnsName string, priKey crsa.PrivateKey) voting.VoteBool{
	mvb := av.BaseGet(ipnsName, priKey)
	return voting.UnmarshalVoteBool(mvb)
}

func (av *ApprovalVoting) Count(nameList map[string]struct{}, proKey crsa.PrivateKey){
	;
}