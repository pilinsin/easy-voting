package util

import(
	"errors"
	"strings"
	pv "github.com/pilinsin/p2p-verse"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
)

func ParseConfigAddr(cfgAddr string) (string, string, error){
	addrs := strings.Split(strings.TrimPrefix(cfgAddr, "/"), "/")
	pre := addrs[0]
	if len(addrs)<3 || (pre != "r" && pre != "v"){
		return "", "", errors.New("invalid cfgAddr")
	}
	return addrs[1], addrs[2], nil
}


func NewIpfs(hGen pv.HostGenerator, bAddr, dirName string, save bool) (ipfs.Ipfs, error){
	bootstraps := pv.AddrInfosFromString(bAddr)
	ipfsDir := dirName+"_ipfs"
	return ipfs.NewIpfsStore(hGen, ipfsDir, "ipfs_kw", save, false, bootstraps...)
}

