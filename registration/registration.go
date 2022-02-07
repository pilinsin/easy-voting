package registration

import (
	"fmt"
	"time"
	"strings"

	"github.com/pilinsin/easy-voting/ipfs"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	"github.com/pilinsin/easy-voting/util"
	"github.com/pilinsin/easy-voting/util/crypto"
)

type IRegistration interface {
	Close()
	VerifyHashNameMap() bool
	VerifyUserIdentity(identity *rutil.UserIdentity) bool
	Registrate(userData ...string) (*rutil.UserIdentity, error)
}

type registration struct {
	is          *ipfs.IPFS
	psTopic     string
	hnmCid      string
	rPubKey     crypto.IPubKey
	salt1       string
	salt2       string
	chmCid      string
	hnmIpnsName string
}

func NewRegistration(rCfgCid string, is *ipfs.IPFS) (*registration, error) {
	rCfg, err := rutil.ConfigFromCid(rCfgCid, is)
	if err != nil {
		return nil, err
	}
	hnmCid, err := ipfs.CidFromName(rCfg.HnmIpnsName(), is)
	if err != nil {
		return nil, err
	}

	r := &registration{
		is:          is,
		psTopic:     "registration_pubsub/" + rCfgCid,
		hnmCid:      hnmCid,
		rPubKey:     rCfg.RPubKey(),
		salt1:       rCfg.Salt1(),
		salt2:       rCfg.Salt2(),
		chmCid:      rCfg.ChMapCid(),
		hnmIpnsName: rCfg.HnmIpnsName(),
	}
	return r, nil
}
func (r *registration) Close() {
	r.is = nil
}
func (r *registration) VerifyHashNameMap() bool {
	chm := &rutil.ConstHashMap{}
	err := chm.FromCid(r.chmCid, r.is)
	if err != nil {
		fmt.Println("chm unmarshal error")
		return false
	}

	hnm := &rutil.HashNameMap{}
	pth, _ := r.is.NameResolve(r.hnmIpnsName)
	mhnm, err := r.is.FileGet(pth)
	if err != nil {
		fmt.Println("hnmName error")
		return false
	}
	err = hnm.Unmarshal(mhnm)
	if err != nil {
		fmt.Println("hnm unmarshal error")
		return false
	}

	if ok := hnm.VerifyHashes(chm, r.is); !ok {
		fmt.Println("invalid uhHash is contained in hnm")
		return false
	}
	if hnm.VerifyCid(r.hnmCid, r.is) {
		r.hnmCid = strings.TrimPrefix(pth.String(), "/ipfs/")
		return true
	} else {
		fmt.Println("invalid hnm cid")
		return false
	}
}
func (r *registration) VerifyUserIdentity(identity *rutil.UserIdentity) bool {
	hnm := &rutil.HashNameMap{}
	if err := hnm.FromName(r.hnmIpnsName, r.is); err != nil {
		return false
	}
	return hnm.VerifyUserIdentity(identity, r.salt2, r.is)
}
func (r *registration) Registrate(userData ...string) (*rutil.UserIdentity, error) {
	userHash := rutil.NewUserHash(r.is, r.salt1, userData...)
	uhHash := rutil.NewUhHash(r.is, r.salt2, userHash)

	chm := &rutil.ConstHashMap{}
	if err := chm.FromCid(r.chmCid, r.is); err != nil {
		return nil, err
	}
	if ok := chm.ContainHash(uhHash, r.is); !ok {
		return nil, util.NewError("uhHash is not contained")
	}
	hnm := &rutil.HashNameMap{}
	if err := hnm.FromName(r.hnmIpnsName, r.is); err != nil {
		return nil, err
	}
	if _, ok := hnm.ContainHash(uhHash, r.is); ok {
		return nil, util.NewError("uhHash is already registrated")
	}

	rKeyFile := ipfs.NewKeyFile()
	userEncKeyPair := crypto.NewEncryptKeyPair()
	userSignKeyPair := crypto.NewSignKeyPair()

	rb := rutil.NewRegistrationBox(userEncKeyPair.Public())
	rIpnsName := ipfs.ToNameWithKeyFile(rb.Marshal(), rKeyFile, r.is)

	id := rutil.NewUserIdentity(userHash, rKeyFile, userEncKeyPair.Private(), userSignKeyPair.Sign())
	uInfo := rutil.NewUserInfo(userHash, rIpnsName)

	encInfo, err := r.rPubKey.Encrypt(uInfo.Marshal())
	if err != nil {
		return nil, util.AddError(err, "encUInfo err in r.Registrate")
	}
	r.is.PubSubPublish(encInfo, r.psTopic)
	//return id, nil
	
	ticker := time.NewTicker(30*time.Second)
	defer ticker.Stop()
	for {
		hnm := &rutil.HashNameMap{}
		if err := hnm.FromName(r.hnmIpnsName, r.is); err != nil {
			return nil, err
		}
		if hnm.VerifyUserInfo(uInfo, r.salt2, r.is) {
			fmt.Println("uInfo verified")
			return id, nil
		}
		//fmt.Println("wait for registration")
		<-ticker.C
	}

}
