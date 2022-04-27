package registration

import (
	"errors"
	"context"
	"path/filepath"
	"time"

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
	salt1   string
	idStr	string
	is ipfs.Ipfs
	uhm    	crdt.IStore
	cfg    *rutil.Config
}

func NewRegistration(ctx context.Context, rCfgAddr, idStr string) (riface.IRegistration, error) {
	bAddr, rCfgCid, err := evutil.ParseConfigAddr(rCfgAddr)
	if err != nil{return nil, err}

	ipfsDir, storeDir, save := parseIdStr(idStr)
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

func parseIdStr(idStr string) (string, string, bool){
	mi := &rutil.ManIdentity{}
	if err := mi.FromString(idStr); err == nil{
		return mi.IpfsDir, mi.StoreDir, true
	}
	title := pv.RandString(8)
	ipfsDir := filepath.Join(title, pv.RandString(8))
	storeDir := filepath.Join(title, pv.RandString(8))
	return ipfsDir, storeDir, false
}

func (r *registration) Close() {
	r.cancel()
	r.is.Close()
	r.uhm.Close()
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
	if _, _,  man := parseIdStr(r.idStr); man{
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
//ClhkZ2thdGxXLU5tQmhMY0hiRWZrNzlpYm9UVHFGWExXQXB5ZFZScGNOLVhfSHhqRF9INmt5YXdOS2FkMDJjSnhKeHo3VGh6RWZoUGctNkhRVlg4aFJRQT09EusCCugCCs8C3cZ0ZGf5e6Vr-puVTHmkTBhMhmi4RQ6zV8SRhh71OkVppPJCDw24-ldbP2yLWk8BSPk3-DbrCxoh0804yNaVHMQtTEqd23rawLYG7bnqbu7xom6bcUbpYQbuoUN1F-ICpyBQokUtQ4ShiyaN4YHcNSAt83XClt4Fcauu05tGwCls8J-hlOIA3SZWNFGsZjoFRp4GISmq5MnE9-UQpBWbGBeKHSddvgrhpTY0fZtJlCgcoyECGGx7b0wL6fxdS8lpMPQsJWHAkwGR1n5ngLO-NSkeAVD-d93oeItUeQ40Wj3Jt_N1DjkflSrf80g8UY4igwoY25Ev7lLmMkEjAUhhid8v6zqLAJaR89nmoT8cg0DUkUegC02oJUgZFjr9Mw4BP9zexOqowOWb9BQGLkRDv9bjRrCEHpck0KNo2loVlpfWbFjSW0rWhJltRysBAgISFFNJREgtcDc1MS1jb21wcmVzc2VkGk0KSwowOA0CduNuWgKyzsEOHoo_uvQer83tgQBYAKagClygRswZY5WZAVOsIyTtt5jP8gQAEhRTSURILXA3NTEtY29tcHJlc3NlZBjPAg==
//ClhPU1hSbzNnZ1RCSG92U0w5TW5vWGkyVDU5amRrLTlNTkV4MjE2a0JXWDlCMEhwaXQtd2NNQkpRT2h1ZUcyNnMxV3dJalliNDFlLWZZdXlmZTBHTmdQdz09EusCCugCCs8CrUq9T6D9SGAdI1OqBTp_utWrzJUQu1X5zsuYdPrDLWDWsKII7iDaB7Aq8paM__EDNoXkk_uyeuOpah-OUpXZKd8ikdJlpTooR0zrDcq3AxOfFC2upsVx7Oe_ShsJcqIElfvZV9EnAe39YzQhrqemwZ0T2zGB2R_KdBtWXkIPOb1PbByJeTO5duJBzKaFCfADHk0yAppSRL_J0l35zPpdy5vdn70FR0IptSxS2kqfvc7uHJRkKERPeF-cCMH7YxjB1csnnGng139FOX2l3JAtVK9s5TdB1tJq5P9z674wROFSFaXlUV1Pq81SsuMBTYtys_mzlDmVdt2GucDZnVFoPZ5h62NEnmThjSfCmchx4ChV5IhfrA8dIFd04FulfIpkK8m-AQemRaUkZ1mjoXeYkYOXjWj4p3JOBvnNPrEAzPtstFX-Qopw5L9drwiBAgESFFNJREgtcDc1MS1jb21wcmVzc2VkGk0KSwowxiC6pjEOAlnYcfWFSOCj1IIs_HmWTHpsk1qX4nyZKxXRcFFaVL7TYt7cxjZFQwgAEhRTSURILXA3NTEtY29tcHJlc3NlZBjPAg==
