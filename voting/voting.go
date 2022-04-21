package voting

import (
	"context"
	"errors"
	"path/filepath"

	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
	evutil "github.com/pilinsin/easy-voting/util"
	vutil "github.com/pilinsin/easy-voting/voting/util"
	viface "github.com/pilinsin/easy-voting/voting/interface"
	module "github.com/pilinsin/easy-voting/voting/modules"
)

type votingWithIpfs struct{
	viface.IVoting
	is ipfs.Ipfs
}
func (v *votingWithIpfs) Close(){
	v.IVoting.Close()
	v.is.Close()
}

func NewVoting(ctx context.Context, vCfgAddr, idStr string) (viface.IVoting, error) {
	bAddr, vCfgCid, err := evutil.ParseConfigAddr(vCfgAddr)
	if err != nil{return nil, err}

	ipfsDir, storeDir, save := parseIdStr(idStr)
	is, err := evutil.NewIpfs(i2p.NewI2pHost, bAddr, ipfsDir, save)
	if err != nil{return nil, err}
	vCfg := &vutil.Config{}
	if err := vCfg.FromCid(vCfgCid, is); err != nil{return nil, err}

	var v viface.IVoting
	switch vCfg.Type {
	case vutil.Single:
		v, err = module.NewSingleVoting(ctx, vCfg, idStr, storeDir, bAddr, save)
	case vutil.Block:
		v, err = module.NewBlockVoting(ctx, vCfg, idStr, storeDir, bAddr, save)
	case vutil.Approval:
		v, err = module.NewApprovalVoting(ctx, vCfg, idStr, storeDir, bAddr, save)
	case vutil.Range:
		v, err = module.NewRangeVoting(ctx, vCfg, idStr, storeDir, bAddr, save)
	case vutil.Preference:
		v, err = module.NewPreferenceVoting(ctx, vCfg, idStr, storeDir, bAddr, save)
	case vutil.Cumulative:
		v, err = module.NewCumulativeVoting(ctx, vCfg, idStr, storeDir, bAddr, save)
	default:
		return nil, errors.New("invalid VType")
	}
	if err != nil{
		is.Close()
		return nil, err
	}

	return &votingWithIpfs{v, is}, nil
}

func parseIdStr(idStr string) (string, string, bool){
	mi := &vutil.ManIdentity{}
	if err := mi.FromString(idStr); err == nil{
		return mi.IpfsDir, mi.StoreDir, true
	}

	baseDir := pv.RandString(8)
	ipfsDir := filepath.Join(baseDir, pv.RandString(8))
	storeDir := filepath.Join(baseDir, pv.RandString(8))
	return ipfsDir, storeDir, false
}
