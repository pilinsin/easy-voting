package crypto

import(
	"golang.org/x/crypto/argon2"
)

func Hash(txt []byte, salt []byte) []byte {
	return argon2.IDKey(txt, salt, 1, 64*1024, 4, 128)
}