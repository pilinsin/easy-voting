package voting

import (
	"errors"
	"path/filepath"

	evutil "github.com/pilinsin/easy-voting/util"
	viface "github.com/pilinsin/easy-voting/voting/interface"
	module "github.com/pilinsin/easy-voting/voting/modules"
	vutil "github.com/pilinsin/easy-voting/voting/util"
	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
	crdt "github.com/pilinsin/p2p-verse/crdt"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
)

type votingWithAddress struct {
	viface.ITypedVoting
	addr string
}

func (v *votingWithAddress) Close() {
	v.ITypedVoting.Close()
}
func (v *votingWithAddress) Address() string {
	return v.addr
}

func NewVoting(vCfgAddr, baseDir string) (viface.IVoting, error) {
	ipfsDir := filepath.Join(baseDir, "ipfs")
	storeDir := filepath.Join(baseDir, "store")
	save := true

	bAddr, vCfgCid, err := evutil.ParseConfigAddr(vCfgAddr)
	if err != nil {
		return nil, err
	}
	bootstraps := pv.AddrInfosFromString(bAddr)

	is, err := evutil.NewIpfs(i2p.NewI2pHost, ipfsDir, save, bootstraps)
	if err != nil {
		return nil, err
	}

	vCfg := &vutil.Config{}
	if err := vCfg.FromCid(vCfgCid, is); err != nil {
		is.Close()
		return nil, err
	}

	stInfo := [][2]string{
		{vCfg.HkmAddr, "signature"},
		{vCfg.IvmAddr, "updatableSignature"},
	}
	stores, err := evutil.NewStore(i2p.NewI2pHost, stInfo, storeDir, save, bootstraps)
	if err != nil {
		is.Close()
		return nil, err
	}
	if len(stores) < 2 {
		is.Close()
		for _, st := range stores{
			st.Close()
		}
		return nil, errors.New("too few stores loaded")
	}
	hkm, tmp := stores[0], stores[1]
	ivm := tmp.(crdt.IUpdatableSignatureStore)

	var v viface.ITypedVoting
	switch vCfg.Type {
	case vutil.Single:
		v, err = module.NewSingleVoting(vCfg, is, hkm, ivm)
	case vutil.Block:
		v, err = module.NewBlockVoting(vCfg, is, hkm, ivm)
	case vutil.Approval:
		v, err = module.NewApprovalVoting(vCfg, is, hkm, ivm)
	case vutil.Range:
		v, err = module.NewRangeVoting(vCfg, is, hkm, ivm)
	case vutil.Preference:
		v, err = module.NewPreferenceVoting(vCfg, is, hkm, ivm)
	case vutil.Cumulative:
		v, err = module.NewCumulativeVoting(vCfg, is, hkm, ivm)
	default:
		return nil, errors.New("invalid VType")
	}
	if err != nil {
		is.Close()
		hkm.Close()
		ivm.Close()
		return nil, err
	}

	return &votingWithAddress{v, vCfgAddr}, nil
}


func NewVotingWithStores(vCfgAddr string, is ipfs.Ipfs, hkm crdt.IStore, ivm crdt.IUpdatableSignatureStore) (viface.IVoting, error){
	_, vCfgCid, err := evutil.ParseConfigAddr(vCfgAddr)
	if err != nil {
		return nil, err
	}
	vCfg := &vutil.Config{}
	if err := vCfg.FromCid(vCfgCid, is); err != nil {
		return nil, err
	}

	var v viface.ITypedVoting
	switch vCfg.Type {
	case vutil.Single:
		v, err = module.NewSingleVoting(vCfg, is, hkm, ivm)
	case vutil.Block:
		v, err = module.NewBlockVoting(vCfg, is, hkm, ivm)
	case vutil.Approval:
		v, err = module.NewApprovalVoting(vCfg, is, hkm, ivm)
	case vutil.Range:
		v, err = module.NewRangeVoting(vCfg, is, hkm, ivm)
	case vutil.Preference:
		v, err = module.NewPreferenceVoting(vCfg, is, hkm, ivm)
	case vutil.Cumulative:
		v, err = module.NewCumulativeVoting(vCfg, is, hkm, ivm)
	default:
		return nil, errors.New("invalid VType")
	}
	if err != nil{return nil, err}

	return &votingWithAddress{v, vCfgAddr}, nil
}