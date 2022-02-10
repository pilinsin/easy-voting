package interfacevoting

import (
	"fyne.io/fyne/v2"

	"github.com/pilinsin/util/crypto"
	vutil "github.com/pilinsin/easy-voting/voting/util"
)

type iBaseVoting interface {
	Close()
	VerifyIdentity() bool
	VerifyHashVerfMap() bool
	VerifyMyResult() (bool, error)
	VerifyManResult() (bool, error)
}
type IVoting interface {
	iBaseVoting
	NewVotingForm(ngs []string) IVotingForm
	Type() string
	Vote(data vutil.VoteInt) error
	CountMyResult(manPriKey crypto.IPriKey) (string, error)
	CountManResult(manPriKey crypto.IPriKey) (string, error)
}

type IManager interface {
	Load() error
	Close()
	IsValidUser(userData ...string) bool
	Registrate() error
	Log()
	UploadResultBox() error
	VerifyResultBox() (bool, error)
}

type IVotingForm interface {
	fyne.CanvasObject
	VoteInt() vutil.VoteInt
}
