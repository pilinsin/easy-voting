package votingutil

import (
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)

type PubVidHash string
func NewPubVidHash(pubKey crypto.IPubKey, votingID string) NameVidHash {
	hash := crypto.HashWithSize(pubKey.Marshal(), []byte(votingID), 384)
	hash := crypto.HashWithSize([]byte(votingID), hash, 384)
	return NameVidHash(util.AnyBytes64ToStr(hash))
}

type UidVidHash string
func NewUidVidHash(userID, votingID string) UidVidHash {
	hash := crypto.Hash([]byte(userID), []byte(votingID))
	hash = crypto.HashWithSize([]byte(votingID), hash, 64)
	return UidVidHash(util.AnyBytes64ToStr(hash))
}

type UvhHash string
func NewUvhHash(uvHash UidVidHash, votingID string) UvhHash{
	hash := crypto.Hash([]byte(uvHash), []byte(votingID))
	hash = crypto.HashWithSize([]byte(votingID), hash, 64)
	return UvhHash(util.AnyBytes64ToStr(hash))
}

