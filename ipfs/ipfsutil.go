package ipfs

import (
	"crypto/rand"
	"errors"
	"os"

	files "github.com/ipfs/go-ipfs-files"
	p2pcrypt "github.com/libp2p/go-libp2p-core/crypto"

	"EasyVoting/util"
)

func bootstrapNodes() []string {
	return []string{
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",

		"/ip4/138.201.67.219/tcp/4001/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
		"/ip4/138.201.67.219/udp/4001/quic/p2p/QmUd6zHcbkbcs7SMxwLs48qZVX3vpcM8errYS7xEczwRMA",
		"/ip4/138.201.67.220/tcp/4001/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
		"/ip4/138.201.67.220/udp/4001/quic/p2p/QmNSYxZAiJHeLdkBg38roksAR9So7Y5eojks1yjEcUtZ7i",
		"/ip4/138.201.68.74/tcp/4001/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
		"/ip4/138.201.68.74/udp/4001/quic/p2p/QmdnXwLrC8p1ueiq2Qya8joNvk3TVVDAut7PrikmZwubtR",
		"/ip4/94.130.135.167/tcp/4001/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
		"/ip4/94.130.135.167/udp/4001/quic/p2p/QmUEMvxS2e7iDrereVYc5SWPauXPyNwxcy9BXZrC1QTcHE",
	}
}

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

type KeyFile struct {
	keyFile p2pcrypt.PrivKey
}

func GenKeyFile() KeyFile {
	priv, _, err := p2pcrypt.GenerateEd25519Key(rand.Reader)
	util.CheckError(err)
	return KeyFile{priv}
}

func (kf KeyFile) Equals(kf2 KeyFile) bool {
	return kf.keyFile.Equals(kf2.keyFile)
}

func (kf KeyFile) Marshal() []byte {
	kb, err := kf.keyFile.Bytes()
	util.CheckError(err)
	return kb
}

func UnmarshalKeyFile(b []byte) KeyFile {
	kFile, err := p2pcrypt.UnmarshalPrivateKey(b)
	util.CheckError(err)

	return KeyFile{kFile}
}
