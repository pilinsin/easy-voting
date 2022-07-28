package registrationutil

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	evutil "github.com/pilinsin/easy-voting/util"
	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
	crdt "github.com/pilinsin/p2p-verse/crdt"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
	"github.com/pilinsin/util"
	hash "github.com/pilinsin/util/hash"

	pb "github.com/pilinsin/easy-voting/registration/util/pb"
	proto "google.golang.org/protobuf/proto"
)

type RegistrationStores struct {
	Is  ipfs.Ipfs
	Uhm crdt.IHashStore
}

func (rs *RegistrationStores) Close() {
	if rs.Is != nil {
		rs.Is.Close()
	}
	if rs.Uhm != nil {
		rs.Uhm.Close()
	}
}

type Config struct {
	Title   string
	Salt1   string
	Salt2   []byte
	UhmAddr string
	Labels  []string
}

func NewConfig(title string, userDataset <-chan []string, userDataLabels []string, bAddr string) (string, *RegistrationStores, error) {
	salt2 := hash.HashWithSize([]byte(title), util.GenRandomBytes(30), 32)
	cfg := &Config{
		Title:  title,
		Salt1:  title + util.GenUniqueID(30, 30),
		Salt2:  salt2,
		Labels: userDataLabels,
	}

	uhHashes := make(chan string)
	go func() {
		defer close(uhHashes)
		for userData := range userDataset {
			userHash := NewUserHash(cfg.Salt1, userData...)
			uhHashes <- crdt.MakeHashKey(userHash, salt2)
		}
	}()

	bootstraps := pv.AddrInfosFromString(bAddr)
	baseDir := evutil.BaseDir("registration", "setup")

	storeDir := filepath.Join(baseDir, "store")
	os.RemoveAll(storeDir)
	v := crdt.NewVerse(i2p.NewI2pHost, storeDir, false, bootstraps...)
	uhm, err := v.NewStore(pv.RandString(8), "hash", &crdt.StoreOpts{Salt: salt2})
	if err != nil {
		return "", nil, err
	}
	uhm, err = v.NewAccessStore(uhm, uhHashes)
	if err != nil {
		return "", nil, err
	}
	cfg.UhmAddr = uhm.Address()

	ipfsDir := filepath.Join(baseDir, "ipfs")
	os.RemoveAll(ipfsDir)
	is, err := evutil.NewIpfs(i2p.NewI2pHost, ipfsDir, false, bootstraps)
	if err != nil {
		uhm.Close()
		return "", nil, err
	}
	cid, err := cfg.toCid(is)
	if err != nil {
		uhm.Close()
		is.Close()
		return "", nil, err
	}

	rs := &RegistrationStores{
		Is:  is,
		Uhm: uhm.(crdt.IHashStore),
	}

	return "r/" + cid, rs, nil
}

func (cfg Config) Marshal() []byte {
	pbCfg := &pb.Config{
		Title:   cfg.Title,
		Salt1:   cfg.Salt1,
		Salt2:   cfg.Salt2,
		Labels:  cfg.Labels,
		UhmAddr: cfg.UhmAddr,
	}
	m, _ := proto.Marshal(pbCfg)
	return m
}
func (cfg *Config) Unmarshal(m []byte) error {
	pbCfg := &pb.Config{}
	if err := proto.Unmarshal(m, pbCfg); err != nil {
		return err
	}

	cfg.Title = pbCfg.Title
	cfg.Salt1 = pbCfg.Salt1
	cfg.Salt2 = pbCfg.Salt2
	cfg.Labels = pbCfg.Labels
	cfg.UhmAddr = pbCfg.UhmAddr
	return nil
}

func (cfg *Config) toCid(is ipfs.Ipfs) (string, error) {
	return is.Add(cfg.Marshal(), time.Second*10)
}
func (cfg *Config) FromCid(rCfgCid string, is ipfs.Ipfs) error {
	m, err := is.Get(rCfgCid, time.Second*5)
	if err != nil {
		return errors.New("get from rCfgCid error")
	}

	if err := cfg.Unmarshal(m); err != nil {
		return errors.New("unmarshal rCfg error")
	}
	return nil
}
