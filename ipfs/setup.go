package ipfs

import (
	"io/ioutil"
	"os"
	"path/filepath"

	config "github.com/ipfs/go-ipfs-config"
	assets "github.com/ipfs/go-ipfs/assets"
	core "github.com/ipfs/go-ipfs/core"
	plugin "github.com/ipfs/go-ipfs/plugin"
	flatfs "github.com/ipfs/go-ipfs/plugin/plugins/flatfs"
	levelds "github.com/ipfs/go-ipfs/plugin/plugins/levelds"
	repo "github.com/ipfs/go-ipfs/repo"
	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/ipfs/go-ipfs/repo/fsrepo/migrations"
	"github.com/ipfs/go-ipfs/repo/fsrepo/migrations/ipfsfetcher"
	path "github.com/ipfs/go-path"
	unixfs "github.com/ipfs/go-unixfs"
	options "github.com/ipfs/interface-go-ipfs-core/options"

	"github.com/pilinsin/easy-voting/util"
)

func newRepo(repoPath string) (repo.Repo, error) {
	cfg, err := newConfig()
	if err != nil {
		return nil, err
	}
	loadPlugins()
	if err := doInit(repoPath, cfg); err != nil {
		return nil, err
	}
	//if err := migration(repoPath); err != nil{return nil, err}

	return fsrepo.Open(repoPath)
}

func loadPlugins() {
	//ignore "already have a datastore named ~" errors
	fp := flatfs.Plugins[0].(plugin.PluginDatastore)
	fsrepo.AddDatastoreConfigHandler("flatfs", fp.DatastoreConfigParser())
	lp := levelds.Plugins[0].(plugin.PluginDatastore)
	fsrepo.AddDatastoreConfigHandler("levelds", lp.DatastoreConfigParser())
}

func newConfig() (*config.Config, error) {
	keyGenOpts := []options.KeyGenerateOption{options.Key.Type(options.Ed25519Key)}
	id, err := config.CreateIdentity(ioutil.Discard, keyGenOpts)
	if err != nil {
		return nil, err
	}
	return config.InitWithIdentity(id)
}

func doInit(repoPath string, cfg *config.Config) error {
	if err := checkWritable(repoPath); err != nil {
		return err
	}
	if fsrepo.IsInitialized(repoPath) {
		return util.NewError("ipfs config file already exists.")
	}
	if err := fsrepo.Init(repoPath, cfg); err != nil {
		return err
	}
	if err := addDefaultAssets(repoPath); err != nil {
		return err
	}

	return initializeIpnsKeyspace(repoPath)
}

func checkWritable(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(dir, 0775)
		}
		if os.IsPermission(err) {
			return util.NewError("unexpected error:", err)
		}
		return err
	}
	testFile := filepath.Join(dir, "test")
	f, err := os.Create(testFile)
	if err != nil {
		if os.IsPermission(err) {
			return util.NewError(dir, "is not writable by the user")
		}
		return util.NewError("unexpected error:", err)
	}
	f.Close()
	return os.Remove(testFile)
}
func addDefaultAssets(repoPath string) error {
	ctx, cancel := util.CancelContext()
	defer cancel()

	r, err := fsrepo.Open(repoPath)
	if err != nil {
		return err
	}
	node, err := core.NewNode(ctx, &core.BuildCfg{Repo: r})
	if err != nil {
		return err
	}
	defer node.Close()
	_, err = assets.SeedInitDocs(node)
	if err != nil {
		return util.NewError("seeding init docs failed:", err)
	}

	return nil
}
func initializeIpnsKeyspace(repoPath string) error {
	ctx, cancel := util.CancelContext()
	defer cancel()

	r, err := fsrepo.Open(repoPath)
	if err != nil {
		return err
	}
	node, err := core.NewNode(ctx, &core.BuildCfg{Repo: r})
	if err != nil {
		return err
	}
	defer node.Close()

	emptyDir := unixfs.EmptyDirNode()
	if err := node.Pinning.Pin(ctx, emptyDir, true); err != nil {
		return err
	}
	if err := node.Pinning.Flush(ctx); err != nil {
		return err
	}

	return node.Namesys.Publish(ctx, node.PrivateKey, path.FromCid(emptyDir.Cid()))
}

func migration(repoPath string) error {
	migrationCfg, err := migrations.ReadMigrationConfig(repoPath)
	if err != nil {
		return err
	}

	fetcherFunc := func(distPath string) migrations.Fetcher {
		return ipfsfetcher.NewIpfsFetcher(distPath, 0, &repoPath)
	}
	fetchDistPath := migrations.GetDistPathEnv(migrations.CurrentIpfsDist)
	fetcher, err := migrations.GetMigrationFetcher(migrationCfg.DownloadSources, fetchDistPath, fetcherFunc)
	if err != nil {
		return err
	}
	defer fetcher.Close()

	cache := migrationCfg.Keep == "cache"
	pin := migrationCfg.Keep == "pin"
	if cache || pin {
		migrations.DownloadDirectory, err = ioutil.TempDir("", "migrations")
		if err != nil {
			return err
		}
		defer func() {
			if migrations.DownloadDirectory != "" {
				os.RemoveAll(migrations.DownloadDirectory)
			}
		}()
	}
	return migrations.RunMigration(util.NewContext(), fetcher, fsrepo.RepoVersion, "", false)
}
