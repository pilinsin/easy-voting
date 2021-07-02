package util

import (
	"context"
	"log"
	"io"
	"strings"
	"golang.org/x/crypto/sha3"
	"encoding/base64"
	"crypto/rand"
	"math"
)

func CheckError(err error){
	if err != nil{
		log.Panic(err)
	}
}


func BoolPtr(b bool) *bool{
	return &b
}

func StrPtr(s string) *string{
	return &s
}

func Str2Reader(str string) io.Reader{
	return strings.NewReader(str)
}

func Reader2Str(reader io.Reader) string{
	buf := new(strings.Builder)
	io.Copy(buf, reader)
	return buf.String()
}

func NewContext() context.Context{
	ctx, _ := context.WithCancel(context.Background())
	return ctx
}


func Hash(str string) string{
	bStr := []byte(str)
	var base int
	if len(str) <= 64{
		base = 64
	}else{
		base = len(str)
	}

	hashed := make([]byte, base)
	sha3.ShakeSum256(hashed, bStr)
	return base64.StdEncoding.EncodeToString(hashed)
} 

func GenRandomBytes(length int) []byte{
	rb := make([]byte, length)
	_, err := rand.Read(rb)
	CheckError(err)

	return rb
}

func GenUniqueID(length int, step int) string{
	idChars := "0123456789abcdefghijkmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ#%&"
	idBytes := []byte(idChars)

	nSteps := int(math.Ceil(float64(length)/float64(step)))
	uidSize := length + nSteps - 1
	uid := GenRandomBytes(uidSize)

	for i, st := 0, 0; i<uidSize; i++{
		if st == step{
			uid[i] = []byte("-")[0]
			st = 0
		}else{
			//1byte = 8bit, 8bit >>2 = 6bit ([0, 63])
			uid[i] = idBytes[int(uid[i]) >> 2]
			st++
		}
	}

	return string(uid)
}

