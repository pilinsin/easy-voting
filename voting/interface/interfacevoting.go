package interfacevoting

import (
	"fyne.io/fyne/v2"

	vutil "EasyVoting/voting/util"
)

type iBaseVoting interface {
	Close()
	VerifyIdentity() bool
	VerifyIdVerfKeyMap() bool
	VerifyResultMap() (bool, error)
}
type IVoting interface {
	iBaseVoting
	NewVotingForm(ngs []string) IVotingForm
	Type() string
	Vote(data vutil.VoteInt) error
	GetMyVote() (vutil.VoteInt, error)
	Count() (map[string]map[string]int, int, int, error)
}

type IManager interface {
	Close()
	IsValidUser(userData ...string) bool
	Registrate() error
	GetResultMap() error
	VerifyResultMap() (bool, error)
}

type IVotingForm interface {
	fyne.CanvasObject
	VoteInt() vutil.VoteInt
}
