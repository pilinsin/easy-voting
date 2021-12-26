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
	repo "github.com/ipfs/go-ipfs/repo"
	iface "github.com/ipfs/interface-go-ipfs-core"
	options "github.com/ipfs/interface-go-ipfs-core/options"
	nsopts "github.com/ipfs/interface-go-ipfs-core/options/namesys"
	path "github.com/ipfs/interface-go-ipfs-core/path"

	"EasyVoting/util"
)

type IPFS struct {
	ipfsApi iface.CoreAPI
	ctx     context.Context
	kStore  kstore.Keystore
}

func New(repoStr string) *IPFS {
	repoPath, _ := ioutil.TempDir("", repoStr)

	keyGenOpts := []options.KeyGenerateOption{options.Key.Type(options.Ed25519Key)}
	id, _ := config.CreateIdentity(ioutil.Discard, keyGenOpts)
	cfg, _ := config.InitWithIdentity(id)

	ds := dsync.MutexWrap(dstore.NewMapDatastore())
	ks, _ := kstore.NewFSKeystore(repoPath)
	r := &repo.Mock{D: ds, C: *cfg, K: ks}
	exOpts := map[string]bool{
		//"discovery": false,
		"dht": true,
		//"pubsub":    true,
	}
	buildCfg := core.BuildCfg{
		Online:    true,
		Repo:      r,
		ExtraOpts: exOpts,
	}

	ctx := context.Background()
	node, _ := core.NewNode(ctx, &buildCfg)
	coreApi, _ := coreapi.NewCoreAPI(node)

	return &IPFS{coreApi, ctx, r.Keystore()}
}

func (ipfs *IPFS) CoreApi() *iface.CoreAPI {
	return &ipfs.ipfsApi
}

func (ipfs *IPFS) FileAdd(data []byte, pn bool) path.Resolved {
	file := files.NewBytesFile(data)
	pth, _ := ipfs.ipfsApi.Unixfs().Add(ipfs.ctx, file, options.Unixfs.Pin(pn))
	return pth
}
func (ipfs *IPFS) FileHash(data []byte) path.Resolved {
	file := files.NewBytesFile(data)
	pth, _ := ipfs.ipfsApi.Unixfs().Add(ipfs.ctx, file, options.Unixfs.HashOnly(true))
	return pth
}

func (ipfs *IPFS) FileGet(pth path.Path) ([]byte, error) {
	f, err := ipfs.ipfsApi.Unixfs().Get(ipfs.ctx, pth)
	if err != nil {
		return nil, err
	} else {
		return ipfsFileNodeToBytes(f)
	}
}

func (ipfs *IPFS) hasKey(kw string) bool {
	keys, _ := ipfs.ipfsApi.Key().List(ipfs.ctx)
	for _, key := range keys {
		if key.Name() == kw {
			return true
		}
	}
	return false
}
func (ipfs *IPFS) keySet(kw string) {
	if kw == "self" {
		return
	}
	if !ipfs.hasKey(kw) {
		ipfs.ipfsApi.Key().Generate(ipfs.ctx, kw, options.Key.Type("ed25519"))
	}
}
func parseDuration(vt string) time.Duration {
	t, err := time.ParseDuration(vt)
	if err != nil {
		t, _ = time.ParseDuration("8760h")
	}
	return t
}

func (ipfs *IPFS) NamePublishWithKeyFile(pth path.Path, vt string, kFile *KeyFile) iface.IpnsEntry {
	t := parseDuration(vt)

	var kw string
	for {
		kw = util.GenUniqueID(50, 50)
		if ng := ipfs.hasKey(kw); !ng {
			break
		}
	}
	ipfs.kStore.Put(kw, kFile.keyFile)
	ipnsEntry, _ := ipfs.ipfsApi.Name().Publish(ipfs.ctx, pth, options.Name.ValidTime(t), options.Name.Key(kw))
	ipfs.kStore.Delete(kw)
	return ipnsEntry
}

func (ipfs *IPFS) NamePublish(pth path.Path, vt string, kw string) iface.IpnsEntry {
	t := parseDuration(vt)

	ipfs.keySet(kw)
	ipnsEntry, _ := ipfs.ipfsApi.Name().Publish(ipfs.ctx, pth, options.Name.ValidTime(t), options.Name.Key(kw))
	return ipnsEntry
}

func (ipfs *IPFS) NameResolve(name string) (path.Path, error) {
	return ipfs.ipfsApi.Name().Resolve(ipfs.ctx, name, options.Name.ResolveOption(nsopts.DhtRecordCount(1)))
}

func (ipfs *IPFS) Close() {
	ipfs.ipfsApi = nil
}
