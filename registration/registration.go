package registration

import (
	"errors"
	"context"
	"path/filepath"
	"time"
	"os"

	query "github.com/ipfs/go-datastore/query"

	"github.com/pilinsin/util/crypto"
	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
	crdt "github.com/pilinsin/p2p-verse/crdt"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
	evutil "github.com/pilinsin/easy-voting/util"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	riface "github.com/pilinsin/easy-voting/registration/interface"
)

type registration struct {
	ctx context.Context
	cancel func()
	dirCancel func()
	salt1   string
	idStr	string
	is ipfs.Ipfs
	uhm    	crdt.IStore
	cfg    *rutil.Config
}

func NewRegistration(ctx context.Context, rCfgAddr, idStr string) (riface.IRegistration, error) {
	bAddr, rCfgCid, err := evutil.ParseConfigAddr(rCfgAddr)
	if err != nil{return nil, err}

	dirCancel := func(){}
	baseDir, ipfsDir, storeDir, save := parseIdStr(idStr)
	if !save{
		dirCancel = func(){os.RemoveAll(baseDir)}
	}
	is, err := evutil.NewIpfs(i2p.NewI2pHost, bAddr, ipfsDir, save)
	if err != nil{return nil, err}
	rCfg := &rutil.Config{}
	if err := rCfg.FromCid(rCfgCid, is); err != nil{return nil, err}

	bootstraps := pv.AddrInfosFromString(bAddr)
	v := crdt.NewVerse(i2p.NewI2pHost, storeDir, save, false, bootstraps...)
	opt := &crdt.StoreOpts{Salt:rCfg.Salt2}
	uhm, err := v.LoadStore(ctx, rCfg.UhmAddr, "hash", opt)
	if err != nil{return nil, err}

	ctx, cancel := context.WithCancel(context.Background())
	autoSync(ctx, uhm)

	return &registration{
		ctx:	ctx,
		cancel:	cancel,
		dirCancel: dirCancel,
		salt1:  rCfg.Salt1,
		idStr:	idStr,
		is:		is,
		uhm: 	uhm,
		cfg:	rCfg,
	}, nil
}
func autoSync(ctx context.Context, uhm crdt.IStore){
	ticker := time.NewTicker(time.Second*10)
	go func(){
		for{
			select{
			case <-ctx.Done():
				return
			case <-ticker.C:
				uhm.Sync()
			}
		}
	}()
}

func parseIdStr(idStr string) (string, string, string, bool){
	mi := &rutil.ManIdentity{}
	if err := mi.FromString(idStr); err == nil{
		return "", mi.IpfsDir, mi.StoreDir, true
	}
	title := pv.RandString(8)
	ipfsDir := filepath.Join(title, pv.RandString(8))
	storeDir := filepath.Join(title, pv.RandString(8))
	return title, ipfsDir, storeDir, false
}

func (r *registration) Close() {
	r.cancel()
	r.is.Close()
	r.uhm.Close()

	time.Sleep(time.Second*10)
	r.dirCancel()
}

func (r *registration) Config() *rutil.Config{
	return r.cfg
}

func (r *registration) hasPubKey(pubKey crypto.IPubKey) bool{
	rs, err := r.uhm.Query()
	if err != nil{return false}
	mpub, err := crypto.MarshalPubKey(pubKey)
	if err != nil{return false}
	rs = query.NaiveFilter(rs, crdt.ValueMatchFilter{mpub})
	resList, err := rs.Rest()
	rs.Close()

	return len(resList) > 0 && err == nil
}
func (r *registration) Registrate(userData ...string) (string, error) {
	if _, _, _,  man := parseIdStr(r.idStr); man{
		return "", errors.New("manager can not registrate")
	}

	userHash := rutil.NewUserHash(r.salt1, userData...)
	if ok, err := r.uhm.Has(userHash); ok && err == nil{
		return "", errors.New("already registrated")
	}

	var userEncKeyPair crypto.IPubEncryptKeyPair
	for{
		userEncKeyPair = crypto.NewPubEncryptKeyPair()
		if has := r.hasPubKey(userEncKeyPair.Public()); !has{break}
	}

	id := rutil.NewUserIdentity(
		userHash,
		userEncKeyPair.Public(),
		userEncKeyPair.Private(),
	)

	mpub, err := crypto.MarshalPubKey(userEncKeyPair.Public())
	if err != nil{return "", err}
	if err := r.uhm.Put(userHash, mpub); err != nil{
		return "", err
	}

	time.Sleep(15*time.Second)
	return id.ToString(), nil
}
