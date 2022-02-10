package voting

import (
	"fmt"
	"time"

	iface "github.com/ipfs/interface-go-ipfs-core"

	"github.com/pilinsin/ipfs-util"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	vutil "github.com/pilinsin/easy-voting/voting/util"
)

type manager struct {
	is            *ipfs.IPFS
	tInfo         *util.TimeInfo
	salt1         string
	salt2         string
	verfSub         iface.PubSubSubscription
	voteSub         iface.PubSubSubscription
	logSub         iface.PubSubSubscription
	logTopic string
	hashVoteMap *vutil.HashVoteMap
	hbmCid string
	uhmCid        string
	manPriKey     crypto.IPriKey
	verfMapKeyFile *ipfs.KeyFile
	resBoxKeyFile *ipfs.KeyFile
}

func NewManager(vCfgCid string, manIdentity *vutil.ManIdentity, is *ipfs.IPFS) (*manager, error) {
	vCfg, err := vutil.ConfigFromCid(vCfgCid, is)
	if err != nil {
		return nil, util.NewError("invalid vCfgCid")
	}

	man := &manager{
		is:            is,
		tInfo:         vCfg.TimeInfo(),
		salt1:         vCfg.Salt1(),
		salt2:         vCfg.Salt2(),
		verfSub:         is.PubSub().Subscribe("verfKey_pubsub/" + vCfgCid),
		voteSub:         is.PubSub().Subscribe("voting_pubsub/" + vCfgCid),
		logSub:         is.PubSub().Subscribe("log_pubsub/" + vCfgCid),
		logTopic: "log_pubsub/" + vCfgCid,
		hashVoteMap: vutil.NewHashVoteMap(100000, vCfg.TimeInfo(), vCfg.VotingID)
		hbmCid: vCfg.HbmCid(),
		uhmCid:        vCfg.UhmCid(),
		manPriKey:     manIdentity.Private(),
		verfMapKeyFile: manIdentity.VerfMapKeyFile(),
		resBoxKeyFile: manIdentity.ResBoxKeyFile(),
	}
	return man, nil
}
func (m *manager) Load() error{
	if err := updateHashVerfMap(); err != nil{return err}
	hvkmName, _ := m.verfMapKeyFile.Name()
	hvkm, err := vutil.HashVerfMapFromName(hvkmName, m.is)

	cids := m.is.PubSub().NextAll(m.logSub)
	for idx, _ := range cids{
		cid := util.Bytes64ToAnyStr(cids[len(cids) - idx - 1])
		hvtm := HashVoteMapFromCid(cid, m.is)
		if ok := hvtm.VerifyMap(hvkm, m.is); ok{
			m.hashVoteMap = hvtm
			return nil
		}
	}
	return util.NewError("load failed: no valid log")
}

func (m *manager) Close() {
	m.verfSub.Close()
	m.verfSub = nil
	m.voteSub.Close()
	m.voteSub = nil
	m.logSub.Close()
	m.logSub = nil
	m.is = nil
}

func (m manager) IsValidUser(userData ...string) bool {
	mhbm, err := ipfs.File.Get(m.hbmCid, m.is)
	if err != nil {
		return false
	}
	hbm, err := UnmarshalUvhHashMap(mhbm, m.is)
	if err != nil {
		return false
	}

	userHash := rutil.NewUserHash(m.is, m.salt1, userData...)
	uhHash := rutil.NewUhHash(m.is, m.salt2, userHash)
	_, ok := hbm.ContainHash(uhHash, m.is)
	return ok
}

func (m *manager) Registrate() error {
	if err := updateHashVerfMap(); err != nil{return err}
	updateHashVoteMap()
	return nil
}
func (m *manager) updateHashVerfMap() error{
	uvhm, err := UvhHashMapFromCid(m.uhmCid, m.is)
	if err != nil{return err}

	hvkmName, _ := m.verfMapKeyFile.Name()
	hvkm, err := vutil.HashVerfMapFromName(hvkmName, m.is)
	if err != nil{return err}
	hvkmUpdate := false
	uInfos := m.is.PubSub().NextAll(m.verfSub)
	for _, encUInfo := range uInfos{
		mui, err := m.manPriKey.Decrypt(encUInfo)
		if err != nil{continue}
		uInfo := &vutil.UserInfo{}
		if err := uInfo.Unmarshal(mui); err != nil{continue}
		if err := hvkm.Append(uInfo, uvhm, m.is); err == nil{hvkmUpdate = true}
	}
	if hvkmUpdate{
		ipfs.Name.PublishWithKeyFile(hvkm.Marshal(), m.verfMapKeyFile, m.is)
	}
	return nil
}
func (m *manager) updateHashVoteMap(){
	vInfos := m.is.PubSub().NextAll(m.voteSub)
	for _, mvi := range vInfos{
		vInfo := &vutil.VoteInfo{}
		if err := vInfo.Unmarshal(mvi); err != nil{continue}
		m.hashVoteMap.Append(vInfo, hvkm, m.is)
	}
}

func (m *manager) Log(){
	cid := ipfs.File.Add(m.hashVoteMap.Marshal(), m.is)
	m.is.PubSub().Publish(util.AnyStrToBytes64(cid), m.logTopic)
}

func (m manager) UploadResultBox() error {
	if ok := m.tInfo.AfterTime(time.Now()); !ok {
		return util.NewError("now is the voting time")
	}
	resBox := vutil.NewResultBox(m.hashVoteMap, m.manPriKey)
	ipfs.Name.PublishWithKeyFile(resBox.Marshal(), m.resBoxKeyFile, m.is)
	return nil
}

func (m manager) VerifyResultBox() (bool, error) {
	hvkmName, _ := m.verfMapKeyFile.Name()
	hvkm, err := vutil.HashVerfMapFromName(hvkmName, m.is)
	if err != nil{return false, err}
	return m.hashVoteMap.VerifyMap(hvkm, m.is), nil
}
