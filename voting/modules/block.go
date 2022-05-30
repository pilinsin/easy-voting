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

type blockVoting struct {
	voting
	total int
}

func NewBlockVoting(ctx context.Context, vCfg *vutil.Config, storeDir string, bs []peer.AddrInfo, save bool) (viface.ITypedVoting, error) {
	bv := &blockVoting{
		total: vCfg.Params.Total,
	}
	if err := bv.init(ctx, vCfg, storeDir, bs, save); err != nil {
		return nil, err
	}
	return bv, nil
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
		return errors.New("invalid vote")
	}
}

func (bv blockVoting) GetMyVote() (*vutil.VoteInt, error) {
	vi, err := bv.baseGetMyVote()
	if err != nil {
		return nil, err
	} else if vi != nil && bv.isValidData(*vi) {
		return vi, nil
	} else {
		return nil, errors.New("invalid vote")
	}
}

func (bv blockVoting) newResult() vutil.VoteResult {
	result := make(vutil.VoteResult, len(bv.cands))
	for _, name := range bv.candNameGroups() {
		result[name] = map[string]int{"n_votes": 0}
	}
	return result
}
func (bv blockVoting) addVoteToResult(vi vutil.VoteInt, result vutil.VoteResult) vutil.VoteResult {
	for k, v := range vi {
		if v > 0 {
			result[k]["n_votes"]++
		}
	}
	return result
}
func (bv blockVoting) GetResult() (*vutil.VoteResult, int, int, error) {
	result := bv.newResult()

	viChan, nVoters, err := bv.baseGetVotes()
	if err != nil {
		return nil, -1, -1, err
	}

	nVoted := 0
	for vi := range viChan {
		if bv.isValidData(*vi) {
			result = bv.addVoteToResult(*vi, result)
			nVoted++
		}
	}
	return &result, nVoters, nVoted, nil
}
