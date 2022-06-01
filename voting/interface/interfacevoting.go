package interfacevoting

import (
	"fyne.io/fyne/v2"

	vutil "github.com/pilinsin/easy-voting/voting/util"
)

type iBaseVoting interface {
	Close()
	Config() *vutil.Config
	SetIdentity(string)
	GetIdentity() string
}
type ITypedVoting interface {
	iBaseVoting
	NewVotingForm(ngs []string) IVotingForm
	Type() string
	Vote(data vutil.VoteInt) error
	GetMyVote() (*vutil.VoteInt, error)
	GetResult() (*vutil.VoteResult, int, int, error)
}
type IVoting interface {
	ITypedVoting
	Address() string
}

type IVotingForm interface {
	fyne.CanvasObject
	VoteInt() vutil.VoteInt
}
