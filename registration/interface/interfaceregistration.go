package interfaceregistration


type IRegistration interface {
	Close()
	Registrate(userData ...string) (string, error)
}
