package voting

import (
	crsa "crypto/rsa"
	"encoding/json"

	"github.com/ipfs/interface-go-ipfs-core/path"

	"EasyVoting/ipfs"
	"EasyVoting/util"
	"EasyVoting/util/rsa"
)

func VerifyUserID(is *ipfs.IPFS, userID string, huidListAddrs []string) bool {
	for _, huidListAddr := range huidListAddrs {
		huidListPath := path.New(huidListAddr)
		huidListBytes := is.FileGet(huidListPath)

		var huidList map[string]struct{}
		err := json.Unmarshal(huidListBytes, &huidList)
		util.CheckError(err)

		huid := util.Hash(userID)
		_, ok := huidList[huid]
		if ok {
			return true
		}

	}
	return false
}

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
	encData := rsa.Encrypt(string(data), pubKey)
	resolved := bv.iPFS.FileAdd([]byte(encData), true)
	ipnsEntry := bv.iPFS.NamePublish(resolved, bv.validTime, bv.key)
	return ipnsEntry.Name()
}

func (bv *BaseVoting) BaseGet(ipnsName string, priKey crsa.PrivateKey) []byte {
	pth := bv.iPFS.NameResolve(ipnsName)
	file := bv.iPFS.FileGet(pth)
	return []byte(rsa.Decrypt(string(file), priKey))
}
