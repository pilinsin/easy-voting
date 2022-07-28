package votingmodule

import (
	"errors"
	"log"

	"fyne.io/fyne/v2/widget"

	viface "github.com/pilinsin/easy-voting/voting/interface"
	vutil "github.com/pilinsin/easy-voting/voting/util"
	crdt "github.com/pilinsin/p2p-verse/crdt"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
)

type approvalVoting struct {
	voting
}

func NewApprovalVoting(vCfg *vutil.Config, is ipfs.Ipfs, hkm crdt.ISignatureStore, ivm crdt.IUpdatableSignatureStore) (viface.ITypedVoting, error) {
	av := &approvalVoting{}
	if err := av.init(vCfg, is, hkm, ivm); err != nil {
		return nil, err
	}
	return av, nil
}

type approvalForm struct {
	widget.Form
	vi vutil.VoteInt
}

func newApprovalForm(options []string, vi vutil.VoteInt) *approvalForm {
	af := &approvalForm{
		Form: widget.Form{},
		vi:   vi,
	}

	for _, opt := range options {
		chk := &widget.Check{
			DisableableWidget: widget.DisableableWidget{},
			Text:              opt,
		}
		chk.OnChanged = func(isChecked bool) {
			if isChecked {
				af.vi[chk.Text] = 1
			} else {
				af.vi[chk.Text] = 0
			}
			log.Println(af.vi)
		}
		chk.ExtendBaseWidget(chk)
		af.Items = append(af.Items, widget.NewFormItem("", chk))
	}

	af.ExtendBaseWidget(af)
	return af
}
func (af *approvalForm) VoteInt() vutil.VoteInt {
	return af.vi
}
func (av *approvalVoting) NewVotingForm(ngs []string) viface.IVotingForm {
	return newApprovalForm(ngs, av.newDefaultVoteInt())
}

func (av *approvalVoting) newDefaultVoteInt() vutil.VoteInt {
	vi := make(vutil.VoteInt)
	for _, ng := range av.candNameGroups() {
		vi[ng] = 0
	}

	return vi
}

func (av *approvalVoting) isValidData(vi vutil.VoteInt) bool {
	return av.isCandsMatch(vi)
}

func (av *approvalVoting) Type() string {
	return "approvalvoting"
}

func (av *approvalVoting) Vote(data vutil.VoteInt) error {
	if av.isValidData(data) {
		return av.baseVote(data)
	} else {
		return errors.New("invalid vote")
	}
}
func (av approvalVoting) GetMyVote() (*vutil.VoteInt, error) {
	vi, err := av.baseGetMyVote()
	if err != nil {
		return nil, err
	} else if vi != nil && av.isValidData(*vi) {
		return vi, nil
	} else {
		return nil, errors.New("invalid vote")
	}
}

func (av approvalVoting) newResult() vutil.VoteResult {
	result := make(vutil.VoteResult, len(av.cands))
	for _, name := range av.candNameGroups() {
		result[name] = map[string]int{"n_votes": 0}
	}
	return result
}
func (av approvalVoting) addVoteToResult(vi vutil.VoteInt, result vutil.VoteResult) vutil.VoteResult {
	for k, v := range vi {
		if v > 0 {
			result[k]["n_votes"]++
		}
	}
	return result
}
func (av approvalVoting) GetResult() (*vutil.VoteResult, int, int, error) {
	result := av.newResult()

	viChan, nVoters, err := av.baseGetVotes()
	if err != nil {
		return nil, -1, -1, err
	}

	nVoted := 0
	for vi := range viChan {
		if av.isValidData(*vi) {
			result = av.addVoteToResult(*vi, result)
			nVoted++
		}
	}
	return &result, nVoters, nVoted, nil
}
