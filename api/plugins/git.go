package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"gopkg.in/src-d/go-billy.v4"

	"github.com/blang/semver"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	"github.com/Zenika/marcel/api/db/plugins"
	"github.com/Zenika/marcel/config"
)

const (
	ErrPluginNotFound errPluginNotFound = "NO_PLUGIN_FOUND"
	masterRef                           = "refs/heads/master"
)

var (
	master = Version{ReferenceName: masterRef}
)

type errPluginNotFound string

func (err errPluginNotFound) Error() string {
	return string(err)
}

// FetchVersionsFromGit returns a sorted list of versions found in the remote tag list
func FetchVersionsFromGit(url string) (Versions, error) {
	repo, err := CloneGitRepository(url, "", nil)
	if err != nil {
		return nil, fmt.Errorf("Error while cloning %s: %s", url, err)
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		return nil, fmt.Errorf("Error retrieving origin remote from %s: %s", url, err)
	}

	log.Debug("Fetching tags...")
	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Error fetching tags from %s: %s", url, err)
	}

	var versions Versions
	for _, ref := range refs {
		name := ref.Name()
		if name.IsTag() {
			if version, err := semver.ParseTolerant(name.Short()); err != nil {
				log.Debugf("Ignoring non semver tag: %s", name.Short())
			} else {
				versions = append(versions, Version{name, version})
			}
		}
	}

	sort.Sort(versions)

	return versions, nil
}

// FetchManifestFromGit reads the marcel's manifest file from the given repository
func FetchManifestFromGit(repo *git.Repository, ref plumbing.ReferenceName) (*plugins.Plugin, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("Error while getting WorkTree : %s", err)
	}

	if err = wt.Checkout(&git.CheckoutOptions{Branch: ref}); err != nil {
		return nil, fmt.Errorf("Error while checking out manifest: %s", err)
	}

	manifest, err := wt.Filesystem.Open("marcel.json")
	if err != nil {
		return nil, fmt.Errorf("Error while opening manifest : %s", err)
	}
	defer manifest.Close()

	plugin := &plugins.Plugin{}
	if err := json.NewDecoder(manifest).Decode(plugin); err != nil {
		return nil, fmt.Errorf("Error while reading manifest : %s", err)
	}

	return plugin, nil
}

// CloneGitRepository returns a repo initialised for url and checked out for ref.
// ref can be omited to fetch default remote's HEAD
// fs can be omitted to avoid checking out repository's content
func CloneGitRepository(url string, ref plumbing.ReferenceName, fs billy.Filesystem) (*git.Repository, error) {
	if fs != nil {
		log.Debugf("Cloning %s (%s) into %s ...", url, ref.Short(), fs.Root())
	} else {
		log.Debugf("Cloning %s (%s)...", url, ref.Short())
	}

	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:           url,
		SingleBranch:  true,
		NoCheckout:    true,
		Depth:         1,
		Tags:          git.NoTags,
		ReferenceName: ref,
	})
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// FetchFromGit returns the plugin found in the git repo pointed by url
// It also returns the fullpath of the temporary directory where the plugin's repo content is stored
// The caller should take care of the temporary directory removal
func FetchFromGit(url string) (plugin *plugins.Plugin, tempDir string, err error) {

	versions, err := FetchVersionsFromGit(url)
	if err != nil {
		return nil, tempDir, fmt.Errorf("Error while retreiving versions: %s", err)
	}

	latest, err := versions.Last()
	if err != nil {
		latest = master
		log.Warnf("No versions were found on %s. Using default reference (%s)", url, latest.Short())
	}

	tempDir, err = ioutil.TempDir(config.Config.PluginsPath, "new_plugin")
	if err != nil {
		return nil, tempDir, fmt.Errorf("Error while trying to create temporary directory: %s", err)
	}

	repo, err := CloneGitRepository(url, latest.ReferenceName, osfs.New(tempDir))
	if err != nil {
		return nil, tempDir, fmt.Errorf("Error while cloning %s into %s : %s", latest.Short(), tempDir, err)
	}

	log.Debug("Checking out manifest...")

	plugin, err = FetchManifestFromGit(repo, latest.ReferenceName)
	if err != nil {
		return nil, tempDir, fmt.Errorf("Error while fetching manifest: %s", err)
	}

	plugin.URL = url
	for _, version := range versions {
		plugin.Versions = append(plugin.Versions, version.String())
	}

	return plugin, tempDir, nil
}
