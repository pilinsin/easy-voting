package votingutil

import (
	"EasyVoting/util"
	"EasyVoting/util/crypto"
)

type UidVidHash string

func NewUidVidHash(userID, votingID string) UidVidHash {
	return UidVidHash(util.AnyBytes64ToStr(crypto.Hash([]byte(userID), []byte(votingID))))
}

type NameVidHash string

func NewNameVidHash(ipnsName, votingID string) NameVidHash {
	return NameVidHash(util.AnyBytes64ToStr(crypto.Hash([]byte(ipnsName), []byte(votingID))))
}
