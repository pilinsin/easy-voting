package voting

import (
	"EasyVoting/ipfs"
	"EasyVoting/util/ecies"
)

type BaseVoting struct {
	iPFS      *ipfs.IPFS
	validTime string
	key       string
	nCands    int
	pubKey    ecies.PubKey
}

func (bv *BaseVoting) BaseInit(is *ipfs.IPFS, vt string, userID string, nCands int, pubKey ecies.PubKey) {
	bv.iPFS = is
	bv.validTime = vt
	bv.key = userID
	bv.nCands = nCands
	bv.pubKey = pubKey
}

func (bv *BaseVoting) NumCandsMatch(num int) bool {
	return num == bv.nCands
}

func (bv *BaseVoting) BaseVote(data []byte) string {
	encData := bv.pubKey.Encrypt(data)
	resolved := bv.iPFS.FileAdd(encData, true)
	ipnsEntry := bv.iPFS.NamePublish(resolved, bv.validTime, bv.key)
	return ipnsEntry.Name()
}

func (bv *BaseVoting) BaseGet(ipnsName string, priKey ecies.PriKey) []byte {
	pth := bv.iPFS.NameResolve(ipnsName)
	file := bv.iPFS.FileGet(pth)
	return priKey.Decrypt(file)
}
