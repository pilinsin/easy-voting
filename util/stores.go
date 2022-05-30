package util

import (
	"context"
	"errors"
	"strings"

	peer "github.com/libp2p/go-libp2p-core/peer"

	pv "github.com/pilinsin/p2p-verse"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
	crdt "github.com/pilinsin/p2p-verse/crdt"
)

func ParseConfigAddr(cfgAddr string) (string, string, error) {
	addrs := strings.Split(strings.TrimPrefix(cfgAddr, "/"), "/")
	pre := addrs[1]
	if len(addrs) < 3 || (pre != "r" && pre != "v") {
		return "", "", errors.New("invalid cfgAddr")
	}
	return addrs[0], addrs[2], nil
}

func NewIpfs(hGen pv.HostGenerator, dirName string, save bool, bootstraps []peer.AddrInfo) (ipfs.Ipfs, error) {
	return ipfs.NewIpfsStore(hGen, dirName, save, bootstraps...)
}
func NewStore(ctx context.Context, hGen pv.HostGenerator, stInfo [][2]string, dirName string, save bool, bootstraps []peer.AddrInfo, opts ...*crdt.StoreOpts) ([]crdt.IStore, error){
	v := crdt.NewVerse(hGen, dirName, save, bootstraps...)

	opt := &crdt.StoreOpts{}
	if len(opts) > 0{
		opt = opts[0]
	}
	stores := make([]crdt.IStore, 0)
	for idx := range stInfo{
		name := stInfo[idx][0]
		tp := stInfo[idx][1]
		st, err := v.LoadStore(ctx, name, tp, opt)
		if err != nil{return nil, err}
		stores = append(stores, st)
	}
	return stores, nil
}