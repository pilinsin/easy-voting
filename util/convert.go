package util

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"io"
	"strconv"
	"strings"
)

func BytesToReader(b []byte) io.Reader {
	return bytes.NewBuffer(b)
}

func ReaderToBytes(reader io.Reader) []byte {
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return nil
	}

	return buf.Bytes()
}

func Bytes64ToStr(b []byte) string {
	b, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}
func StrToBytes64(str string) []byte {
	b64Str := base64.StdEncoding.EncodeToString([]byte(str))
	return []byte(b64Str)
}
func Bytes64ToStrs(b []byte) []string {
	s := Bytes64ToStr(b)
	return strings.Split(s, " |_| ")
}
func StrsToBytes64(strs []string) []byte {
	s := strings.Join(strs, " |_| ")
	return StrToBytes64(s)
}
func Bytes64ToInt(b []byte) int {
	s := Bytes64ToStr(b)
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	} else {
		return i
	}
}
func IntToBytes64(i int) []byte {
	s := strconv.Itoa(i)
	return StrToBytes64(s)
}

func Marshal(objWithPublicMembers interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(objWithPublicMembers)
	return buf.Bytes(), err
}
func Unmarshal(b []byte, objWithPublicMembers interface{}) error {
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	err := dec.Decode(objWithPublicMembers)
	return err
}
