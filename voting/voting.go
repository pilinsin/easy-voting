package voting

import (
	"errors"
	//i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
	//crdt "github.com/pilinsin/p2p-verse/crdt"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
	evutil "github.com/pilinsin/easy-voting/util"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	vutil "github.com/pilinsin/easy-voting/voting/util"
	viface "github.com/pilinsin/easy-voting/voting/interface"
	module "github.com/pilinsin/easy-voting/voting/modules"
)

type votingWithIpfs struct{
	viface.IVoting
	is ipfs.Ipfs
}
func (v *votingWithIpfs) Close(){
	IVoting.Close()
	is.Close()
}

func NewVoting(vCfgAddr, idStr string) (viface.IVoting, error) {
	bAddr, vCfgCid, err := evutil.ParseConfigAddr(vCfgAddr)
	if err != nil{return nil, err}

	ipfsDir, save := parseIdStr(idStr)
	is, err := evutil.NewIpfs(bAddr, ipfsDir, save)
	if err != nil{return nil, err}
	vCfg := &vutil.Config{}
	if err := vCfg.FromCid(vCfgCid, is); err != nil{return nil, err}

	var v viface.IVoting
	switch vCfg.Type {
	case vutil.Single:
		v = module.NewSingleVoting(vCfg, idStr, bAddr)
	case vutil.Block:
		v = module.NewBlockVoting(vCfg, idStr, bAddr)
	case vutil.Approval:
		v = module.NewApprovalVoting(vCfg, idStr, bAddr)
	case vutil.Range:
		v = module.NewRangeVoting(vCfg, idStr, bAddr)
	case vutil.Preference:
		v = module.NewPreferenceVoting(vCfg, idStr, bAddr)
	case vutil.Cumulative:
		v = module.NewCumulativeVoting(vCfg, idStr, bAddr)
	default:
		return nil, errors.New("invalid VType")
	}

	return &votingWithIpfs{v, is}, nil
}

func parseIdStr(idStr string) (string, bool){
	mi := &vutil.ManIdentity{}
	if err := mi.FromString(idStr); err == nil{
		return mi.IpfsDir, true
	}
	
	return pv.RandString(8), false
}