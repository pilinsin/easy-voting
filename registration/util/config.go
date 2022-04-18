package registrationutil

import (
	"errors"
	"path/filepath"

	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
	crdt "github.com/pilinsin/p2p-verse/crdt"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
	evutil "github.com/pilinsin/easy-voting/util"
	
	pb "github.com/pilinsin/easy-voting/registration/util/pb"
	proto "google.golang.org/protobuf/proto"
)

type Config struct{
	Title string
	Salt1 string
	Salt2 []byte
	UhmAddr string
	Labels []string
}
func NewConfig(title string, userDataset <-chan []string, userDataLabels []string, bAddr string) (string, string, error) {
	salt2 := crypto.HashWithSize([]byte(title), util.GenRandomBytes(30), 32)
	cfg := &Config{
		Title:          title,
		Salt1:          title + util.GenUniqueID(30, 30),
		Salt2:          salt2,
		Labels: userDataLabels,
	}

	uhHashes := make(chan string)
	go func(){
		defer close(uhHashes)
		for userData := range userDataset {
			userHash := NewUserHash(cfg.Salt1, userData...)
			uhHashes <- crdt.MakeHashKey(userHash, salt2)
		}
	}()

	bootstraps := pv.AddrInfosFromString(bAddr)
	storeDir := filepath.Join(title, pv.RandString(8))
	v := crdt.NewVerse(i2p.NewI2pHost, storeDir, true, false, bootstraps...)
	ac, err := v.NewAccessController(pv.RandString(8), uhHashes)
	if err != nil{return "", "", err}
	uhm, err := v.NewStore(pv.RandString(8), "hash", &crdt.StoreOpts{Salt:salt2, Ac:ac})
	if err != nil{return "", "", err}
	defer uhm.Close()

	cfg.UhmAddr = uhm.Address()

	ipfsDir := filepath.Join(title, pv.RandString(8))
	is, err := evutil.NewIpfs(i2p.NewI2pHost, bAddr, ipfsDir, true)
	if err != nil{return "", "", err}
	defer is.Close()
	cid, err := cfg.toCid(is)
	if err != nil{return "", "", err}

	mi := &ManIdentity{ipfsDir, storeDir}

	return "r/"+bAddr+"/"+cid, mi.toString(), nil
}

func (cfg Config) Marshal() []byte{
	pbCfg := &pb.Config{
		Title: cfg.Title,
		Salt1: cfg.Salt1,
		Salt2: cfg.Salt2,
		Labels: cfg.Labels,
		UhmAddr: cfg.UhmAddr,
	}
	m, _ := proto.Marshal(pbCfg)
	return m
}
func (cfg *Config) Unmarshal(m []byte) error{
	pbCfg := &pb.Config{}
	if err := proto.Unmarshal(m, pbCfg); err != nil{return err}

	cfg.Title = pbCfg.Title
	cfg.Salt1 = pbCfg.Salt1
	cfg.Salt2 = pbCfg.Salt2
	cfg.Labels = pbCfg.Labels
	cfg.UhmAddr = pbCfg.UhmAddr
	return nil
}

func (cfg *Config) toCid(is ipfs.Ipfs) (string, error){
	return is.Add(cfg.Marshal())
}
func (cfg *Config) FromCid(rCfgCid string, is ipfs.Ipfs) error {
	m, err := is.Get(rCfgCid)
	if err != nil {
		return errors.New("get from rCfgCid error")
	}

	if err := cfg.Unmarshal(m); err != nil {
		return errors.New("unmarshal rCfg error")
	}
	return nil
}
