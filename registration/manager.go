package registration

import (
	"fmt"

	iface "github.com/ipfs/interface-go-ipfs-core"

	"github.com/pilinsin/ipfs-util"
	"github.com/pilinsin/util/crypto"
	rutil "github.com/pilinsin/easy-voting/registration/util"
)

type IManager interface {
	Close()
	Registrate() error
}

type manager struct {
	is          *ipfs.IPFS
	sub         iface.PubSubSubscription
	priKey      crypto.IPriKey
	keyFile     *ipfs.KeyFile
	salt2       string
	uhmCid      string
	hbmIpnsName string
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
		chmCid:      rCfg.UhmCid(),
		hbmIpnsName: rCfg.HbmIpnsName(),
	}
	return man, nil
}
func (m *manager) Close() {
	m.sub.Close()
	m.sub = nil
	m.is = nil
}
func (m *manager) Registrate() error {
	uhm, err := UhHashMapFromCid(r.uhmCid, r.is)
	if err != nil {
		fmt.Println("m.Registrate FromCid error", err)
		return err
	}
	hbm, err := HashBoxMapFromName(r.hbmIpnsName, r.is)
	if err != nil {
		fmt.Println("hbm.FromName error", err)
		return err
	}

	subs := m.is.PubSub().NextAll(m.sub)
	fmt.Println("data list: ", subs)
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
		uhHash := rutil.NewUhHash(m.salt2, uInfo.UserHash())
		if ok := uhm.ContainHash(uhHash, m.is); !ok {
			fmt.Println("the uhHash is not contained")
			continue
		}
		if _, ok := hbm.ContainHash(uhHash, m.is); ok {
			fmt.Println("the uhHash is already registrated")
			continue
		}

		if err := hbm.Append(uInfo, m.salt2, m.is); err == nil{
			isHbmUpdated = true
			fmt.Println("uInfo appended")
		}
	}
	if isHbmUpdated{
		name := ipfs.Name.PublishWithKeyFile(hbm.Marshal(), m.keyFile, m.is)
		fmt.Println("ipnsPublished to ", name)
	}
	return nil
}
