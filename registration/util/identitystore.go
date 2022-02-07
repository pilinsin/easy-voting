package registrationutil

import(
	"os"

	"github.com/pilinsin/util"
)

func identityStorePath() string{
	dir := util.ExeDirPath()
	dir = util.PathJoin(dir, "identitystore")

	os.Mkdir(dir, 0755)
	return util.PathJoin(dir, "identities.dat")
}

type identityStore map[StoreHash][]byte
func NewIdentityStore() *identityStore{
	isPath := identityStorePath()
	if m, err := os.ReadFile(isPath); err == nil{
		is := &identityStore{}
		if err := util.Unmarshal(m, is); err == nil{
			return is
		}
	}
	is := make(identityStore)
	return &is
}
func (is identityStore) Close(){
	isPath := identityStorePath()
	m, _ := util.Marshal(is)
	os.WriteFile(isPath, m, 0744)
}
func (is identityStore) Put(kw string, identity []byte){
	hash := NewStoreHash(kw)
	is[hash] = identity
}
func (is identityStore) Get(kw string) ([]byte, bool){
	hash := NewStoreHash(kw)
	v, ok := is[hash]
	return v, ok
}
func (is identityStore) Delete(kw string){
	hash := NewStoreHash(kw)
	delete(is, hash)
}