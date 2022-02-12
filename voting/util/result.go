package votingutil

import(
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)

type resultBox struct{
	hvtm *HashVoteMap
	manPriKey crypto.IPriKey
}
func NewResultBox(hvtm *HashVoteMap, manPriKey crypto.IPriKey) *resultBox{
	return &resultBox{hvtm, manPriKey}
}
func (res resultBox) HashVoteMap() *HashVoteMap{
	return res.hvtm
}
func (res resultBox) ManPriKey() crypto.IPriKey{
	return res.manPriKey
}
func (res *resultBox) Marshal() []byte{
	mRes := &struct{
		M []byte
		K []byte
	}{res.hvtm.Marshal(), res.manPriKey.Marshal()}
	m, _ := util.Marshal(mRes)
	return m
}
func UnmarshalResultBox(m []byte) (*resultBox, error){
	mRes := &struct{
		M []byte
		K []byte
	}{}
	if err := util.Unmarshal(m, mRes); err != nil{return nil, err}

	hvtm, err := UnmarshalHashVoteMap(mRes.M)
	if err != nil{return nil, err}
	priKey, err := crypto.UnmarshalPriKey(mRes.K)
	if err != nil{return nil, err}
	return &resultBox{hvtm, priKey}, nil
}


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