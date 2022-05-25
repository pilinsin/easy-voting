package interfaceregistration

import (
	rutil "github.com/pilinsin/easy-voting/registration/util"
)

type IRegistration interface {
	Close()
	Config() *rutil.Config
	Address() string
	Registrate(userData ...string) (string, error)
}
