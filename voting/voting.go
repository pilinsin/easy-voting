package voting

import (
	"encoding/json"
	"fmt"
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
	topic   string
	key     string
	begin   string
	end     string
	nCands  int
	pubKey  *ecies.PubKey
	signKey *ed25519.SignKey
}

type InitConfig struct {
	Is       *ipfs.IPFS
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
	v.iPFS = cfg.Is
	v.topic = cfg.Topic
	v.key = GenerateKeyHash(cfg.UserID, cfg.VotingID)
	v.begin = cfg.Begin
	v.end = cfg.End
	v.nCands = cfg.NCands
	v.pubKey = cfg.PubKey
	v.signKey = cfg.SignKey
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

type VoteMap struct {
	Counts int
	Votes  map[string]Vote
}

func (vm *VoteMap) Append(vote Vote) {
	vm.Counts++
	vm.Votes[vote.Hash] = vote
}
func (vm *VoteMap) Copy(vm2 *VoteMap) {
	vm.Counts = vm2.Counts
	vm.Votes = vm2.Votes
}
func (v *Voting) Get(vts *VoteMap, usrVerfKeyMap map[string](ed25519.VerfKey), manVerfKey ed25519.VerfKey) *VoteMap {
	var votes *VoteMap
	if vts != nil {
		fmt.Println("vts is nil")
		votes.Copy(vts)
	}

	sub := v.iPFS.PubsubSubscribe(v.topic)
	defer sub.Close()
	subCount := 0
	for {
		msg, err := v.iPFS.SubNext(sub)
		fmt.Println(err != nil)
		//if err == io.EOF || err == context.Canceled {
		if err != nil {
			util.CheckError(err)
			return votes
		}
		if subCount < votes.Counts {
			subCount++
			continue
		}

		data := msg.Data()
		if IsVoteEnd(data, manVerfKey) {
			return votes
		}

		vd, err := UnmarshalVotingData(data)
		if err != nil {
			if _, ok := usrVerfKeyMap[vd.Vote.Hash]; ok {
				if vd.Verify(usrVerfKeyMap[vd.Vote.Hash]) {
					votes.Append(vd.Vote)
				}
			}
		}

	}
}

const voteEnd = "voting end"

type VoteEnd struct {
	Text string
	Sign []byte
}

func (v *Voting) MarshalVoteEnd() []byte {
	sign := v.signKey.Sign([]byte(voteEnd))
	ve := VoteEnd{voteEnd, sign}
	mve, err := json.Marshal(ve)
	util.CheckError(err)
	return mve
}
func IsVoteEnd(mve []byte, verfKey ed25519.VerfKey) bool {
	var ve VoteEnd
	err := json.Unmarshal(mve, &ve)
	if err != nil {
		return false
	}
	return verfKey.Verify([]byte(voteEnd), ve.Sign)
}

func (v *Voting) GenVotingData(data VoteInt) VotingData {
	vt := Vote{v.key, v.pubKey.Encrypt(data.Marshal())}
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
