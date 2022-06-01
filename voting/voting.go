package voting

import (
	"context"
	"errors"
	"path/filepath"

	evutil "github.com/pilinsin/easy-voting/util"
	viface "github.com/pilinsin/easy-voting/voting/interface"
	module "github.com/pilinsin/easy-voting/voting/modules"
	vutil "github.com/pilinsin/easy-voting/voting/util"
	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
)

type votingWithIpfs struct {
	viface.ITypedVoting
	addr string
	is   ipfs.Ipfs
}

func (v *votingWithIpfs) Close() {
	v.ITypedVoting.Close()
	v.is.Close()
}
func (v *votingWithIpfs) Address() string {
	return v.addr
}

func NewVoting(ctx context.Context, vCfgAddr, baseDir string) (viface.IVoting, error) {
	bAddr, vCfgCid, err := evutil.ParseConfigAddr(vCfgAddr)
	if err != nil {
		return nil, err
	}
	bootstraps := pv.AddrInfosFromString(bAddr)

	ipfsDir := filepath.Join(baseDir, "ipfs")
	storeDir := filepath.Join(baseDir, "store")
	save := true

	is, err := evutil.NewIpfs(i2p.NewI2pHost, ipfsDir, save, bootstraps)
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
		v, err = module.NewSingleVoting(ctx, vCfg, storeDir, bootstraps, save)
	case vutil.Block:
		v, err = module.NewBlockVoting(ctx, vCfg, storeDir, bootstraps, save)
	case vutil.Approval:
		v, err = module.NewApprovalVoting(ctx, vCfg, storeDir, bootstraps, save)
	case vutil.Range:
		v, err = module.NewRangeVoting(ctx, vCfg, storeDir, bootstraps, save)
	case vutil.Preference:
		v, err = module.NewPreferenceVoting(ctx, vCfg, storeDir, bootstraps, save)
	case vutil.Cumulative:
		v, err = module.NewCumulativeVoting(ctx, vCfg, storeDir, bootstraps, save)
	default:
		return nil, errors.New("invalid VType")
	}
	if err != nil {
		is.Close()
		return nil, err
	}

	return &votingWithIpfs{v, vCfgAddr, is}, nil
}
