package registration

import (
	"context"
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
	"github.com/pilinsin/util/crypto"
)

type registration struct {
	ctx    context.Context
	cancel func()
	salt1  string
	addr   string
	is     ipfs.Ipfs
	uhm    crdt.IStore
	cfg    *rutil.Config
}

func NewRegistration(ctx context.Context, rCfgAddr, baseDir string) (riface.IRegistration, error) {
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
	stores, err := evutil.NewStore(ctx, i2p.NewI2pHost, stInfo, storeDir, save, bootstraps, opt)
	if err != nil {
		return nil, err
	}
	uhm := stores[0]

	ctx, cancel := context.WithCancel(context.Background())

	return &registration{
		ctx:    ctx,
		cancel: cancel,
		salt1:  rCfg.Salt1,
		addr:   rCfgAddr,
		is:     is,
		uhm:    uhm,
		cfg:    rCfg,
	}, nil
}

func (r *registration) Close() {
	r.cancel()
	r.is.Close()
	r.uhm.Close()
}

func (r *registration) Config() *rutil.Config {
	return r.cfg
}
func (r *registration) Address() string {
	return r.addr
}

func (r *registration) hasPubKey(pubKey crypto.IPubKey) bool {
	rs, err := r.uhm.Query()
	if err != nil {
		return false
	}
	mpub, err := crypto.MarshalPubKey(pubKey)
	if err != nil {
		return false
	}
	rs = query.NaiveFilter(rs, crdt.ValueMatchFilter{mpub})
	resList, err := rs.Rest()
	rs.Close()

	return len(resList) > 0 && err == nil
}

func (r *registration) Registrate(userData ...string) (string, error) {
	userHash := rutil.NewUserHash(r.salt1, userData...)
	if ok, err := r.uhm.Has(userHash); ok && err == nil {
		return "", errors.New("already registrated")
	}

	var userEncKeyPair crypto.IPubEncryptKeyPair
	for {
		userEncKeyPair = crypto.NewPubEncryptKeyPair()
		if has := r.hasPubKey(userEncKeyPair.Public()); !has {
			break
		}
	}

	id := rutil.NewUserIdentity(
		userHash,
		userEncKeyPair.Public(),
		userEncKeyPair.Private(),
	)

	mpub, err := crypto.MarshalPubKey(userEncKeyPair.Public())
	if err != nil {
		return "", err
	}
	if err := r.uhm.Put(userHash, mpub); err != nil {
		return "", err
	}

	time.Sleep(15 * time.Second)
	return id.ToString(), nil
}
