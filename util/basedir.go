package util

import(
	"encoding/base64"
	"golang.org/x/crypto/argon2"
)

func BaseDir(mode, cfgCid string) string{
	txt := "baseDir_hash(cid, mode+padd)"
	b := argon2.IDKey([]byte(cfgCid), []byte(mode+txt), 1, 64*1024, 4, 16)
	return base64.URLEncoding.EncodeToString(b)
}