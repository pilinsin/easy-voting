package votingutil

import(
	"github.com/pilinsin/util"
)

type result struct{
	res map[string]map[string]int
	nVoted int
	nVoters int
}
func NewResult(res map[string]map[string]int, nVoted, nVoters int) *result{
	return &result{res, nVoted, nVoters}
}
func (r *result) Marshal() []byte{
	mResult := &struct{
		Res map[string]map[string]int
		NVoted int
		NVoters int
	}{r.res, r.nVoted, r.nVoters}
	m, _ := util.Marshal(mResult)
	return m
}
func UnmarshalResult(m []byte) (*result, error){
	mResult := &struct{
		Res map[string]map[string]int
		NVoted int
		NVoters int
	}{}
	err := util.Unmarshal(m, mResult)

	return &result{mResult.Res, mResult.NVoted, mResult.NVoters}, err
}