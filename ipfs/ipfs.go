package ipfs

import (
	"context"
	"time"
	//"fmt"
	"io/ioutil"
	"sync"

	dstore "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	//fstore "github.com/ipfs/go-filestore"
	config "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	kstore "github.com/ipfs/go-ipfs-keystore"
	core "github.com/ipfs/go-ipfs/core"
	coreapi "github.com/ipfs/go-ipfs/core/coreapi"
	//libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	repo "github.com/ipfs/go-ipfs/repo"
	//fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	iface "github.com/ipfs/interface-go-ipfs-core"
	options "github.com/ipfs/interface-go-ipfs-core/options"
	path "github.com/ipfs/interface-go-ipfs-core/path"
	peer "github.com/libp2p/go-libp2p-core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"

	"EasyVoting/ipfs/pubsub"
	"EasyVoting/util"
)

type IPFS struct {
	ipfsApi iface.CoreAPI
	ctx     context.Context
	kStore  kstore.Keystore
	ps      *pubsub.PubSub
}

func New(ctx context.Context, repoStr string, topic string) *IPFS {
	repoPath := repoStr
	//repoPath, err := ioutil.TempDir("", repoStr)
	//util.CheckError(err)
	//keyGenOpts := []options.KeyGenerateOption{options.Key.Type(options.Ed25519Key)}
	//id, err := config.CreateIdentity(ioutil.Discard, keyGenOpts)
	//util.CheckError(err)
	//cfg, err := config.InitWithIdentity(id)
	cfg, err := config.Init(ioutil.Discard, 2048)
	util.CheckError(err)
	//cfg.Datastore.Spec = map[string]interface{}{
	//	"type": "mem",
	//	"path": "mem",
	//}
	//cfg.Pubsub = config.PubsubConfig{"gossipsub", false}
	ds := dsync.MutexWrap(dstore.NewMapDatastore())
	//fm := fstore.NewFileManager(dstore.NewMapDatastore(), repoPath)
	//fm.AllowFiles = true
	ks, err := kstore.NewFSKeystore(repoPath)
	util.CheckError(err)
	r := &repo.Mock{D: ds, C: *cfg, K: ks} //, F: fm}
	//err = fsrepo.Init(repoPath, cfg)
	//util.CheckError(err)
	//r, err := fsrepo.Open(repoPath)
	//util.CheckError(err)
	exOpts := map[string]bool{
		//"discovery": false,
		//"dht":       true,
		//"pubsub":    true,
	}
	buildCfg := core.BuildCfg{
		Online: true,
		//Routing:   libp2p.DHTOption,
		//Host: mock.MockHostOption(mocknet.New(ctx)),
		Repo:      r,
		ExtraOpts: exOpts,
	}
	node, err := core.NewNode(ctx, &buildCfg)
	util.CheckError(err)
	coreApi, err := coreapi.NewCoreAPI(node)
	util.CheckError(err)

	iPFS := &IPFS{coreApi, ctx, r.Keystore(), pubsub.New2(ctx, topic)}
	//go iPFS.connect2Peers()
	return iPFS
}

func (ipfs *IPFS) connect2Peers() {
	var wg sync.WaitGroup
	pInfos := make(map[peer.ID]*peer.AddrInfo, len(bootstrapNodes()))
	for _, addrStr := range bootstrapNodes() {
		addr, err := multiaddr.NewMultiaddr(addrStr)
		util.CheckError(err)
		pInfo, err := peer.AddrInfoFromP2pAddr(addr)
		util.CheckError(err)
		pi, ok := pInfos[pInfo.ID]
		if !ok {
			pi = &peer.AddrInfo{ID: pInfo.ID}
			pInfos[pi.ID] = pi
		}
		pi.Addrs = append(pi.Addrs, pInfo.Addrs...)
	}
	wg.Add(len(pInfos))
	for _, pInfo := range pInfos {
		go func(peerInfo *peer.AddrInfo) {
			defer wg.Done()
			err := ipfs.ipfsApi.Swarm().Connect(ipfs.ctx, *peerInfo)
			util.CheckError(err)
		}(pInfo)
	}
	wg.Wait()
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

func (ipfs *IPFS) keyFileSet(kFile KeyFile, kw string) {
	if ipfs.hasKey(kw) {
		err := ipfs.kStore.Delete(kw)
		util.CheckError(err)
	}
	err := ipfs.kStore.Put(kw, kFile.keyFile)
	util.CheckError(err)

}

func (ipfs *IPFS) keySet(kw string) {
	if !ipfs.hasKey(kw) {
		_, err := ipfs.ipfsApi.Key().Generate(ipfs.ctx, kw)
		util.CheckError(err)
	}
}

func NameGet(kFile KeyFile) string {
	pid, err := peer.IDFromPrivateKey(kFile.keyFile)
	util.CheckError(err)
	name := iface.FormatKeyID(pid)

	return name
}

func (ipfs *IPFS) NamePublishWithKeyFile(pth path.Path, vt string, kFile KeyFile, kw string) iface.IpnsEntry {
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

func (ipfs *IPFS) PubsubLs() []string {
	return ipfs.ps.Ls()
}
func (ipfs *IPFS) PubsubPeers() []peer.ID {
	return ipfs.ps.Peers()
}

func (ipfs *IPFS) PubsubConnect() {
	//ipfs.ps.Connect(ipfs.ctx)
}

func (ipfs *IPFS) PubsubPublish(data []byte) {
	ipfs.ps.Publish(ipfs.ctx, data)
}

func (ipfs *IPFS) PubsubSubTest() {
	ipfs.ps.Subscribe(ipfs.ctx)
}
func (ipfs *IPFS) PubsubSubscribe(topic string) iface.PubSubSubscription {
	sub, err := ipfs.ipfsApi.PubSub().Subscribe(ipfs.ctx, topic)
	util.CheckError(err)
	defer sub.Close()

	return sub
}

func (ipfs *IPFS) SubNext(sub iface.PubSubSubscription) (iface.PubSubMessage, error) {
	return sub.Next(ipfs.ctx)
}
