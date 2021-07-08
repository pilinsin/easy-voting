package util

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
	"io"
	"log"
	"math"
)

func CheckError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func BoolPtr(b bool) *bool {
	return &b
}

func StrPtr(s string) *string {
	return &s
}

func Bytes2Reader(b []byte) io.Reader {
	return bytes.NewBuffer(b)
}

func Reader2Bytes(reader io.Reader) []byte {
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(reader)
	CheckError(err)

	return buf.Bytes()
}

func Bytes64EncodeStr(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func Bytes64DecodeStr(str string) []byte {
	b, err := base64.StdEncoding.DecodeString(str)
	CheckError(err)
	return b
}

func NewContext() context.Context {
	ctx, _ := context.WithCancel(context.Background())
	return ctx
}

func Hash(txt []byte, salt []byte) []byte {
	keyLen := uint32(len(txt) + len(salt))
	if keyLen <= 64 {
		keyLen = 64
	}

	return argon2.IDKey(txt, salt, 1, 64*1024, 4, keyLen)
}

func GenRandomBytes(length int) []byte {
	rb := make([]byte, length)
	_, err := rand.Read(rb)
	CheckError(err)

	return rb
}

func GenUniqueID(length int, step int) string {
	idChars := "0123456789abcdefghijkmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ#%&"
	idBytes := []byte(idChars)

	nSteps := int(math.Ceil(float64(length) / float64(step)))
	uidSize := length + nSteps - 1
	uid := GenRandomBytes(uidSize)

	for i, st := 0, 0; i < uidSize; i++ {
		if st == step {
			uid[i] = []byte("-")[0]
			st = 0
		} else {
			//1byte = 8bit, 8bit >>2 = 6bit ([0, 63])
			uid[i] = idBytes[int(uid[i])>>2]
			st++
		}
	}

	return string(uid)
}
