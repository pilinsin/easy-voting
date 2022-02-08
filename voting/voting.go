package voting

import (
	"github.com/pilinsin/util"
	"github.com/pilinsin/ipfs-util"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	vutil "github.com/pilinsin/easy-voting/voting/util"
	viface "github.com/pilinsin/easy-voting/voting/interface"
	module "github.com/pilinsin/easy-voting/voting/modules"
)

func NewVoting(vCfgCid string, identity *rutil.UserIdentity, is *ipfs.IPFS) (viface.IVoting, error) {
	vCfg, err := vutil.ConfigFromCid(vCfgCid, is)
	if err != nil {
		return nil, err
	}

	switch vCfg.VType() {
	case vutil.Single:
		return module.NewSingleVoting(vCfgCid, identity, is), nil
	case vutil.Block:
		return module.NewBlockVoting(vCfgCid, identity, is), nil
	case vutil.Approval:
		return module.NewApprovalVoting(vCfgCid, identity, is), nil
	case vutil.Range:
		return module.NewRangeVoting(vCfgCid, identity, is), nil
	case vutil.Preference:
		return module.NewPreferenceVoting(vCfgCid, identity, is), nil
	case vutil.Cumulative:
		return module.NewCumulativeVoting(vCfgCid, identity, is), nil
	default:
		return nil, util.NewError("invalid VType")
	}
}
