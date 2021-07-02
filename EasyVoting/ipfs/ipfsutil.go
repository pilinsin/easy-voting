package ipfs

import(
	"os"
	"errors"
	"crypto/rand"

	files "github.com/ipfs/go-ipfs-files"
	p2pcrypt  "github.com/libp2p/go-libp2p-core/crypto"
	pb  "github.com/libp2p/go-libp2p-core/crypto/pb"
	options "github.com/ipfs/interface-go-ipfs-core/options"

	"EasyVoting/util"
)

func Str2IpfsFile(str string) files.File{
	return files.NewBytesFile([]byte(str))	
}

func IpfsFileNode2Str(node files.Node) string{
	switch node := node.(type){
	case files.File:
		return util.Reader2Str(node)
	default:
		err := errors.New("node (type: files.Node) does not have files.File type!")
		util.CheckError(err)
		return ""
	}
}
func FilePath2IpfsFile(filePath string) files.File{
	file, err := os.Open(filePath)
	util.CheckError(err)
	defer file.Close()
	
	return files.NewReaderFile(file)
}

func DirPath2IpfsFileNode(dirPath string) files.Node{
	st, err := os.Stat(dirPath)
	util.CheckError(err)
	f, err := files.NewSerialFile(dirPath, false, st)
	util.CheckError(err)

	return f
}


func KeyFileGenerate() p2pcrypt.PrivKey{
	priv, _, err := p2pcrypt.GenerateKeyPairWithReader(p2pcrypt.RSA, options.DefaultRSALen, rand.Reader)
	util.CheckError(err)
	return priv
}

func MarshalKeyFile(kFile p2pcrypt.PrivKey) []byte{
	b, err := p2pcrypt.MarshalPrivateKey(kFile)
	util.CheckError(err)
	return b
}

func UnMarshalKeyFile(b []byte) p2pcrypt.PrivKey{
	unmarshal := p2pcrypt.PrivKeyUnmarshallers[pb.KeyType_RSA]
	kFile, err := unmarshal(b)
	util.CheckError(err)

	return kFile
}