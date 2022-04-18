package votingmodule

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/pilinsin/util"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	vutil "github.com/pilinsin/easy-voting/voting/util"
	viface "github.com/pilinsin/easy-voting/voting/interface"
)

type rangeVoting struct {
	voting
	min int
	max int
}

func NewRangeVoting(vCfg *vutil.Config, idStr, bAddr string) viface.IVoting {
	rv := &rangeVoting{
		min: vCfg.Params.Min,
		max: vCfg.Params.Max,
	}
	rv.init(vCfg, idStr, bAddr)
	return rv
}

type rangeForm struct {
	widget.Form
	vi vutil.VoteInt
}
type rSliderWithText struct {
	widget.Slider
	Text string
}

func newRangeForm(options []string, min, max int, vi vutil.VoteInt) *rangeForm {
	rf := &rangeForm{
		Form: widget.Form{},
		vi:   vi,
	}

	for _, opt := range options {
		valLabel := widget.NewLabel("0")
		sl := &rSliderWithText{
			Slider: widget.Slider{
				Value:       0,
				Min:         float64(min),
				Max:         float64(max),
				Step:        1,
				Orientation: widget.Horizontal,
			},
			Text: opt,
		}
		sl.OnChanged = func(val float64) {
			rf.vi[sl.Text] = int(val)
			valLabel.Text = strconv.Itoa(int(sl.Value))
			log.Println(rf.vi)
		}
		sl.ExtendBaseWidget(sl)

		sl2 := container.NewGridWrap(fyne.NewSize(250.0, 20.0), sl)
		slTxt := widget.NewLabel(sl.Text)
		rf.Items = append(rf.Items, widget.NewFormItem("", container.NewHBox(sl2, valLabel, slTxt)))
	}

	rf.ExtendBaseWidget(rf)
	return rf
}
func (rf *rangeForm) VoteInt() vutil.VoteInt {
	return rf.vi
}
func (rv *rangeVoting) NewVotingForm(ngs []string) viface.IVotingForm {
	return newRangeForm(ngs, rv.min, rv.max, rv.newDefaultVoteInt())
}

func (rv *rangeVoting) newDefaultVoteInt() vutil.VoteInt {
	vi := make(vutil.VoteInt)
	for _, ng := range rv.candNameGroups() {
		vi[ng] = rv.min
	}

	return vi
}

func (rv *rangeVoting) isValidData(vi vutil.VoteInt) bool {
	if !rv.isCandsMatch(vi) {
		return false
	}

	for _, vote := range vi {
		if vote < rv.min || vote > rv.max {
			return false
		}
	}

	return true
}

func (rv *rangeVoting) Type() string {
	return "rangevoting"
}

func (rv *rangeVoting) Vote(data vutil.VoteInt) error {
	if rv.isValidData(data) {
		return rv.baseVote(data)
	} else {
		return util.NewError("invalid vote")
	}
}
func (rv rangeVoting) GetMyVote() (*vutil.VoteInt, error) {
	vi, err := rv.baseGetMyVote()
	if err != nil {
		return nil, err
	} else if vi != nil && rv.isValidData(*vi) {
		return vi, nil
	} else {
		return nil, util.NewError("invalid vote")
	}
}

func (rv rangeVoting) newResult() vutil.VoteResult {
	result := make(vutil.VoteResult, len(rv.cands))
	for _, name := range rv.candNameGroups() {
		result[name] = map[string]int{"score": 0}
	}
	return result
}
func (rv rangeVoting) addVoteToResult(vi vutil.VoteInt, result vutil.VoteResult) vutil.VoteResult {
	for k, v := range vi {
		result[k]["score"] += v
	}
	return result
}

func (rv rangeVoting) GetResult() (*vutil.VoteResult, int, error) {
	result := rv.newResult()

	viChan, nVoters, err := rv.baseGetVotes()
	if err != nil {
		return nil, -1, err
	}

	nVoted := 0
	for vi := range viChan {
		if rv.isValidData(*vi) {
			result = rv.addVoteToResult(*vi, result)
			nVoted++
		}
	}
	return &result, nVoted, nil
}

