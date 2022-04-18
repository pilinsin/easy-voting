package interfacevoting

import (
	"fyne.io/fyne/v2"

	vutil "github.com/pilinsin/easy-voting/voting/util"
)
type iBaseVoting interface {
	Close()
}
type IVoting interface {
	iBaseVoting
	NewVotingForm(ngs []string) IVotingForm
	Type() string
	Vote(data vutil.VoteInt) error
	GetMyVote() (*vutil.VoteInt, error)
	GetResult() (*vutil.VoteResult, int, error)
}


type IVotingForm interface {
	fyne.CanvasObject
	VoteInt() vutil.VoteInt
}
