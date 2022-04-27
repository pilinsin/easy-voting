package interfaceregistration

import(
	rutil "github.com/pilinsin/easy-voting/registration/util"
)

type IRegistration interface {
	Close()
	Config() *rutil.Config
	Registrate(userData ...string) (string, error)
}
