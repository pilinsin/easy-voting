package util

import (
	"encoding/base64"
	"os"
	"path/filepath"

	"golang.org/x/crypto/argon2"
)

func exeDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}

func BaseDir(mode, addr string) string {
	hashAddr := hashDir(addr, mode)
	return filepath.Join(exeDir(), "stores", hashAddr)
}

func hashDir(addr, salt string) string {
	txt := "baseDir_hash(cid, salt+padd)"
	b := argon2.IDKey([]byte(addr), []byte(salt+txt), 1, 64*1024, 4, 16)
	return base64.URLEncoding.EncodeToString(b)
}
