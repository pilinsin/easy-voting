package votingmodule

import (
	"errors"
	"log"
	"sort"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	viface "github.com/pilinsin/easy-voting/voting/interface"
	vutil "github.com/pilinsin/easy-voting/voting/util"
	crdt "github.com/pilinsin/p2p-verse/crdt"
	ipfs "github.com/pilinsin/p2p-verse/ipfs"
)

type preferenceVoting struct {
	voting
}

func NewPreferenceVoting(vCfg *vutil.Config, is ipfs.Ipfs, hkm crdt.ISignatureStore, ivm crdt.IUpdatableSignatureStore) (viface.ITypedVoting, error) {
	pv := &preferenceVoting{}
	if err := pv.init(vCfg, is, hkm, ivm); err != nil {
		return nil, err
	}
	return pv, nil
}

type preferenceForm struct {
	widget.Form
	vi vutil.VoteInt
}
type pSliderWithText struct {
	widget.Slider
	Text string
}

func newPreferenceForm(options []string, vi vutil.VoteInt) *preferenceForm {
	pf := &preferenceForm{
		Form: widget.Form{},
		vi:   vi,
	}

	for _, opt := range options {
		valLabel := widget.NewLabel("1")
		sl := &pSliderWithText{
			Slider: widget.Slider{
				Value:       1,
				Min:         1,
				Max:         float64(len(options)),
				Step:        1,
				Orientation: widget.Horizontal,
			},
			Text: opt,
		}
		sl.OnChanged = func(val float64) {
			pf.vi[sl.Text] = int(val)
			valLabel.Text = strconv.Itoa(int(sl.Value))
			log.Println(pf.vi)
		}
		sl.ExtendBaseWidget(sl)

		sl2 := container.NewGridWrap(fyne.NewSize(250.0, 20.0), sl)
		slTxt := widget.NewLabel(sl.Text)
		pf.Items = append(pf.Items, widget.NewFormItem("", container.NewHBox(sl2, valLabel, slTxt)))
	}

	pf.ExtendBaseWidget(pf)
	return pf
}
func (pf *preferenceForm) VoteInt() vutil.VoteInt {
	return pf.vi
}
func (pv *preferenceVoting) NewVotingForm(ngs []string) viface.IVotingForm {
	return newPreferenceForm(ngs, pv.newDefaultVoteInt())
}

func (pv *preferenceVoting) newDefaultVoteInt() vutil.VoteInt {
	vi := make(vutil.VoteInt)
	for _, ng := range pv.candNameGroups() {
		vi[ng] = 0
	}
	return vi
}

func (pv *preferenceVoting) isValidData(vi vutil.VoteInt) bool {
	if !pv.isCandsMatch(vi) {
		return false
	}

	var ps []int
	for _, vote := range vi {
		ps = append(ps, vote)
	}
	sort.Ints(ps)
	for i, v := range ps {
		if i != v {
			return false
		}
	}
	return true
}

func (pv *preferenceVoting) Type() string {
	return "preferencevoting"
}

func (pv *preferenceVoting) Vote(data vutil.VoteInt) error {
	if pv.isValidData(data) {
		return pv.baseVote(data)
	} else {
		return errors.New("invalid vote")
	}
}
func (pv preferenceVoting) GetMyVote() (*vutil.VoteInt, error) {
	vi, err := pv.baseGetMyVote()
	if err != nil {
		return nil, err
	} else if vi != nil && pv.isValidData(*vi) {
		return vi, nil
	} else {
		return nil, errors.New("invalid vote")
	}
}

func (pv preferenceVoting) newResult() vutil.VoteResult {
	result := make(vutil.VoteResult, len(pv.cands))
	for _, name := range pv.candNameGroups() {
		result[name] = make(map[string]int, len(pv.cands))
		for idx := 0; idx < len(pv.cands); idx++ {
			result[name][strconv.Itoa(idx)] = 0
		}
	}
	return result
}
func (pv preferenceVoting) addVoteToResult(vi vutil.VoteInt, result vutil.VoteResult) vutil.VoteResult {
	for k, v := range vi {
		result[k][strconv.Itoa(v)]++
	}
	return result
}

func (pv preferenceVoting) GetResult() (*vutil.VoteResult, int, int, error) {
	result := pv.newResult()

	viChan, nVoters, err := pv.baseGetVotes()
	if err != nil {
		return nil, -1, -1, err
	}

	nVoted := 0
	for vi := range viChan {
		if pv.isValidData(*vi) {
			result = pv.addVoteToResult(*vi, result)
			nVoted++
		}
	}
	return &result, nVoters, nVoted, nil
}
