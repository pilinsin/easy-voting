package util

import (
	"bytes"
	"strconv"
	"os"
	"path/filepath"
)

func ExeDirPath() string{
	exe, _ := os.Executable()
	return filepath.Dir(exe)
}
func PathJoin(base string, adders ...string) string{
	for _, adder := range adders{
		base = filepath.Join(base, adder)
	}
	return base
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

//m1 > m2
func MapContainMap(m1, m2 map[string][]byte) bool {
	for k2, v2 := range m2 {
		if v1, ok := m1[k2]; !ok {
			return false
		} else if v2 != nil && !bytes.Equal(v1, v2) {
			return false
		}
	}
	return true
}

func Arange(start, stop, step int) []int {
	if start < 0 {
		start = 0
	}
	if stop < 0 {
		stop = 1
	}
	if step < 0 {
		step = 1
	}

	var arr []int
	for i := start; i < stop; i += step {
		arr = append(arr, i)
	}
	return arr
}

func ArangeStr(start, stop, step int) []string {
	if start < 0 {
		start = 0
	}
	if stop < 0 {
		stop = 1
	}
	if step < 0 {
		step = 1
	}

	var arr []string
	for i := start; i < stop; i += step {
		arr = append(arr, strconv.Itoa(i))
	}
	return arr
}
