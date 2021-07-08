package ipfs

import (
	"crypto/rand"
	"errors"
	"os"

	files "github.com/ipfs/go-ipfs-files"
	p2pcrypt "github.com/libp2p/go-libp2p-core/crypto"

	"EasyVoting/util"
)

func Bytes2IpfsFile(b []byte) files.File {
	return files.NewBytesFile(b)
}

func IpfsFileNode2Bytes(node files.Node) []byte {
	switch node := node.(type) {
	case files.File:
		return util.Reader2Bytes(node)
	default:
		err := errors.New("node (type: files.Node) does not have files.File type!")
		util.CheckError(err)
		return nil
	}
}
func FilePath2IpfsFile(filePath string) files.File {
	file, err := os.Open(filePath)
	util.CheckError(err)
	defer file.Close()

	return files.NewReaderFile(file)
}

func DirPath2IpfsFileNode(dirPath string) files.Node {
	st, err := os.Stat(dirPath)
	util.CheckError(err)
	f, err := files.NewSerialFile(dirPath, false, st)
	util.CheckError(err)

	return f
}

func GenKeyFile() p2pcrypt.PrivKey {
	priv, _, err := p2pcrypt.GenerateRSAKeyPair(2048, rand.Reader)
	util.CheckError(err)
	return priv
}

func MarshalKeyFile(kFile p2pcrypt.PrivKey) []byte {
	kb, err := kFile.Bytes()
	util.CheckError(err)
	return kb
}

func UnmarshalKeyFile(b []byte) p2pcrypt.PrivKey {
	kFile, err := p2pcrypt.UnmarshalPrivateKey(b)
	util.CheckError(err)

	return kFile
}
