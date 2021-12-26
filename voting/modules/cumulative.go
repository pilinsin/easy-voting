package votingmodule

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
	viface "EasyVoting/voting/interface"
	vutil "EasyVoting/voting/util"
)

type cumulativeVoting struct {
	voting
	min   int
	total int
}

func NewCumulativeVoting(vCfgCid string, identity *rutil.UserIdentity, is *ipfs.IPFS) viface.IVoting {
	vCfg, _ := vutil.ConfigFromCid(vCfgCid, is)
	cv := &cumulativeVoting{
		min:   vCfg.VParam().Min,
		total: vCfg.VParam().Total,
	}
	cv.init(vCfgCid, identity, is)
	return cv
}

type cumulativeForm struct {
	widget.Form
	vi        vutil.VoteInt
	totalVals int
}
type cSliderWithText struct {
	widget.Slider
	PreVal int
	Text   string
}

func newCumulativeForm(options []string, min, total int, vi vutil.VoteInt) *cumulativeForm {
	totalLabel := widget.NewLabel(strconv.Itoa(total) + "/" + strconv.Itoa(total))
	totalItems := []*widget.FormItem{widget.NewFormItem("", totalLabel)}
	cf := &cumulativeForm{
		Form:      widget.Form{Items: totalItems},
		vi:        vi,
		totalVals: 0,
	}

	for _, opt := range options {
		valLabel := widget.NewLabel("0")
		sl := &cSliderWithText{
			Slider: widget.Slider{
				Value:       0,
				Min:         float64(min),
				Max:         float64(total),
				Step:        1,
				Orientation: widget.Horizontal,
			},
			PreVal: 0,
			Text:   opt,
		}
		sl.OnChanged = func(val float64) {
			if cf.totalVals+int(val)-sl.PreVal > total {
				val = float64(sl.PreVal)
				sl.SetValue(val)
			} else {
				cf.totalVals += int(val) - sl.PreVal
				sl.PreVal = int(val)
			}
			cf.vi[sl.Text] = int(val)
			totalLabel.Text = strconv.Itoa(total-cf.totalVals) + "/" + strconv.Itoa(total)
			valLabel.Text = strconv.Itoa(int(val))
			log.Println(cf.vi)
		}
		sl.ExtendBaseWidget(sl)

		sl2 := container.NewGridWrap(fyne.NewSize(250.0, 20.0), sl)
		slTxt := widget.NewLabel(sl.Text)
		cf.Items = append(cf.Items, widget.NewFormItem("", container.NewHBox(sl2, valLabel, slTxt)))
	}

	cf.ExtendBaseWidget(cf)
	return cf
}
func (cf *cumulativeForm) VoteInt() vutil.VoteInt {
	return cf.vi
}
func (cv *cumulativeVoting) NewVotingForm(ngs []string) viface.IVotingForm {
	return newCumulativeForm(ngs, cv.min, cv.total, cv.newDefaultVoteInt())
}
func (cv *cumulativeVoting) newDefaultVoteInt() vutil.VoteInt {
	vi := make(vutil.VoteInt)
	for _, id := range cv.candNameGroups() {
		vi[id] = cv.min
	}
	return vi
}

func (cv *cumulativeVoting) isValidData(vi vutil.VoteInt) bool {
	if !cv.isCandsMatch(vi) {
		return false
	}

	tl := 0
	for _, vote := range vi {
		if vote < cv.min {
			return false
		}
		tl += vote
	}
	return tl <= cv.total
}

func (cv *cumulativeVoting) Type() string {
	return "cumulativevoting"
}

func (cv *cumulativeVoting) Vote(data vutil.VoteInt) error {
	if cv.isValidData(data) {
		return cv.baseVote(data)
	} else {
		return util.NewError("invalid vote")
	}
}
func (cv cumulativeVoting) GetMyVote() (vutil.VoteInt, error) {
	vi, err := cv.baseGetMyVote()
	if err != nil {
		return cv.newDefaultVoteInt(), err
	} else if vi != nil && cv.isValidData(*vi) {
		return *vi, nil
	} else {
		return cv.newDefaultVoteInt(), util.NewError("invalid vote")
	}
}

func (cv cumulativeVoting) newResult() map[string]map[string]int {
	result := make(map[string]map[string]int, len(cv.cands))
	for _, name := range cv.candNameGroups() {
		result[name] = map[string]int{"n_votes": 0}
	}
	return result
}
func (cv cumulativeVoting) addVote2Result(vi vutil.VoteInt, result map[string]map[string]int) map[string]map[string]int {
	for k, v := range vi {
		result[k]["n_votes"] += v
	}
	return result
}
func (cv cumulativeVoting) Count() (map[string]map[string]int, int, int, error) {
	result := cv.newResult()

	viChan, nVoters, err := cv.baseGetVotes()
	if err != nil {
		return make(map[string]map[string]int), -1, -1, err
	}

	nVoted := 0
	for vi := range viChan {
		if cv.isValidData(*vi) {
			result = cv.addVote2Result(*vi, result)
			nVoted++
		}
	}
	return result, nVoted, nVoters, nil
}