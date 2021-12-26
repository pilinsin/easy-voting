package registration

import (
	"time"

	"EasyVoting/ipfs"
	"EasyVoting/ipfs/pubsub"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
	"EasyVoting/util/ecies"
	"EasyVoting/util/ed25519"
)

type IRegistration interface {
	Close()
	VerifyHashNameMap() bool
	Registrate(userData ...string) (*rutil.UserIdentity, error)
}

type registration struct {
	is          *ipfs.IPFS
	psTopic     string
	hnmCid      string
	rPubKey     *ecies.PubKey
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
	r.is.Close()
	r.is = nil
}
func (r *registration) VerifyHashNameMap() bool {
	chm := &rutil.ConstHashMap{}
	err := chm.FromCid(r.chmCid, r.is)
	if err != nil {
		return false
	}

	hnm := &rutil.HashNameMap{}
	pth, _ := r.is.NameResolve(r.hnmIpnsName)
	mhnm, err := r.is.FileGet(pth)
	if err != nil {
		return false
	}
	err = hnm.Unmarshal(mhnm)
	if err != nil {
		return false
	}

	if ok := hnm.VerifyHashes(chm, r.is); !ok {
		return false
	}
	if hnm.VerifyCid(r.hnmCid, r.is) {
		r.hnmCid = pth.String()
		return true
	} else {
		return false
	}
}
func (r *registration) Registrate(userData ...string) (*rutil.UserIdentity, error) {
	userHash := rutil.NewUserHash(r.is, r.salt1, userData...)
	uhHash := rutil.NewUhHash(r.is, r.salt2, userHash)

	chm := &rutil.ConstHashMap{}
	err := chm.FromCid(r.chmCid, r.is)
	if err != nil {
		return nil, err
	}
	if ok := chm.ContainHash(uhHash, r.is); !ok {
		return nil, util.NewError("uhHash is not contained")
	}
	hnm := &rutil.HashNameMap{}
	err = hnm.FromName(r.hnmIpnsName, r.is)
	if err != nil {
		return nil, err
	}
	if _, ok := hnm.ContainHash(uhHash, r.is); ok {
		return nil, util.NewError("uhHash is already registrated")
	}

	rKeyFile := ipfs.NewKeyFile()
	userEncKeyPair := ecies.NewKeyPair()
	userSignKeyPair := ed25519.NewKeyPair()

	rb := rutil.NewRegistrationBox(userEncKeyPair.Public(), userSignKeyPair.Verify())
	rIpnsName := ipfs.ToNameWithKeyFile(rb.Marshal(), rKeyFile, r.is)

	id := rutil.NewUserIdentity(userHash, rKeyFile, userEncKeyPair.Private(), userSignKeyPair.Sign())
	uInfo := rutil.NewUserInfo(userHash, rIpnsName)

	encInfo, err := r.rPubKey.Encrypt(uInfo.Marshal())
	if err != nil {
		return nil, util.AddError(err, "encUInfo err in r.Registrate")
	}
	ps, err := pubsub.New(r.psTopic)
	if err != nil {
		return nil, util.AddError(err, "pubsub.New error")
	}
	defer ps.Close()
	ps.Publish(encInfo)

	for {
		hnm := &rutil.HashNameMap{}
		err = hnm.FromName(r.hnmIpnsName, r.is)
		if err != nil {
			return nil, err
		}
		if hnm.VerifyUserInfo(uInfo, r.salt2, r.is) {
			return id, nil
		} else {
			<-time.After(30 * time.Second)
		}
	}
}
