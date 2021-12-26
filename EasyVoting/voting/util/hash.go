package votingutil

import (
	"EasyVoting/util"
)

type UidVidHash string

func NewUidVidHash(userID, votingID string) UidVidHash {
	return UidVidHash(util.Bytes64ToStr(util.Hash([]byte(userID), []byte(votingID))))
}

type NameVidHash string

func NewNameVidHash(ipnsName, votingID string) NameVidHash {
	return NameVidHash(util.Bytes64ToStr(util.Hash([]byte(ipnsName), []byte(votingID))))
}
