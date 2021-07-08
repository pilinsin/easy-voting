package voting

import (
	crsa "crypto/rsa"

	"EasyVoting/ipfs"
	"EasyVoting/util/rsa"
)

type BaseVoting struct {
	iPFS      *ipfs.IPFS
	validTime string
	key       string
	nCands    int
}

func (bv *BaseVoting) BaseInit(is *ipfs.IPFS, vt string, userID string, nCands int) {
	bv.iPFS = is
	bv.validTime = vt
	bv.key = userID
	bv.nCands = nCands
}

func (bv *BaseVoting) NumCandsMatch(num int) bool {
	return num == bv.nCands
}

func (bv *BaseVoting) BaseVote(data []byte, pubKey crsa.PublicKey) string {
	encData := rsa.Encrypt(data, pubKey)
	resolved := bv.iPFS.FileAdd(encData, true)
	ipnsEntry := bv.iPFS.NamePublish(resolved, bv.validTime, bv.key)
	return ipnsEntry.Name()
}

func (bv *BaseVoting) BaseGet(ipnsName string, priKey crsa.PrivateKey) []byte {
	pth := bv.iPFS.NameResolve(ipnsName)
	file := bv.iPFS.FileGet(pth)
	return rsa.Decrypt(file, priKey)
}
