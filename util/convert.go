package util

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"io"
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

func AnyBytes64ToStr(b []byte) string {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.Encode(&struct{ B []byte }{b})
	s := buf.String()
	s = strings.Split(s, ":")[1]
	return strings.Split(s, "\"")[1]
}
func StrToAnyBytes64(str string) []byte {
	str = "{\"B\":\"" + str + "\"}"
	buf := bytes.NewBufferString(str)
	dec := json.NewDecoder(buf)
	obj := &struct{ B []byte }{}
	dec.Decode(obj)
	return obj.B
}
func Bytes64ToAnyStr(b []byte) string {
	b, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}
func AnyStrToBytes64(str string) []byte {
	b64Str := base64.StdEncoding.EncodeToString([]byte(str))
	return []byte(b64Str)
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
