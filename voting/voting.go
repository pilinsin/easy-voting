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
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
)

type votingWithIpfs struct {
	viface.IVoting
	is        ipfs.Ipfs
}

func (v *votingWithIpfs) Close() {
	v.IVoting.Close()
	v.is.Close()
}

func NewVoting(ctx context.Context, vCfgAddr, baseDir string) (viface.IVoting, error) {
	bAddr, vCfgCid, err := evutil.ParseConfigAddr(vCfgAddr)
	if err != nil {
		return nil, err
	}

	ipfsDir := filepath.Join("stores", baseDir, "ipfs")
	storeDir := filepath.Join("stores", baseDir, "store")
	save := true

	is, err := evutil.NewIpfs(i2p.NewI2pHost, bAddr, ipfsDir, save)
	if err != nil {
		return nil, err
	}
	vCfg := &vutil.Config{}
	if err := vCfg.FromCid(vCfgCid, is); err != nil {
		return nil, err
	}

	var v viface.IVoting
	switch vCfg.Type {
	case vutil.Single:
		v, err = module.NewSingleVoting(ctx, vCfg, storeDir, bAddr, save)
	case vutil.Block:
		v, err = module.NewBlockVoting(ctx, vCfg, storeDir, bAddr, save)
	case vutil.Approval:
		v, err = module.NewApprovalVoting(ctx, vCfg, storeDir, bAddr, save)
	case vutil.Range:
		v, err = module.NewRangeVoting(ctx, vCfg, storeDir, bAddr, save)
	case vutil.Preference:
		v, err = module.NewPreferenceVoting(ctx, vCfg, storeDir, bAddr, save)
	case vutil.Cumulative:
		v, err = module.NewCumulativeVoting(ctx, vCfg, storeDir, bAddr, save)
	default:
		return nil, errors.New("invalid VType")
	}
	if err != nil {
		is.Close()
		return nil, err
	}

	return &votingWithIpfs{v, is}, nil
}
