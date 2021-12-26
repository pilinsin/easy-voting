package votingutil

import (
	"time"

	"EasyVoting/ipfs"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
	"EasyVoting/util/ecies"
	"EasyVoting/util/ed25519"
)

type VotingBox struct {
	manEncVote  []byte
	userEncVote []byte
}

func NewVotingBox() *VotingBox {
	return &VotingBox{}
}
func (vb *VotingBox) Vote(vi VoteInt, manPubKey *ecies.PubKey, identity *rutil.UserIdentity) {
	sv := newSignedVote(vi, identity).marshal()
	mev, err := manPubKey.Encrypt(sv)
	if err != nil {
		return
	}
	uev, err := identity.Private().Public().Encrypt(sv)
	if err != nil {
		return
	}

	vb.manEncVote = mev
	vb.userEncVote = uev
}
func (vb VotingBox) GetVote(manPriKey *ecies.PriKey) ([]byte, error) {
	mSignedVote, err := manPriKey.Decrypt(vb.manEncVote)
	if err != nil {
		return nil, err
	}

	if _, err := UnmarshalSignedVote(mSignedVote); err != nil {
		return nil, err
	}
	return mSignedVote, nil
}
func (vb VotingBox) GetMyVote(identity *rutil.UserIdentity) ([]byte, error) {
	mSignedVote, err := identity.Private().Decrypt(vb.userEncVote)
	if err != nil {
		return nil, err
	}

	if _, err := UnmarshalSignedVote(mSignedVote); err != nil {
		return nil, err
	}
	return mSignedVote, nil
}
func (vb *VotingBox) FromName(vdName string, is *ipfs.IPFS) error {
	m, err := ipfs.FromName(vdName, is)
	if err != nil {
		return err
	}
	err = vb.Unmarshal(m)
	if err != nil {
		return err
	}
	return nil
}
func (vb VotingBox) Marshal() []byte {
	mvb := &struct {
		Mev []byte
		Uev []byte
	}{vb.manEncVote, vb.userEncVote}
	m, _ := util.Marshal(mvb)
	return m
}
func (vb *VotingBox) Unmarshal(m []byte) error {
	mvb := &struct {
		Mev []byte
		Uev []byte
	}{}
	err := util.Unmarshal(m, mvb)
	if err != nil {
		return err
	}

	vb.manEncVote = mvb.Mev
	vb.userEncVote = mvb.Uev
	return nil
}

type signedVote struct {
	*vote
	sign []byte
}

func newSignedVote(vi VoteInt, identity *rutil.UserIdentity) *signedVote {
	vote := newVote(vi)
	sign := identity.Sign().Sign(vote.marshal())
	return &signedVote{
		vote: vote,
		sign: sign,
	}
}
func (sv signedVote) Verify(verfKey *ed25519.VerfKey) bool {
	return verfKey.Verify(sv.vote.marshal(), sv.sign)
}
func (sv signedVote) marshal() []byte {
	msv := &struct{ V, S []byte }{sv.vote.marshal(), sv.sign}
	m, _ := util.Marshal(msv)
	return m
}
func UnmarshalSignedVote(m []byte) (*signedVote, error) {
	msv := &struct{ V, S []byte }{}
	if err := util.Unmarshal(m, msv); err != nil {
		return nil, err
	}

	V, err := unmarshalVote(msv.V)
	if err != nil {
		return nil, err
	}
	return &signedVote{V, msv.S}, nil
}

type vote struct {
	vote VoteInt
	t    time.Time
}

//todo: timezone
func newVote(vi VoteInt) *vote {
	return &vote{
		vote: vi,
		t:    time.Now(),
	}
}
func (vt vote) Vote(tInfo *util.TimeInfo) (*VoteInt, error) {
	if ok := tInfo.WithinTime(vt.t); ok {
		return &vt.vote, nil
	} else {
		return nil, util.NewError("invalid time error")
	}
}
func (vt vote) marshal() []byte {
	mvt := &struct {
		V VoteInt
		T time.Time
	}{vt.vote, vt.t}
	m, _ := util.Marshal(mvt)
	return m
}
func unmarshalVote(m []byte) (*vote, error) {
	mvt := &struct {
		V VoteInt
		T time.Time
	}{}
	err := util.Unmarshal(m, mvt)
	if err != nil {
		return nil, err
	}
	return &vote{mvt.V, mvt.T}, nil
}

type VoteInt map[string]int