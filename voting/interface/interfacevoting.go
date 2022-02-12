package interfacevoting

import (
	"fyne.io/fyne/v2"

	vutil "github.com/pilinsin/easy-voting/voting/util"
)
type GetVotesFunc func() (<-chan *vutil.VoteInt, int, error)
type iBaseVoting interface {
	Load() error
	Close()
	Update()
	Log()
	VerifyIdentity() bool
	VerifyHashVerfMap() bool
}
type IVoting interface {
	iBaseVoting
	NewVotingForm(ngs []string) IVotingForm
	Type() string
	Vote(data vutil.VoteInt) error
	CountMyResult() (string, error)
	CountManResult() (string, error)
}

type IManager interface {
	Load() error
	Close()
	IsValidUser(userData ...string) bool
	Registrate() error
	Log()
	UploadResultBox() error
}

type IVotingForm interface {
	fyne.CanvasObject
	VoteInt() vutil.VoteInt
}
