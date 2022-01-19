package votingutil

import (
	"time"

	"EasyVoting/ipfs"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
	"EasyVoting/util/crypto"
)

type votingData struct {
	encKeyFile  []byte
	name        string
}

func newVotingData(rIpnsName string, is *ipfs.IPFS) *votingData {
	kf := ipfs.NewKeyFile()
	name, _ := kf.Name()
	rb, _ := rutil.RBoxFromName(rIpnsName, is)
	encKf, _ := rb.Public().Encrypt(kf.Marshal())
	return &votingData{
		encKeyFile:  encKf,
		name:        name,
	}
}
func (vd votingData) verifyIdentity(identity *rutil.UserIdentity) bool {
	if identity == nil {
		return false
	}
	priKey := identity.Private()
	if priKey == nil {
		return false
	}
	m, err := priKey.Decrypt(vd.encKeyFile)
	if err != nil {
		return false
	}
	kf := &ipfs.KeyFile{}
	if err := kf.Unmarshal(m); err != nil {
		return false
	}
	return identity.KeyFile().Equals(kf)
}
func (vd votingData) keyFile(identity *rutil.UserIdentity) *ipfs.KeyFile {
	m, _ := identity.Private().Decrypt(vd.encKeyFile)
	kf := &ipfs.KeyFile{}
	kf.Unmarshal(m)
	return kf
}
func (vd votingData) Marshal() []byte {
	mvd := &struct {
		EncKeyFile []byte
		Name       string
	}{vd.encKeyFile, vd.name}
	m, _ := util.Marshal(mvd)
	return m
}
func (vd *votingData) Unmarshal(m []byte) error {
	mvd := &struct {
		EncKeyFile []byte
		Name       string
	}{}
	err := util.Unmarshal(m, mvd)
	if err != nil {
		return err
	}

	vd.encKeyFile = mvd.EncKeyFile
	vd.name = mvd.Name
	return nil
}

type keyValue struct {
	uvHash  UidVidHash
	vb      *VotingBox
}

func (kv keyValue) Key() UidVidHash {
	return kv.uvHash
}
func (kv keyValue) Value() *VotingBox {
	return kv.vb
}

type idVotingMap struct {
	sm    *ipfs.ScalableMap
	tInfo *util.TimeInfo
}

func NewIdVotingMap(capacity int, tInfo *util.TimeInfo) *idVotingMap {
	return &idVotingMap{
		sm:    ipfs.NewScalableMap(capacity),
		tInfo: tInfo,
	}
}
func (ivm idVotingMap) Len(is *ipfs.IPFS) int {
	return ivm.sm.Len(is)
}
func (ivm idVotingMap) withinTime() bool {
	return ivm.tInfo.WithinTime(time.Now())
}

//votingData->Name->VotingBox
func (ivm idVotingMap) Next(is *ipfs.IPFS) <-chan *VotingBox {
	ch := make(chan *VotingBox)
	go func() {
		defer close(ch)
		for m := range ivm.sm.Next(is) {
			vd := &votingData{}
			err := vd.Unmarshal(m)
			if err != nil {
				continue
			}
			vb := &VotingBox{}
			err = vb.FromName(vd.name, is)
			if err == nil {
				ch <- vb
			}
		}
	}()
	return ch
}

//{uvHash, VotingBox}
func (ivm idVotingMap) NextKeyValue(is *ipfs.IPFS) <-chan *keyValue {
	ch := make(chan *keyValue)
	go func() {
		defer close(ch)
		for kv := range ivm.sm.NextKeyValue(is) {
			vd := &votingData{}
			err := vd.Unmarshal(kv.Value())
			if err != nil {
				continue
			}
			vb := &VotingBox{}
			err = vb.FromName(vd.name, is)
			if err == nil {
				ch <- &keyValue{UidVidHash(kv.Key()), vb}
			}
		}
	}()
	return ch
}
func (ivm *idVotingMap) Vote(hash UidVidHash, vote VoteInt, id *rutil.UserIdentity, manPubKey crypto.IPubKey, is *ipfs.IPFS) {
	if m, ok := ivm.sm.ContainKey(hash, is); !ok {
		return
	} else if ok := ivm.withinTime(); !ok {
		return
	} else {
		vd := &votingData{}
		err := vd.Unmarshal(m)
		if err != nil {
			return
		}
		vb := &VotingBox{}
		err = vb.FromName(vd.name, is)
		if err != nil {
			return
		}
		if ok := vd.verifyIdentity(id); !ok {
			return
		}
		if manPubKey == nil {
			return
		}

		vb.Vote(vote, manPubKey, id)
		kf := vd.keyFile(id)
		ipfs.ToNameWithKeyFile(vb.Marshal(), kf, is)
	}
}
func (ivm *idVotingMap) Append(hash UidVidHash, rIpnsName string, is *ipfs.IPFS) {
	_, err := rutil.RBoxFromName(rIpnsName, is)
	if err == nil {
		data := newVotingData(rIpnsName, is)
		ivm.sm.Append(hash, data.Marshal(), is)
	}
}

//*votingBox, bool
func (ivm idVotingMap) ContainHash(hash UidVidHash, is *ipfs.IPFS) (*VotingBox, bool) {
	if m, ok := ivm.sm.ContainKey(hash, is); !ok {
		return nil, false
	} else {
		vd := &votingData{}
		err := vd.Unmarshal(m)
		if err != nil {
			return nil, false
		}
		vb := &VotingBox{}
		err = vb.FromName(vd.name, is)
		if err != nil {
			return nil, false
		}
		return vb, true
	}
}
func (ivm idVotingMap) Marshal() []byte {
	mivm := &struct {
		Mrm      []byte
		TimeInfo *util.TimeInfo
	}{ivm.sm.Marshal(), ivm.tInfo}
	m, _ := util.Marshal(mivm)
	return m
}
func UnmarshalIdVotingMap(m []byte) (*idVotingMap, error) {
	mivm := &struct {
		Mrm      []byte
		TimeInfo *util.TimeInfo
	}{}
	err := util.Unmarshal(m, mivm)
	if err != nil {
		return nil, err
	}

	sm := &ipfs.ScalableMap{}
	if err = sm.Unmarshal(mivm.Mrm); err != nil {
		return nil, err
	}

	ivm := &idVotingMap{sm, mivm.TimeInfo}
	return ivm, nil
}
func IdVotingMapFromCid(ivmCid string, is *ipfs.IPFS) (*idVotingMap, error) {
	m, err := ipfs.FromCid(ivmCid, is)
	if err != nil {
		return nil, err
	}
	return UnmarshalIdVotingMap(m)
}
