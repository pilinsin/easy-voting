package voting

import (
	"encoding/json"
	"time"

	"EasyVoting/ipfs"
	"EasyVoting/util"
	"EasyVoting/util/ecies"
	"EasyVoting/util/ed25519"
)

func VerifyUserID(kFile ipfs.KeyFile, ipnsAddrs []string) bool {
	ipnsAddr := ipfs.NameGet(kFile)
	for _, addr := range ipnsAddrs {
		if ipnsAddr == addr {
			return true
		}
	}
	return false
}

func GenerateKeyHash(userID string, votingID string) string {
	return util.Bytes64EncodeStr(util.Hash([]byte(userID), []byte(votingID)))
}

type Voting struct {
	iPFS    *ipfs.IPFS
	key     string
	begin   string
	end     string
	nCands  int
	pubKey  *ecies.PubKey
	signKey *ed25519.SignKey
}

type InitConfig struct {
	RepoStr  string
	Topic    string
	VotingID string
	UserID   string
	Begin    string
	End      string
	NCands   int
	PubKey   *ecies.PubKey
	SignKey  *ed25519.SignKey
}

func (v *Voting) Init(cfg *InitConfig) {
	v.iPFS = ipfs.New(util.NewContext(), cfg.RepoStr)
	v.key = GenerateKeyHash(cfg.UserID, cfg.VotingID)
	v.begin = cfg.Begin
	v.end = cfg.End
	v.nCands = cfg.NCands
	v.pubKey = cfg.PubKey
	v.signKey = cfg.SignKey

	v.iPFS.PubsubConnect(cfg.Topic)
}
func (v *Voting) Close() {
	v.iPFS.PubsubClose()
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

func (v *Voting) NumCandsMatch(num int) bool {
	return num == v.nCands
}

func (v *Voting) BaseVote(data []byte) {
	v.iPFS.PubsubPublish(data)
}

func (v *Voting) Get(vts map[string]Vote, usrVerfKeyMap map[string](ed25519.VerfKey)) map[string]Vote {
	votes := make(map[string]Vote)
	if vts != nil {
		for k, v := range vts {
			votes[k] = v
		}
	}

	dataset := v.iPFS.PubsubSubscribe()
	if dataset == nil {
		return votes
	}
	for _, data := range dataset {
		vd, err := UnmarshalVotingData(data)
		if err == nil {
			if _, ok := usrVerfKeyMap[vd.Vote.Hash]; ok {
				if vd.Verify(usrVerfKeyMap[vd.Vote.Hash]) {
					if _, ok := votes[vd.Vote.Hash]; ok {
						oldTime := votes[vd.Vote.Hash].Time
						newTime := vd.Vote.Time
						if newTime.After(oldTime) || newTime.Equal(oldTime) {
							votes[vd.Vote.Hash] = vd.Vote
						}
					} else {
						votes[vd.Vote.Hash] = vd.Vote
					}
				}
			}
		}

	}
	return votes
}

func (v *Voting) GenVotingData(data VoteInt) VotingData {
	vt := Vote{v.key, time.Now(), v.pubKey.Encrypt(data.Marshal())}
	sign := v.signKey.Sign(vt.Marshal())
	return VotingData{vt, sign}
}

type VotingData struct {
	Vote Vote
	Sign []byte
}

func (vd *VotingData) Verify(verfKey ed25519.VerfKey) bool {
	return verfKey.Verify(vd.Vote.Marshal(), vd.Sign)
}

func (vd VotingData) Marshal() []byte {
	mvd, err := json.Marshal(vd)
	util.CheckError(err)
	return mvd
}
func UnmarshalVotingData(mvd []byte) (VotingData, error) {
	var vd VotingData
	err := json.Unmarshal(mvd, &vd)
	return vd, err
}

type Vote struct {
	Hash string
	Time time.Time
	Data []byte
}

func (vt Vote) Marshal() []byte {
	mvt, err := json.Marshal(vt)
	util.CheckError(err)
	return mvt
}

//TODO: exclude null data
func UnmarshalVote(mvt []byte) Vote {
	var vt Vote
	err := json.Unmarshal(mvt, &vt)
	util.CheckError(err)
	return vt
}

type VoteInt map[string]int

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
