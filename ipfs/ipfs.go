package ipfs

import (
	"context"
	"io/ioutil"
	"time"

	files "github.com/ipfs/go-ipfs-files"
	kstore "github.com/ipfs/go-ipfs-keystore"
	core "github.com/ipfs/go-ipfs/core"
	coreapi "github.com/ipfs/go-ipfs/core/coreapi"
	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	iface "github.com/ipfs/interface-go-ipfs-core"
	options "github.com/ipfs/interface-go-ipfs-core/options"
	nsopts "github.com/ipfs/interface-go-ipfs-core/options/namesys"
	ipath "github.com/ipfs/interface-go-ipfs-core/path"

	"EasyVoting/util"
)

type IPFS struct {
	ipfsApi iface.CoreAPI
	ctx     context.Context
	kStore  kstore.Keystore
}

func New(repoStr string) (*IPFS, error) {
	repoPath, _ := ioutil.TempDir("", repoStr)
	r, err := newRepo(repoPath)
	if err != nil {
		return nil, err
	}

	exOpts := map[string]bool{
		"discovery": true,
		"dht":       true,
		"pubsub":    true,
	}
	buildCfg := core.BuildCfg{
		Online:    true,
		Repo:      r,
		Routing: libp2p.DHTOption,
		ExtraOpts: exOpts,
	}

	ctx := context.Background()
	node, _ := core.NewNode(ctx, &buildCfg)
	api, _ := coreapi.NewCoreAPI(node)

	return &IPFS{api, ctx, r.Keystore()}, nil
}
func (ipfs *IPFS) Close() {
	ipfs.ipfsApi = nil
}

func (ipfs *IPFS) CoreApi() *iface.CoreAPI {
	return &ipfs.ipfsApi
}

func (ipfs *IPFS) FileAdd(data []byte, pn bool) ipath.Resolved {
	file := files.NewBytesFile(data)
	pth, _ := ipfs.ipfsApi.Unixfs().Add(ipfs.ctx, file, options.Unixfs.Pin(pn))
	return pth
}
func (ipfs *IPFS) FileHash(data []byte) ipath.Resolved {
	file := files.NewBytesFile(data)
	pth, _ := ipfs.ipfsApi.Unixfs().Add(ipfs.ctx, file, options.Unixfs.HashOnly(true))
	return pth
}

func (ipfs *IPFS) FileGet(pth ipath.Path) ([]byte, error) {
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

func (ipfs *IPFS) NamePublishWithKeyFile(pth ipath.Path, vt string, kFile *KeyFile) iface.IpnsEntry {
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

func (ipfs *IPFS) NamePublish(pth ipath.Path, vt string, kw string) iface.IpnsEntry {
	t := parseDuration(vt)

	ipfs.keySet(kw)
	ipnsEntry, _ := ipfs.ipfsApi.Name().Publish(ipfs.ctx, pth, options.Name.ValidTime(t), options.Name.Key(kw))
	return ipnsEntry
}

func (ipfs *IPFS) NameResolve(name string) (ipath.Path, error) {
	return ipfs.ipfsApi.Name().Resolve(ipfs.ctx, name, options.Name.ResolveOption(nsopts.DhtRecordCount(1)))
}

func (ipfs *IPFS) PubSubPublish(data []byte, topic string) {
	ipfs.ipfsApi.PubSub().Publish(ipfs.ctx, topic, data)
}
func (ipfs *IPFS) PubSubSubscribe(topic string) iface.PubSubSubscription {
	sub, _ := ipfs.ipfsApi.PubSub().Subscribe(ipfs.ctx, topic, options.PubSub.Discover(true))
	return sub
}
func (ipfs *IPFS) PubSubNext(sub iface.PubSubSubscription) []byte {
	ctx, cancel := util.CancelTimerContext(5*time.Second)
	defer cancel()

	msg, err := sub.Next(ctx)
	if err != nil {
		return nil
	}
	return msg.Data()
}
func (ipfs *IPFS) PubSubNextAll(sub iface.PubSubSubscription) [][]byte {
	var dataset [][]byte
	for {
		data := ipfs.PubSubNext(sub)
		if data == nil {
			if len(dataset) > 0 {
				return dataset
			} else {
				return nil
			}
		}
		dataset = append(dataset, data)
	}
}
