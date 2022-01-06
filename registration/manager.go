package registration

import (
	"fmt"

	iface "github.com/ipfs/interface-go-ipfs-core"

	"EasyVoting/ipfs"
	rutil "EasyVoting/registration/util"
	//"EasyVoting/util"
	"EasyVoting/util/ecies"
)

type IManager interface {
	Close()
	Registrate() error
}

type manager struct {
	is          *ipfs.IPFS
	sub         iface.PubSubSubscription
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
	rCfgCid := ipfs.ToCid(mCfg.Config().Marshal(), is)

	man := &manager{
		is:          is,
		sub:         is.PubSubSubscribe("registration_pubsub/" + rCfgCid),
		priKey:      mCfg.Private(),
		keyFile:     mCfg.KeyFile(),
		salt2:       mCfg.Salt2(),
		chmCid:      mCfg.ChMapCid(),
		hnmIpnsName: mCfg.HnmIpnsName(),
	}
	return man, nil
}
func (m *manager) Close() {
	m.sub.Close()
	m.sub = nil
	m.is = nil
}
func (m *manager) Registrate() error {
	chm := &rutil.ConstHashMap{}
	err := chm.FromCid(m.chmCid, m.is)
	if err != nil {
		fmt.Println("m.Registrate FromCid error", err)
		return err
	}
	hnm := &rutil.HashNameMap{}
	hnm.FromName(m.hnmIpnsName, m.is)
	if err != nil {
		fmt.Println("hnm.FromName error", err)
		return err
	}

	subs := m.is.PubSubNextAll(m.sub)
	fmt.Println("data group: ", subs)
	if len(subs) <= 0 {
		return nil
	}

	for _, encUInfo := range subs {
		mUInfo, err := m.priKey.Decrypt(encUInfo)
		if err != nil {
			fmt.Println("decrypt uInfo error")
			continue
		}
		uInfo := &rutil.UserInfo{}
		err = uInfo.Unmarshal(mUInfo)
		if err != nil {
			fmt.Println("uInfo unmarshal error")
			continue
		}

		uhHash := rutil.NewUhHash(m.is, m.salt2, uInfo.UserHash())
		if ok := chm.ContainHash(uhHash, m.is); !ok {
			fmt.Println("the uhHash is not contained")
			continue
		}
		if _, ok := hnm.ContainHash(uhHash, m.is); ok {
			fmt.Println("the uhHash is already registrated")
			continue
		}
		hnm.Append(uInfo, m.salt2, m.is)
		fmt.Println("uInfo appended")
	}

	name := ipfs.ToNameWithKeyFile(hnm.Marshal(), m.keyFile, m.is)
	fmt.Println("ipnsPublished to ", name)
	return nil
}
