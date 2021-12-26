package registration

import (
	"EasyVoting/ipfs"
	"EasyVoting/ipfs/pubsub"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
	"EasyVoting/util/ecies"
)

type IManager interface {
	Close()
	Registrate() error
}

type manager struct {
	is          *ipfs.IPFS
	ps          *pubsub.PubSub
	priKey      *ecies.PriKey
	keyFile     *ipfs.KeyFile
	salt2       string
	chmCid      string
	hnmIpnsName string
}

func NewManager(mCfgCid string, is *ipfs.IPFS) (*manager, error) {
	mCfg, err := rutil.ManConfigFromCid(mCfgCid, is)
	if err != nil {
		return nil, err
	}
	ps, err := pubsub.New("registration_pubsub/" + mCfgCid)
	if err != nil {
		return nil, util.AddError(err, "pubsub.New error")
	}

	man := &manager{
		is:          is,
		ps:          ps,
		priKey:      mCfg.Private(),
		keyFile:     mCfg.KeyFile(),
		salt2:       mCfg.Salt2(),
		chmCid:      mCfg.ChMapCid(),
		hnmIpnsName: mCfg.HnmIpnsName(),
	}
	return man, nil
}
func (m *manager) Close() {
	m.is.Close()
	m.is = nil
	m.ps.Close()
	m.ps = nil
}
func (m *manager) Registrate() error {
	chm := &rutil.ConstHashMap{}
	err := chm.FromCid(m.chmCid, m.is)
	if err != nil {
		return err
	}
	hnm := &rutil.HashNameMap{}
	hnm.FromName(m.hnmIpnsName, m.is)
	if err != nil {
		return err
	}

	for _, encUInfo := range m.ps.Subscribe() {
		mUInfo, err := m.priKey.Decrypt(encUInfo)
		if err != nil {
			continue
		}
		uInfo := &rutil.UserInfo{}
		err = uInfo.Unmarshal(mUInfo)
		if err != nil {
			continue
		}

		uhHash := rutil.NewUhHash(m.is, m.salt2, uInfo.UserHash())
		if ok := chm.ContainHash(uhHash, m.is); !ok {
			continue
		}
		if _, ok := hnm.ContainHash(uhHash, m.is); ok {
			continue
		}
		hnm.Append(uInfo, m.salt2, m.is)
	}

	ipfs.ToNameWithKeyFile(hnm.Marshal(), m.keyFile, m.is)
	return nil
}
