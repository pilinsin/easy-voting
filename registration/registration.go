package registration

import (
	"errors"
	"path/filepath"
	"time"

	query "github.com/ipfs/go-datastore/query"

	riface "github.com/pilinsin/easy-voting/registration/interface"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	evutil "github.com/pilinsin/easy-voting/util"
	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
	crdt "github.com/pilinsin/p2p-verse/crdt"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
)

type registration struct {
	salt1 string
	addr  string
	is    ipfs.Ipfs
	uhm   crdt.IStore
	cfg   *rutil.Config
}

func NewRegistration(rCfgAddr, baseDir string) (riface.IRegistration, error) {
	bAddr, rCfgCid, err := evutil.ParseConfigAddr(rCfgAddr)
	if err != nil {
		return nil, err
	}
	save := true
	bootstraps := pv.AddrInfosFromString(bAddr)

	ipfsDir := filepath.Join(baseDir, "ipfs")
	is, err := evutil.NewIpfs(i2p.NewI2pHost, ipfsDir, save, bootstraps)
	if err != nil {
		return nil, err
	}
	rCfg := &rutil.Config{}
	if err := rCfg.FromCid(rCfgCid, is); err != nil {
		return nil, err
	}

	storeDir := filepath.Join(baseDir, "store")
	stInfo := [][2]string{{rCfg.UhmAddr, "hash"}}
	opt := &crdt.StoreOpts{Salt: rCfg.Salt2}
	stores, err := evutil.NewStore(i2p.NewI2pHost, stInfo, storeDir, save, bootstraps, opt)
	if err != nil {
		return nil, err
	}
	uhm := stores[0]

	return &registration{
		salt1: rCfg.Salt1,
		addr:  rCfgAddr,
		is:    is,
		uhm:   uhm,
		cfg:   rCfg,
	}, nil
}
func NewRegistrationWithStores(rCfgAddr string, is ipfs.Ipfs, uhm crdt.IStore) (riface.IRegistration, error) {
	_, rCfgCid, err := evutil.ParseConfigAddr(rCfgAddr)
	if err != nil {
		return nil, err
	}
	rCfg := &rutil.Config{}
	if err := rCfg.FromCid(rCfgCid, is); err != nil {
		return nil, err
	}
	return &registration{
		salt1: rCfg.Salt1,
		addr:  rCfgAddr,
		is:    is,
		uhm:   uhm,
		cfg:   rCfg,
	}, nil
}

func (r *registration) Close() {
	r.is.Close()
	r.uhm.Close()
}

func (r *registration) Config() *rutil.Config {
	return r.cfg
}
func (r *registration) Address() string {
	return r.addr
}

func (r *registration) hasPubKey(pubKey evutil.IPubKey) bool {
	rs, err := r.uhm.Query()
	if err != nil {
		return false
	}
	mpub, err := pubKey.Raw()
	if err != nil {
		return false
	}
	rs = query.NaiveFilter(rs, crdt.ValueMatchFilter{Val: mpub})
	resList, err := rs.Rest()
	rs.Close()

	return len(resList) > 0 && err == nil
}

func (r *registration) Registrate(userData ...string) (string, error) {
	userHash := rutil.NewUserHash(r.salt1, userData...)
	if ok, err := r.uhm.Has(userHash); ok && err == nil {
		return "", errors.New("already registrated")
	}

	userEncKeyPair := evutil.NewPubKeyPair()
	for {
		if has := r.hasPubKey(userEncKeyPair.Public()); !has {
			break
		}
		userEncKeyPair = evutil.NewPubKeyPair()
	}

	id := rutil.NewUserIdentity(
		userHash,
		userEncKeyPair.Public(),
		userEncKeyPair.Private(),
	)

	mpub, err := userEncKeyPair.Public().Raw()
	if err != nil {
		return "", err
	}
	if err := r.uhm.Put(userHash, mpub); err != nil {
		return "", err
	}

	time.Sleep(15 * time.Second)
	return id.ToString(), nil
}
