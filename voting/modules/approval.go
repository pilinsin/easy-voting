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

type approvalVoting struct {
	voting
}

func NewApprovalVoting(vCfgCid string, identity *rutil.UserIdentity, is *ipfs.IPFS) viface.IVoting {
	av := &approvalVoting{}
	av.init(vCfgCid, identity, is)
	return av
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
		return util.NewError("invalid vote")
	}
}
func (av approvalVoting) GetMyVote() (vutil.VoteInt, error) {
	vi, err := av.baseGetMyVote()
	if err != nil {
		return av.newDefaultVoteInt(), err
	} else if vi != nil && av.isValidData(*vi) {
		return *vi, nil
	} else {
		return av.newDefaultVoteInt(), util.NewError("invalid vote")
	}
}

func (av approvalVoting) newResult() map[string]map[string]int {
	result := make(map[string]map[string]int, len(av.cands))
	for _, name := range av.candNameGroups() {
		result[name] = map[string]int{"n_votes": 0}
	}
	return result
}
func (av approvalVoting) addVote2Result(vi vutil.VoteInt, result map[string]map[string]int) map[string]map[string]int {
	for k, v := range vi {
		if v > 0 {
			result[k]["n_votes"]++
		}
	}
	return result
}
func (av approvalVoting) Count() (map[string]map[string]int, int, int, error) {
	result := av.newResult()

	viChan, nVoters, err := av.baseGetVotes()
	if err != nil {
		return make(map[string]map[string]int), -1, -1, err
	}

	nVoted := 0
	for vi := range viChan {
		if av.isValidData(*vi) {
			result = av.addVote2Result(*vi, result)
			nVoted++
		}
	}
	return result, nVoted, nVoters, nil
}
