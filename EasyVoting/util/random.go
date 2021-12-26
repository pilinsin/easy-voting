package util

import (
	"crypto/rand"
	"math"
	"math/big"
)

//only for 64bit
func RandInt(max int) int {
	if max < 0 {
		max = 2147483647
	}
	bInt, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(bInt.Int64())
}

func GenRandomBytes(length int) []byte {
	rb := make([]byte, length)
	for {
		_, err := rand.Read(rb)
		if err == nil {
			return rb
		}
	}
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
