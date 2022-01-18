package crypto

type ISignKeyPair interface{
	Sign() ISignKey
	Verify() IVerfKey
}
type ISignKey interface{
	Sign(m []byte) []byte
	Verify() IVerfKey
	Equals(sk2 ISignKey) bool
	Marshal() []byte
}
type IVerfKey interface{
	Verify(data, sign []byte) bool
	Equals(vk2 IVerfKey) bool
	Marshal() []byte
}

type IEncryptKeyPair interface{
	Public() IPubKey
	Private() IPriKey
}
type IPriKey interface{
	Decrypt(m []byte) ([]byte, error)
	Public() IPubKey
	Equals(pri2 IPriKey) bool
	Marshal() []byte
}
type IPubKey interface{
	Encrypt(data []byte) ([]byte, error)
	Equals(pub2 IPubKey) bool
	Marshal() []byte
}