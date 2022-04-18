package votingutil

import(
	"github.com/pilinsin/util"
)

type VoteInt map[string]int
func (vi VoteInt) Marshal() []byte{
	m, _ := util.Marshal(vi)
	return m
}
func (vi *VoteInt) Unmarshal(m []byte) error{
	return util.Unmarshal(m, vi)
}

type VoteResult map[string]map[string]int
