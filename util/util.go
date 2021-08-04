package util

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"io"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"
)

func CheckError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
func RaiseError(a ...interface{}) {
	err := errors.New(fmt.Sprintln(a...))
	log.Panic(err)
}

func BoolPtr(b bool) *bool {
	return &b
}

func StrPtr(s string) *string {
	return &s
}

func ConstTimeBytesEqual(b1, b2 []byte) bool {
	len1 := len(b1)
	len2 := len(b2)

	var bLong []byte
	var bShort []byte
	if len1 < len2 {
		bLong = b2
		bShort = b1
	} else {
		bLong = b1
		bShort = b2
	}

	res := len1 == len2
	for idx := range bShort {
		res = bShort[idx] == bLong[idx] && res
	}

	return res
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
	return context.Background()
}
func CancelContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
func SignalContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}
func WithSignal(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
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
