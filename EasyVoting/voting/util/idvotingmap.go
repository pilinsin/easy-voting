package votingutil

import (
	"time"

	"EasyVoting/ipfs"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
	"EasyVoting/util/ecies"
	"EasyVoting/util/ed25519"
)

type votingData struct {
	encKeyFile  []byte
	name        string
	userVerfKey *ed25519.VerfKey
}

func NewVotingData(rIpnsName string, is *ipfs.IPFS) *votingData {
	kf := ipfs.NewKeyFile()
	name, _ := kf.Name()
	rb, _ := rutil.RBoxFromName(rIpnsName, is)
	encKf, _ := rb.Public().Encrypt(kf.Marshal())
	return &votingData{
		encKeyFile:  encKf,
		name:        name,
		userVerfKey: rb.Verify(),
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
	return identity.Sign().Verify().Equals(vd.userVerfKey)
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
		VerfKey    []byte
	}{vd.encKeyFile, vd.name, vd.userVerfKey.Marshal()}
	m, _ := util.Marshal(mvd)
	return m
}
func (vd *votingData) Unmarshal(m []byte) error {
	mvd := &struct {
		EncKeyFile []byte
		Name       string
		VerfKey    []byte
	}{}
	err := util.Unmarshal(m, mvd)
	if err != nil {
		return err
	}

	vd.encKeyFile = mvd.EncKeyFile
	vd.name = mvd.Name
	if err := vd.userVerfKey.Unmarshal(mvd.VerfKey); err != nil {
		return err
	}
	return nil
}

type keyValue struct {
	uvHash  UidVidHash
	vb      *VotingBox
	verfKey *ed25519.VerfKey
}

func (kv keyValue) Key() UidVidHash {
	return kv.uvHash
}
func (kv keyValue) Value() (*VotingBox, *ed25519.VerfKey) {
	return kv.vb, kv.verfKey
}

type idVotingMap struct {
	rm    *ipfs.ReccurentMap
	tInfo *util.TimeInfo
}

func NewIdVotingMap(capacity int, tInfo *util.TimeInfo) *idVotingMap {
	return &idVotingMap{
		rm:    ipfs.NewReccurentMap(capacity),
		tInfo: tInfo,
	}
}
func (ivm idVotingMap) Len(is *ipfs.IPFS) int {
	return ivm.rm.Len(is)
}
func (ivm idVotingMap) withinTime() bool {
	return ivm.tInfo.WithinTime(time.Now())
}

//votingData->Name->VotingBox
func (ivm idVotingMap) Next(is *ipfs.IPFS) <-chan *VotingBox {
	ch := make(chan *VotingBox)
	go func() {
		defer close(ch)
		for m := range ivm.rm.Next(is) {
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

//{uvHash, VotingBox, VerfKey}
func (ivm idVotingMap) NextKeyValue(is *ipfs.IPFS) <-chan *keyValue {
	ch := make(chan *keyValue)
	go func() {
		defer close(ch)
		for kv := range ivm.rm.NextKeyValue(is) {
			vd := &votingData{}
			err := vd.Unmarshal(kv.Value())
			if err != nil {
				continue
			}
			vb := &VotingBox{}
			err = vb.FromName(vd.name, is)
			if err == nil {
				ch <- &keyValue{UidVidHash(kv.Key()), vb, vd.userVerfKey}
			}
		}
	}()
	return ch
}
func (ivm *idVotingMap) Vote(hash UidVidHash, vote VoteInt, id *rutil.UserIdentity, manPubKey *ecies.PubKey, is *ipfs.IPFS) {
	if m, ok := ivm.rm.ContainKey(string(hash), is); !ok {
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
		data := NewVotingData(rIpnsName, is)
		ivm.rm.Append(string(hash), data.Marshal(), is)
	}
}

//*votingBox, *verfKey, bool
func (ivm idVotingMap) ContainHash(hash UidVidHash, is *ipfs.IPFS) (*VotingBox, *ed25519.VerfKey, bool) {
	if m, ok := ivm.rm.ContainKey(hash, is); !ok {
		return nil, nil, false
	} else {
		vd := &votingData{}
		err := vd.Unmarshal(m)
		if err != nil {
			return nil, nil, false
		}
		vb := &VotingBox{}
		err = vb.FromName(vd.name, is)
		if err != nil {
			return nil, nil, false
		}
		return vb, vd.userVerfKey, true
	}
}
func (ivm idVotingMap) Marshal() []byte {
	mivm := &struct {
		Mrm      []byte
		TimeInfo *util.TimeInfo
	}{ivm.rm.Marshal(), ivm.tInfo}
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

	ivm := &idVotingMap{}
	err = ivm.rm.Unmarshal(mivm.Mrm)
	ivm.tInfo = mivm.TimeInfo
	if err != nil {
		return nil, err
	}
	return ivm, nil
}
func IdVotingMapFromCid(ivmCid string, is *ipfs.IPFS) (*idVotingMap, error) {
	m, err := ipfs.FromCid(ivmCid, is)
	if err != nil {
		return nil, err
	}
	return UnmarshalIdVotingMap(m)
}
