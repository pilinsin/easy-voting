package votingutil

import (
	"EasyVoting/ipfs"
	"EasyVoting/util"
	"EasyVoting/util/ecies"
)

type resultMap struct {
	rm      *ipfs.ReccurentMap
	tInfo   *util.TimeInfo
	nVoters int
}

func NewResultMap(capacity int, ivm *idVotingMap, manPriKey *ecies.PriKey, is *ipfs.IPFS) (*resultMap, error) {
	resMap := &resultMap{
		rm:      ipfs.NewReccurentMap(capacity),
		tInfo:   ivm.tInfo,
		nVoters: ivm.Len(is),
	}
	for kv := range ivm.NextKeyValue(is) {
		uvHash := kv.Key()
		vb, _ := kv.Value()
		msv, err := vb.GetVote(manPriKey)
		if err != nil {
			return nil, err
		}
		resMap.rm.Append(uvHash, msv, is)
	}
	return resMap, nil
}
func (resMap resultMap) NumVoters() int { return resMap.nVoters }
func (resMap resultMap) VerifyVotes(ivm *idVotingMap, is *ipfs.IPFS) bool {
	for kv := range resMap.rm.NextKeyValue(is) {
		sv, err := UnmarshalSignedVote(kv.Value())
		if err != nil {
			return false
		}
		_, verfKey, ok := ivm.ContainHash(UidVidHash(kv.Key()), is)
		if !ok {
			return false
		}
		if ok := sv.Verify(verfKey); !ok {
			return false
		}
	}
	return true
}
func (resMap resultMap) Next(is *ipfs.IPFS) <-chan *VoteInt {
	ch := make(chan *VoteInt)
	go func() {
		defer close(ch)
		for m := range resMap.rm.Next(is) {
			if sv, err := UnmarshalSignedVote(m); err == nil {
				if vi, err := sv.Vote(resMap.tInfo); err == nil {
					ch <- vi
				}
			}
		}
	}()
	return ch
}
func (resMap resultMap) Marshal() []byte {
	mResMap := &struct {
		Mrm      []byte
		TimeInfo *util.TimeInfo
		N        int
	}{resMap.rm.Marshal(), resMap.tInfo, resMap.nVoters}
	m, _ := util.Marshal(mResMap)
	return m
}
func UnmarshalResultMap(m []byte) (*resultMap, error) {
	mResMap := &struct {
		Mrm      []byte
		TimeInfo *util.TimeInfo
		N        int
	}{}
	if err := util.Unmarshal(m, mResMap); err != nil {
		return nil, err
	}

	resMap := &resultMap{}
	if err := resMap.rm.Unmarshal(mResMap.Mrm); err != nil {
		return nil, err
	}
	resMap.tInfo = mResMap.TimeInfo
	resMap.nVoters = mResMap.N
	return resMap, nil
}
func ResultMapFromName(resMapName string, is *ipfs.IPFS) (*resultMap, error) {
	m, err := ipfs.FromName(resMapName, is)
	if err != nil {
		return nil, err
	}
	return UnmarshalResultMap(m)
}
