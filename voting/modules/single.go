package votingmodule

import (
	"context"
	"errors"
	"log"

	"fyne.io/fyne/v2/widget"
	peer "github.com/libp2p/go-libp2p-core/peer"

	viface "github.com/pilinsin/easy-voting/voting/interface"
	vutil "github.com/pilinsin/easy-voting/voting/util"
)

type singleVoting struct {
	voting
}

func NewSingleVoting(ctx context.Context, vCfg *vutil.Config, storeDir string, bs []peer.AddrInfo, save bool) (viface.ITypedVoting, error) {
	sv := &singleVoting{}
	if err := sv.init(ctx, vCfg, storeDir, bs, save); err != nil {
		return nil, err
	}
	return sv, nil
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
		return errors.New("invalid vote")
	}
}
func (sv singleVoting) GetMyVote() (*vutil.VoteInt, error) {
	vi, err := sv.baseGetMyVote()
	if err != nil {
		return nil, err
	} else if vi != nil && sv.isValidData(*vi) {
		return vi, nil
	} else {
		return nil, errors.New("invalid vote")
	}
}

func (sv singleVoting) newResult() vutil.VoteResult {
	result := make(vutil.VoteResult, len(sv.cands))
	for _, name := range sv.candNameGroups() {
		result[name] = map[string]int{"n_votes": 0}
	}
	return result
}
func (sv singleVoting) addVoteToResult(vi vutil.VoteInt, result vutil.VoteResult) vutil.VoteResult {
	for k, v := range vi {
		if v > 0 {
			result[k]["n_votes"]++
		}
	}
	return result
}

func (sv singleVoting) GetResult() (*vutil.VoteResult, int, int, error) {
	result := sv.newResult()

	viChan, nVoters, err := sv.baseGetVotes()
	if err != nil {
		return nil, -1, -1, err
	}

	nVoted := 0
	for vi := range viChan {
		if sv.isValidData(*vi) {
			result = sv.addVoteToResult(*vi, result)
			nVoted++
		}
	}
	return &result, nVoters, nVoted, nil
}
