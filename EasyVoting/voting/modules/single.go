package votingmodule

import (
	"log"

	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
	viface "EasyVoting/voting/interface"
	vutil "EasyVoting/voting/util"
)

type singleVoting struct {
	voting
}

func NewSingleVoting(vCfgCid string, identity *rutil.UserIdentity, is *ipfs.IPFS) viface.IVoting {
	sv := &singleVoting{}
	sv.init(vCfgCid, identity, is)
	return sv
}

type singleForm struct {
	widget.Select
	vi    vutil.VoteInt
	defVI func() vutil.VoteInt
}

func newSingleForm(options []string, defVI func() vutil.VoteInt) *singleForm {
	sf := &singleForm{
		Select: widget.Select{
			Options:     options,
			PlaceHolder: "(Select one)",
		},
		vi:    defVI(),
		defVI: defVI,
	}
	sf.OnChanged = sf.onChanged
	sf.ExtendBaseWidget(sf)
	return sf
}
func (sf *singleForm) onChanged(selectNG string) {
	sf.vi = sf.defVI()
	sf.vi[selectNG] = 1
	log.Println(sf.vi)
}
func (sf *singleForm) VoteInt() vutil.VoteInt {
	return sf.vi
}
func (sv *singleVoting) NewVotingForm(ngs []string) viface.IVotingForm {
	return newSingleForm(ngs, sv.newDefaultVoteInt)
}
func (sv *singleVoting) newDefaultVoteInt() vutil.VoteInt {
	vi := make(vutil.VoteInt)
	for _, ng := range sv.candNameGroups() {
		vi[ng] = 0
	}
	return vi
}

func (sv *singleVoting) isValidData(vi vutil.VoteInt) bool {
	if !sv.isCandsMatch(vi) {
		return false
	}

	numTrue := 0
	for _, vote := range vi {
		if vote > 0 {
			numTrue++
		}
	}
	return numTrue == 1
}

func (sv *singleVoting) Type() string {
	return "singlevoting"
}

func (sv *singleVoting) Vote(data vutil.VoteInt) error {
	if sv.isValidData(data) {
		return sv.baseVote(data)
	} else {
		return util.NewError("invalid vote")
	}
}
func (sv singleVoting) GetMyVote() (vutil.VoteInt, error) {
	vi, err := sv.baseGetMyVote()
	if err != nil {
		return sv.newDefaultVoteInt(), err
	} else if vi != nil && sv.isValidData(*vi) {
		return *vi, nil
	} else {
		return sv.newDefaultVoteInt(), util.NewError("invalid vote")
	}
}

func (sv singleVoting) newResult() map[string]map[string]int {
	result := make(map[string]map[string]int, len(sv.cands))
	for _, name := range sv.candNameGroups() {
		result[name] = map[string]int{"n_votes": 0}
	}
	return result
}
func (sv singleVoting) addVote2Result(vi vutil.VoteInt, result map[string]map[string]int) map[string]map[string]int {
	for k, v := range vi {
		if v > 0 {
			result[k]["n_votes"]++
		}
	}
	return result
}
func (sv singleVoting) Count() (map[string]map[string]int, int, int, error) {
	result := sv.newResult()

	viChan, nVoters, err := sv.baseGetVotes()
	if err != nil {
		return make(map[string]map[string]int), -1, -1, err
	}

	nVoted := 0
	for vi := range viChan {
		if sv.isValidData(*vi) {
			result = sv.addVote2Result(*vi, result)
			nVoted++
		}
	}
	return result, nVoted, nVoters, nil
}
