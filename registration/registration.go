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
	VerifyHashNameMap() bool
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
	hnmCid      string
	hnmIpnsName string
}

func NewRegistration(rCfgCid string, is *ipfs.IPFS) (*registration, error) {
	rCfg, err := rutil.ConfigFromCid(rCfgCid, is)
	if err != nil {
		return nil, err
	}
	hnmCid, err := ipfs.Name.GetCid(rCfg.HnmIpnsName(), is)
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
		hnmCid:      hnmCid,
		hnmIpnsName: rCfg.HnmIpnsName(),
	}
	return r, nil
}
func (r *registration) Close() {
	r.is = nil
}
func (r *registration) VerifyHashNameMap() bool {
	pth, _ := r.is.Name.Resolve(r.hnmIpnsName)
	mhnm, err := r.is.File.Get(pth)
	if err != nil {
		fmt.Println("hnmName error")
		return false
	}
	hnm, err = UnmarshalHashNameMap(mhnm)
	if err != nil {
		fmt.Println("hnm unmarshal error")
		return false
	}

	if ok := hnm.VerifyHashes(r.uhmCid, r.is); !ok {
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
	if hnm, err := HashNameMapFromName(r.hnmIpnsName, r.is); err != nil {
		return false
	}else{
		return hnm.VerifyUserIdentity(identity, r.salt2, r.is)
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
	hnm, err := HashNameMapFromName(r.hnmIpnsName, r.is)
	if err != nil {
		return nil, err
	}
	if _, ok := hnm.ContainHash(uhHash, r.is); ok {
		return nil, util.NewError("uhHash is already registrated")
	}

	rKeyFile := ipfs.Name.NewKeyFile()
	userEncKeyPair := crypto.NewPubEncryptKeyPair()
	userSignKeyPair := crypto.NewSignKeyPair()

	rb := rutil.NewRegistrationBox(userEncKeyPair.Public())
	rIpnsName := ipfs.Name.PublishWithKeyfile(rb.Marshal(), rKeyFile, r.is)

	id := rutil.NewUserIdentity(userHash, rKeyFile, userEncKeyPair.Private(), userSignKeyPair.Sign())
	uInfo := rutil.NewUserInfo(userHash, rIpnsName)

	encInfo, err := r.rPubKey.Encrypt(uInfo.Marshal())
	if err != nil {
		return nil, util.AddError(err, "encUInfo err in r.Registrate")
	}
	r.is.PubSub().Publish(encInfo, r.psTopic)
	
	ticker := time.NewTicker(30*time.Second)
	defer ticker.Stop()
	for {
		hnm, err := HashNameMapFromName(r.hnmIpnsName, r.is)
		if err != nil {
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
