package registration

import (
	"fmt"

	iface "github.com/ipfs/interface-go-ipfs-core"

	"EasyVoting/ipfs"
	rutil "EasyVoting/registration/util"
	//"EasyVoting/util"
	"EasyVoting/util/crypto/encrypt"
)

type IManager interface {
	Close()
	Registrate() error
}

type manager struct {
	is          *ipfs.IPFS
	sub         iface.PubSubSubscription
	priKey      *encrypt.PriKey
	keyFile     *ipfs.KeyFile
	salt2       string
	chmCid      string
	hnmIpnsName string
}

func NewManager(rCfgCid string, manIdentity *rutil.ManIdentity, is *ipfs.IPFS) (*manager, error) {
	rCfg, err := rutil.ConfigFromCid(rCfgCid, is)
	if err != nil {
		return nil, err
	}

	man := &manager{
		is:          is,
		sub:         is.PubSubSubscribe("registration_pubsub/" + rCfgCid),
		priKey:      manIdentity.Private(),
		keyFile:    	manIdentity.KeyFile(),
		salt2:       rCfg.Salt2(),
		chmCid:      rCfg.ChMapCid(),
		hnmIpnsName: rCfg.HnmIpnsName(),
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
	if err := chm.FromCid(m.chmCid, m.is); err != nil {
		fmt.Println("m.Registrate FromCid error", err)
		return err
	}
	hnm := &rutil.HashNameMap{}
	if err := hnm.FromName(m.hnmIpnsName, m.is); err != nil {
		fmt.Println("hnm.FromName error", err)
		return err
	}

	//it takes 5~6 mins
	subs := m.is.PubSubNextAll(m.sub)
	fmt.Println("data group: ", subs)
	if len(subs) <= 0 {
		return nil
	}
	isHnmUpdated := false
	for _, encUInfo := range subs {
		mUInfo, err := m.priKey.Decrypt(encUInfo)
		if err != nil {
			fmt.Println("decrypt uInfo error")
			continue
		}
		uInfo := &rutil.UserInfo{}
		if err := uInfo.Unmarshal(mUInfo); err != nil {
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

		if err := hnm.Append(uInfo, m.salt2, m.is); err == nil{
			isHnmUpdated = true
			fmt.Println("uInfo appended")
		}
	}
	if isHnmUpdated{
		name := ipfs.ToNameWithKeyFile(hnm.Marshal(), m.keyFile, m.is)
		fmt.Println("ipnsPublished to ", name)
	}
	return nil
}
