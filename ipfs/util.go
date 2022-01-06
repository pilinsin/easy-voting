package ipfs

import (
	"errors"
	"os"
	"strings"

	cid "github.com/ipfs/go-cid"
	files "github.com/ipfs/go-ipfs-files"
	path "github.com/ipfs/interface-go-ipfs-core/path"

	"EasyVoting/util"
)

func bytesToIpfsFile(b []byte) files.File {
	return files.NewBytesFile(b)
}

func ipfsFileNodeToBytes(node files.Node) ([]byte, error) {
	switch node := node.(type) {
	case files.File:
		return util.ReaderToBytes(node), nil
	default:
		err := errors.New("node (type: files.Node) does not have files.File type!")
		return nil, err
	}
}
func filePathToIpfsFile(filePath string) (files.File, error) {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		return nil, err
	} else {
		return files.NewReaderFile(file), nil
	}
}

func dirPathToIpfsFileNode(dirPath string) (files.Node, error) {
	st, err := os.Stat(dirPath)
	if err != nil {
		return nil, err
	}
	return files.NewSerialFile(dirPath, false, st)
}

func ToCid(data []byte, is *IPFS) string {
	return is.FileHash(data).Cid().String()
}
func ToCidWithAdd(data []byte, is *IPFS) string {
	return is.FileAdd(data, true).Cid().String()
}
func FromCid(cidStr string, is *IPFS) ([]byte, error) {
	cid, err := cid.Decode(cidStr)
	if err != nil {
		return nil, err
	}
	pth := path.IpfsPath(cid)
	return is.FileGet(pth)
}

func ToName(data []byte, kw string, is *IPFS) string {
	pth := is.FileAdd(data, true)
	return is.NamePublish(pth, "", kw).Name()
}
func ToNameWithKeyFile(data []byte, kf *KeyFile, is *IPFS) string {
	pth := is.FileAdd(data, true)
	return is.NamePublishWithKeyFile(pth, "", kf).Name()
}
func CidToName(cidStr string, kw string, is *IPFS) string {
	cid, err := cid.Decode(cidStr)
	if err != nil {
		return ""
	}
	pth := path.IpfsPath(cid)
	return is.NamePublish(pth, "", kw).Name()
}
func CidToNameWithKeyFile(cidStr string, kf *KeyFile, is *IPFS) string {
	cid, err := cid.Decode(cidStr)
	if err != nil {
		return ""
	}
	pth := path.IpfsPath(cid)
	return is.NamePublishWithKeyFile(pth, "", kf).Name()
}
func FromName(ipnsName string, is *IPFS) ([]byte, error) {
	pth, err := is.NameResolve(ipnsName)
	if err != nil {
		return nil, err
	} else {
		return is.FileGet(pth)
	}
}
func CidFromName(ipnsName string, is *IPFS) (string, error) {
	pth, err := is.NameResolve(ipnsName)
	if err != nil {
		return "", err
	} else {
		return strings.TrimPrefix(pth.String(), "/ipfs/"), nil
	}
}
