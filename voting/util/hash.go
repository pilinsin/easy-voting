package votingutil

import (
	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)

type BoxVidHash string
func NewBoxVidHash(mRBox []byte, votingID string) NameVidHash {
	hash := crypto.HashWithSize(mRBox, []byte(votingID), 384)
	hash := crypto.HashWithSize([]byte(votingID), hash, 384)
	return NameVidHash(util.AnyBytes64ToStr(hash))
}

type UidVidHash string
func NewUidVidHash(userID, votingID string) UidVidHash {
	hash := crypto.Hash([]byte(userID), []byte(votingID))
	hash = crypto.HashWithSize([]byte(votingID), hash, 64)
	return UidVidHash(util.AnyBytes64ToStr(hash))
}

