package voting

import (
	"time"

	"EasyVoting/ipfs"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
	"EasyVoting/util/ecies"
	vutil "EasyVoting/voting/util"
)

type manager struct {
	is            *ipfs.IPFS
	tInfo         *util.TimeInfo
	salt1         string
	salt2         string
	chmCid        string
	ivmCid        string
	manPriKey     *ecies.PriKey
	resMapKeyFile *ipfs.KeyFile
}

func NewManager(mCfgCid string, is *ipfs.IPFS) (*manager, error) {
	mCfg, err := vutil.ManConfigFromCid(mCfgCid, is)
	if err != nil {
		return nil, util.NewError("invalid mCfgCid")
	}

	man := &manager{
		is:            is,
		tInfo:         mCfg.TimeInfo(),
		salt1:         mCfg.Salt1(),
		salt2:         mCfg.Salt2(),
		chmCid:        mCfg.UchmCid(),
		ivmCid:        mCfg.UivmCid(),
		manPriKey:     mCfg.Private(),
		resMapKeyFile: mCfg.KeyFile(),
	}
	return man, nil
}
func (m *manager) Close() {
	m.is.Close()
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
	ivm, err := vutil.IdVotingMapFromCid(m.ivmCid, m.is)
	if err != nil {
		return false, util.NewError("invalid ivmCid")
	}

	name, err := m.resMapKeyFile.Name()
	if err != nil {
		return false, util.NewError("invalid resMapKeyFile")
	}
	resMap, err := vutil.ResultMapFromName(name, m.is)
	if err != nil {
		return false, util.NewError("resMap does not exist")
	}

	return resMap.VerifyVotes(ivm, m.is), nil
}