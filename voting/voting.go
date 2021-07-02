package voting

import (
	"encoding/json"
	"time"

	"EasyVoting/ipfs"
	"EasyVoting/util"
)

func generateIPNSKey(votingID string, userID string) string {
	return util.Hash(votingID + userID)
}

type Voting struct {
	BaseVoting
	begin string
	end   string
}

type InitConfig struct {
	Is        *ipfs.IPFS
	ValidTime string
	VotingID  string
	UserID    string
	Begin     string
	End       string
	NCands    int
}

func (v *Voting) Init(cfg *InitConfig) {
	v.BaseInit(cfg.Is, cfg.ValidTime, generateIPNSKey(cfg.VotingID, cfg.UserID), cfg.NCands)
	v.begin = cfg.Begin
	v.end = cfg.End
}

func (v *Voting) WithinTime() bool {
	layout := "2006-1-2 3:04pm"
	bTime, err := time.Parse(layout, v.begin)
	util.CheckError(err)
	eTime, err := time.Parse(layout, v.end)
	util.CheckError(err)
	now := time.Now()
	//TODO: timezone

	return (now.Equal(bTime) || now.After(bTime)) && now.Before(eTime)
}

type VoteInt map[string]int

func (vi VoteInt) Cast2Bool() VoteBool {
	vb := VoteBool{}
	for n, v := range vi {
		vb[n] = v > 0
	}
	return vb
}

func (vi VoteInt) Marshal() []byte {
	mvi, err := json.Marshal(vi)
	util.CheckError(err)
	return mvi
}

//TODO: exclude null data
func UnmarshalVoteInt(mvi []byte) VoteInt {
	var vi VoteInt
	err := json.Unmarshal(mvi, &vi)
	util.CheckError(err)
	return vi
}

type VoteBool map[string]bool

func (vb VoteBool) Marshal() []byte {
	mvb, err := json.Marshal(vb)
	util.CheckError(err)
	return mvb
}

//TODO: exclude null data
func UnmarshalVoteBool(mvb []byte) VoteBool {
	var vb VoteBool
	err := json.Unmarshal(mvb, &vb)
	util.CheckError(err)
	return vb
}
