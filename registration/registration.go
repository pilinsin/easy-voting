package registration

import (
	"fmt"
	"time"
	"strings"

	"github.com/pilinsin/ipfs-util"
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
	rutil "github.com/pilinsin/easy-voting/registration/util"
)

type IRegistration interface {
	Close()
	VerifyHashBoxMap() bool
	VerifyUserIdentity(identity *rutil.UserIdentity) bool
	Registrate(userData ...string) (*rutil.UserIdentity, error)
}

type registration struct {
	is          *ipfs.IPFS
	psTopic     string
	rPubKey     crypto.IPubKey
	salt1       string
	salt2       string
	uhmCid      string
	hbmCid      string
	hbmIpnsName string
}

func NewRegistration(rCfgCid string, is *ipfs.IPFS) (*registration, error) {
	rCfg, err := rutil.ConfigFromCid(rCfgCid, is)
	if err != nil {
		return nil, err
	}
	hbmCid, err := ipfs.Name.GetCid(rCfg.HbmIpnsName(), is)
	if err != nil {
		return nil, err
	}

	r := &registration{
		is:          is,
		psTopic:     "registration_pubsub/" + rCfgCid,
		rPubKey:     rCfg.RPubKey(),
		salt1:       rCfg.Salt1(),
		salt2:       rCfg.Salt2(),
		uhmCid:      rCfg.UhmCid(),
		hbmCid:      hbmCid,
		hbmIpnsName: rCfg.HbmIpnsName(),
	}
	return r, nil
}
func (r *registration) Close() {
	r.is = nil
}
func (r *registration) VerifyHashBoxMap() bool {
	pth, _ := r.is.Name.Resolve(r.hbmIpnsName)
	mhbm, err := r.is.File.Get(pth)
	if err != nil {
		fmt.Println("hbmName error")
		return false
	}
	hbm, err = UnmarshalHashBoxMap(mhbm)
	if err != nil {
		fmt.Println("hbm unmarshal error")
		return false
	}

	if ok := hbm.VerifyHashes(r.uhmCid, r.is); !ok {
		fmt.Println("invalid uhHash is contained in hbm")
		return false
	}
	if hbm.VerifyCid(r.hbmCid, r.is) {
		r.hbmCid = pth.Cid().String()
		return true
	} else {
		fmt.Println("invalid hnm cid")
		return false
	}
}
func (r *registration) VerifyUserIdentity(identity *rutil.UserIdentity) bool {
	if hbm, err := HashBoxMapFromName(r.hnmIpnsName, r.is); err != nil {
		return false
	}else{
		return hbm.VerifyUserIdentity(identity, r.salt2, r.is)
	}
}
func (r *registration) Registrate(userData ...string) (*rutil.UserIdentity, error) {
	userHash := rutil.NewUserHash(r.salt1, userData...)
	uhHash := rutil.NewUhHash(r.salt2, userHash)

	uhm, err := UhHashMapFromCid(r.uhmCid, r.is)
	if err != nil {
		return nil, err
	}
	if ok := uhm.ContainHash(uhHash, r.is); !ok {
		return nil, util.NewError("uhHash is not contained")
	}
	hbm, err := HashBoxMapFromName(r.hbmIpnsName, r.is)
	if err != nil {
		return nil, err
	}
	if _, ok := hbm.ContainHash(uhHash, r.is); ok {
		return nil, util.NewError("uhHash is already registrated")
	}

	var userEncKeyPair IPubEncryptKeyPair
	for{
		userEncKeyPair = crypto.NewPubEncryptKeyPair()
		if ng := hbm.ContainPubKey(userEncKeyPair.Public()); !ng{break}
	}
	userSignKeyPair := crypto.NewSignKeyPair()
	rb := rutil.NewRegistrationBox(userEncKeyPair.Public())

	id := rutil.NewUserIdentity(
		userHash,
		userEncKeyPair.Public(),
		userEncKeyPair.Private(),
		userSignKeyPair.Sign(),
		userSignKeyPair.Verify(),
	)
	uInfo := rutil.NewUserInfo(userHash, rb)

	encInfo, err := r.rPubKey.Encrypt(uInfo.Marshal())
	if err != nil {
		return nil, util.AddError(err, "encUInfo err in r.Registrate")
	}
	r.is.PubSub().Publish(encInfo, r.psTopic)
	
	ticker := time.NewTicker(30*time.Second)
	defer ticker.Stop()
	for {
		hbm, err := HashBoxMapFromName(r.hbmIpnsName, r.is)
		if err != nil {
			return nil, err
		}
		if hbm.VerifyUserInfo(uInfo, r.salt2, r.is) {
			fmt.Println("uInfo verified")
			return id, nil
		}
		//fmt.Println("wait for registration")
		<-ticker.C
	}

}
