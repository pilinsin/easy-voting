package voting

import (
	"encoding/json"
	"time"

	crsa "crypto/rsa"
	p2pcrypt "github.com/libp2p/go-libp2p-core/crypto"

	"EasyVoting/ipfs"
	"EasyVoting/util"
	"EasyVoting/util/rsa"
)

func VerifyUserID(kFile p2pcrypt.PrivKey, ipnsAddrs []string) bool {
	ipnsAddr := ipfs.NameGet(kFile)
	for _, addr := range ipnsAddrs {
		if ipnsAddr == addr {
			return true
		}
	}
	return false
}

func generateIPNSKey(userID string, votingID string) string {
	return util.Bytes64EncodeStr(util.Hash([]byte(userID), []byte(votingID)))
}

type Voting struct {
	BaseVoting
	begin   string
	end     string
	keyFile p2pcrypt.PrivKey
	signKey crsa.PrivateKey
}

type InitConfig struct {
	Is        *ipfs.IPFS
	ValidTime string
	VotingID  string
	UserID    string
	Begin     string
	End       string
	NCands    int
	KeyFile   p2pcrypt.PrivKey
	SignKey   crsa.PrivateKey
}

func (v *Voting) Init(cfg *InitConfig) {
	v.BaseInit(cfg.Is, cfg.ValidTime, generateIPNSKey(cfg.UserID, cfg.VotingID), cfg.NCands)
	v.begin = cfg.Begin
	v.end = cfg.End
	v.keyFile = cfg.KeyFile
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

func (v *Voting) BaseVote(data []byte, pubKey crsa.PublicKey) string {
	encData := rsa.Encrypt(data, pubKey)
	resolved := v.iPFS.FileAdd(encData, true)
	ipnsEntry := v.iPFS.NamePublishWithKeyFile(resolved, v.validTime, v.keyFile, v.key)
	return ipnsEntry.Name()
}

func (v *Voting) GenVotingData(data VoteInt) *VotingData {
	sign := rsa.Sign(data.Marshal(), v.signKey)
	vd := &VotingData{data, sign}
	return vd
}

type VotingData struct {
	Data VoteInt
	Sign []byte
}

func (vd *VotingData) Marshal() []byte {
	mvd, err := json.Marshal(vd)
	util.CheckError(err)
	return mvd
}
func UnmarshalVotingData(mvd []byte) VotingData {
	var vd VotingData
	err := json.Unmarshal(mvd, &vd)
	util.CheckError(err)
	return vd
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
