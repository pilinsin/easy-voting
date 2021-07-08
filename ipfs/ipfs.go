package ipfs

import (
	"context"
	"io/ioutil"
	"time"

	dstore "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	config "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	kstore "github.com/ipfs/go-ipfs-keystore"
	core "github.com/ipfs/go-ipfs/core"
	coreapi "github.com/ipfs/go-ipfs/core/coreapi"
	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	repo "github.com/ipfs/go-ipfs/repo"
	iface "github.com/ipfs/interface-go-ipfs-core"
	options "github.com/ipfs/interface-go-ipfs-core/options"
	path "github.com/ipfs/interface-go-ipfs-core/path"
	p2pcrypt "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"

	"EasyVoting/util"
)

type IPFS struct {
	ipfsApi iface.CoreAPI
	ctx     context.Context
	kStore  kstore.Keystore
}

func New(ctx context.Context, repoStr string) *IPFS {
	ds := dsync.MutexWrap(dstore.NewMapDatastore())
	cfg, err := config.Init(ioutil.Discard, 2048)
	util.CheckError(err)
	ks, err := kstore.NewFSKeystore(repoStr)
	util.CheckError(err)
	r := repo.Mock{D: ds, C: *cfg, K: ks}
	exOpts := map[string]bool{
		"discovery": true,
		"dht":       true,
	}
	buildCfg := core.BuildCfg{
		Online:    true,
		Routing:   libp2p.DHTOption,
		Repo:      &r,
		ExtraOpts: exOpts,
	}
	node, err := core.NewNode(ctx, &buildCfg)
	util.CheckError(err)
	coreApi, err := coreapi.NewCoreAPI(node)
	util.CheckError(err)
	return &IPFS{ipfsApi: coreApi, ctx: ctx, kStore: r.Keystore()}
}

func (ipfs *IPFS) CoreApi() *iface.CoreAPI {
	return &ipfs.ipfsApi
}

func (ipfs *IPFS) FileAdd(data []byte, pn bool) path.Resolved {
	file := files.NewBytesFile(data)
	resolved, err := ipfs.ipfsApi.Unixfs().Add(ipfs.ctx, file, options.Unixfs.Pin(pn))
	util.CheckError(err)
	return resolved
}

func (ipfs *IPFS) FileGet(pth path.Path) []byte {
	f, err := ipfs.ipfsApi.Unixfs().Get(ipfs.ctx, pth)
	util.CheckError(err)
	return IpfsFileNode2Bytes(f)
}

func (ipfs *IPFS) hasKey(kw string) bool {
	keys, err := ipfs.ipfsApi.Key().List(ipfs.ctx)
	util.CheckError(err)
	isExsitKey := false
	for _, key := range keys {
		if key.Name() == kw {
			isExsitKey = true
		}
	}
	return isExsitKey
}

func (ipfs *IPFS) keyFileSet(kFile p2pcrypt.PrivKey, kw string) {
	if ipfs.hasKey(kw) {
		err := ipfs.kStore.Delete(kw)
		util.CheckError(err)
	}
	err := ipfs.kStore.Put(kw, kFile)
	util.CheckError(err)

}

func (ipfs *IPFS) keySet(kw string) {
	if !ipfs.hasKey(kw) {
		_, err := ipfs.ipfsApi.Key().Generate(ipfs.ctx, kw)
		util.CheckError(err)
	}
}

func NameGet(kFile p2pcrypt.PrivKey) string {
	pid, err := peer.IDFromPrivateKey(kFile)
	util.CheckError(err)
	name := iface.FormatKeyID(pid)

	return name
}

func (ipfs *IPFS) NamePublishWithKeyFile(pth path.Path, vt string, kFile p2pcrypt.PrivKey, kw string) iface.IpnsEntry {
	t, err := time.ParseDuration(vt)
	util.CheckError(err)

	ipfs.keyFileSet(kFile, kw)

	ipnsEntry, err := ipfs.ipfsApi.Name().Publish(ipfs.ctx, pth, options.Name.ValidTime(t), options.Name.Key(kw))
	util.CheckError(err)
	return ipnsEntry
}

func (ipfs *IPFS) NamePublish(pth path.Path, vt string, kw string) iface.IpnsEntry {
	t, err := time.ParseDuration(vt)
	util.CheckError(err)

	ipfs.keySet(kw)

	ipnsEntry, err := ipfs.ipfsApi.Name().Publish(ipfs.ctx, pth, options.Name.ValidTime(t), options.Name.Key(kw))
	util.CheckError(err)
	return ipnsEntry
}

func (ipfs *IPFS) NameResolve(name string) path.Path {
	pth, err := ipfs.ipfsApi.Name().Resolve(ipfs.ctx, name)
	util.CheckError(err)
	return pth
}
