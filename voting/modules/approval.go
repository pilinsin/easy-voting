package votingmodule

import (
	"log"

	"fyne.io/fyne/v2/widget"

	"github.com/pilinsin/ipfs-util"
	"github.com/pilinsin/util"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	vutil "github.com/pilinsin/easy-voting/voting/util"
	viface "github.com/pilinsin/easy-voting/voting/interface"
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
func (av approvalVoting) GetMyVote() (string, error) {
	vi, err := av.baseGetMyVote()
	if err != nil {
		return "", err
	} else if vi != nil && av.isValidData(*vi) {
		return ipfs.ToCidWithAdd(vi.Marshal(), av.is), nil
	} else {
		return "", util.NewError("invalid vote")
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
func (av approvalVoting) CountMyResult() (string, error){
	return av.countResult(av.baseGetMyVotes)
}
func (av approvalVoting) CountManResult() (string, error){
	return av.countResult(av.baseGetVotes)
}
func (av approvalVoting) countResult(gvf viface.GetVotesFunc) (string, error) {
	result := av.newResult()

	viChan, nVoters, err := gvf()
	if err != nil {
		return "", err
	}

	nVoted := 0
	for vi := range viChan {
		if av.isValidData(*vi) {
			result = av.addVote2Result(*vi, result)
			nVoted++
		}
	}
	m := vutil.NewResult(result, nVoted, nVoters).Marshal()
	return ipfs.File.Add(m, av.is), nil
}

