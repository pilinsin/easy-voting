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

type blockVoting struct {
	voting
	total int
}

func NewBlockVoting(vCfgCid string, identity *rutil.UserIdentity, is *ipfs.IPFS) viface.IVoting {
	vCfg, _ := vutil.ConfigFromCid(vCfgCid, is)
	bv := &blockVoting{
		total: vCfg.VParam().Total,
	}
	bv.init(vCfgCid, identity, is)
	return bv
}

type blockForm struct {
	widget.Form
	vi         vutil.VoteInt
	numChecked int
	total      int
}

func newBlockForm(options []string, total int, vi vutil.VoteInt) *blockForm {
	bf := &blockForm{
		Form:       widget.Form{},
		vi:         vi,
		numChecked: 0,
		total:      total,
	}

	for _, opt := range options {
		chk := &widget.Check{
			DisableableWidget: widget.DisableableWidget{},
			Text:              opt,
		}
		chk.OnChanged = func(isChecked bool) {
			if isChecked {
				if bf.numChecked < bf.total {
					bf.vi[chk.Text] = 1
					bf.numChecked++
				} else {
					chk.Checked = false
				}
			} else {
				bf.vi[chk.Text] = 0
				bf.numChecked--
			}
			log.Println(bf.vi)
		}
		chk.ExtendBaseWidget(chk)
		bf.Items = append(bf.Items, widget.NewFormItem("", chk))
	}

	bf.ExtendBaseWidget(bf)
	return bf
}
func (bf *blockForm) VoteInt() vutil.VoteInt {
	return bf.vi
}
func (bv *blockVoting) NewVotingForm(ngs []string) viface.IVotingForm {
	return newBlockForm(ngs, bv.total, bv.newDefaultVoteInt())
}

func (bv *blockVoting) newDefaultVoteInt() vutil.VoteInt {
	vi := make(vutil.VoteInt)
	for _, ng := range bv.candNameGroups() {
		vi[ng] = 0
	}

	return vi
}

func (bv *blockVoting) isValidData(vi vutil.VoteInt) bool {
	if !bv.isCandsMatch(vi) {
		return false
	}

	numTrue := 0
	for _, vote := range vi {
		if vote > 0 {
			numTrue++
		}
	}
	return numTrue > 0 && numTrue < bv.total
}

func (bv *blockVoting) Type() string {
	return "blockvoting"
}

func (bv *blockVoting) Vote(data vutil.VoteInt) error {
	if bv.isValidData(data) {
		return bv.baseVote(data)
	} else {
		return util.NewError("invalid vote")
	}
}

func (bv blockVoting) GetMyVote() (string, error) {
	vi, err := bv.baseGetMyVote()
	if err != nil {
		return "", err
	} else if vi != nil && bv.isValidData(*vi) {
		return ipfs.ToCidWithAdd(vi.Marshal(), bv.is), nil
	} else {
		return "", util.NewError("invalid vote")
	}
}

func (bv blockVoting) newResult() map[string]map[string]int {
	result := make(map[string]map[string]int, len(bv.cands))
	for _, name := range bv.candNameGroups() {
		result[name] = map[string]int{"n_votes": 0}
	}
	return result
}
func (bv blockVoting) addVote2Result(vi vutil.VoteInt, result map[string]map[string]int) map[string]map[string]int {
	for k, v := range vi {
		if v > 0 {
			result[k]["n_votes"]++
		}
	}
	return result
}
func (bv blockVoting) CountMyResult() (string, error){
	return bv.countResult(bv.baseGetMyVotes)
}
func (bv blockVoting) CountManResult() (string, error){
	return bv.countResult(bv.baseGetVotes)
}
func (bv blockVoting) countResult(gvf viface.GetVotesFunc) (string, error) {
	result := bv.newResult()

	viChan, nVoters, err := gvf()
	if err != nil {
		return "", err
	}

	nVoted := 0
	for vi := range viChan {
		if bv.isValidData(*vi) {
			result = bv.addVote2Result(*vi, result)
			nVoted++
		}
	}
	m := vutil.NewResult(result, nVoted, nVoters).Marshal()
	return ipfs.File.Add(m, bv.is), nil
}

