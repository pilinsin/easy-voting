package voting

import (
	"fmt"
	"time"

	iface "github.com/ipfs/interface-go-ipfs-core"

	"github.com/pilinsin/easy-voting/ipfs"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	"github.com/pilinsin/easy-voting/util"
	"github.com/pilinsin/easy-voting/util/crypto"
	vutil "github.com/pilinsin/easy-voting/voting/util"
)

type manager struct {
	is            *ipfs.IPFS
	sub         iface.PubSubSubscription
	tInfo         *util.TimeInfo
	salt1         string
	salt2         string
	chmCid        string
	ivmCid        string
	manPriKey     crypto.IPriKey
	verfMapKeyFile *ipfs.KeyFile
	resMapKeyFile *ipfs.KeyFile
}

func NewManager(vCfgCid string, manIdentity *vutil.ManIdentity, is *ipfs.IPFS) (*manager, error) {
	vCfg, err := vutil.ConfigFromCid(vCfgCid, is)
	if err != nil {
		return nil, util.NewError("invalid vCfgCid")
	}

	man := &manager{
		is:            is,
		sub:         is.PubSubSubscribe("voting_pubsub/" + vCfgCid),
		tInfo:         vCfg.TimeInfo(),
		salt1:         vCfg.Salt1(),
		salt2:         vCfg.Salt2(),
		chmCid:        vCfg.UchmCid(),
		ivmCid:        vCfg.UivmCid(),
		manPriKey:     manIdentity.Private(),
		verfMapKeyFile: manIdentity.VerfMapKeyFile(),
		resMapKeyFile: manIdentity.ResMapKeyFile(),
	}
	return man, nil
}
func (m *manager) Close() {
	m.sub.Close()
	m.sub = nil
	m.is = nil
}

func (m manager) IsValidUser(userData ...string) bool {
	chm := &rutil.ConstHashMap{}
	err := chm.FromCid(m.chmCid, m.is)
	if err != nil {
		return false
	}

	userHash := rutil.NewUserHash(m.is, m.salt1, userData...)
	uhHash := rutil.NewUhHash(m.is, m.salt2, userHash)
	return chm.ContainHash(uhHash, m.is)
}

func (m *manager) Registrate() error {
	ivm, err := vutil.IdVotingMapFromCid(m.ivmCid, m.is)
	if err != nil {
		return err
	}
	verfName, err := m.verfMapKeyFile.Name()
	if err != nil {
		return util.NewError("invalid verfMapKeyFile")
	}
	verfMap, err := vutil.IdVerfKeyMapFromName(verfName, m.is)
	if err != nil {
		return util.NewError("resMap does not exist")
	}

	//it takes 5~6 mins
	subs := m.is.PubSubNextAll(m.sub)
	fmt.Println("data list: ", subs)
	if len(subs) <= 0 {
		return nil
	}
	isVerfMapUpdated := false
	for _, encUInfo := range subs {
		mUInfo, err := m.manPriKey.Decrypt(encUInfo)
		if err != nil {
			fmt.Println("decrypt uInfo error")
			continue
		}
		uInfo := &vutil.UserInfo{}
		if err := uInfo.Unmarshal(mUInfo); err != nil {
			fmt.Println("uInfo unmarshal error")
			continue
		}
		if err := verfMap.Append(uInfo.UvHash(), uInfo.Verify(), ivm, m.manPriKey, m.is); err == nil{
			isVerfMapUpdated = true
			fmt.Println("uInfo appended")
		}
	}
	if isVerfMapUpdated{
		name := ipfs.ToNameWithKeyFile(verfMap.Marshal(), m.verfMapKeyFile, m.is)
		fmt.Println("ipnsPublished to ", name)
	}
	return nil
}

func (m manager) GetResultMap() error {
	if ok := m.tInfo.AfterTime(time.Now()); !ok {
		return util.NewError("now is the voting time")
	}

	ivm, err := vutil.IdVotingMapFromCid(m.ivmCid, m.is)
	if err != nil {
		return err
	}
	resMap, err := vutil.NewResultMap(100000, ivm, m.manPriKey, m.is)
	if err != nil {
		return err
	}

	ipfs.ToNameWithKeyFile(resMap.Marshal(), m.resMapKeyFile, m.is)
	return nil
}

func (m manager) VerifyResultMap() (bool, error) {
	verfName, err := m.verfMapKeyFile.Name()
	if err != nil {
		return false, util.NewError("invalid verfMapKeyFile")
	}
	verfMap, err := vutil.IdVerfKeyMapFromName(verfName, m.is)
	if err != nil {
		return false, util.NewError("resMap does not exist")
	}
	resName, err := m.resMapKeyFile.Name()
	if err != nil {
		return false, util.NewError("invalid resMapKeyFile")
	}
	resMap, err := vutil.ResultMapFromName(resName, m.is)
	if err != nil {
		return false, util.NewError("resMap does not exist")
	}

	return resMap.VerifyVotes(verfMap, m.is), nil
}
