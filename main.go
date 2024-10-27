package main

import (
	"flag"
	"fmt"
	"os"
	"slices"

	"github.com/markormesher/semver-tagger/internal/git"
	"github.com/markormesher/semver-tagger/internal/log"
	"github.com/markormesher/semver-tagger/internal/semver"
)

var usage = `
Usage: semver-tagger [-a|-M|-m|-p] [options]

-M                Bump the major version number
-m                Bump the minor version number
-p                Bump the patch version number
-rc               Add or increment the release candidate counter.
-no-rc            Remove the release candidate counter.
-init             Create the initial v0.1.0 version.

-v, --verbose     Verbose logging
-f, --force       Create a tag even no matter what
-y, --no-confirm  Create the tag without confirming
-P, --push        Run 'git push --tags' after creating the tag
`

// TODO: actually detect this from the repo (using refs/remotes/origin/HEAD doesn't work on local-only repos)
var defaultBranches = []string{"main", "master", "develop"}

func main() {
	// config

	var majorFlag, minorFlag, patchFlag, rcFlag, noRcFlag, initFlag, verboseFlag, forceFlag, noConfirmFlag, pushFlag bool
	flag.BoolVar(&majorFlag, "M", false, "")
	flag.BoolVar(&minorFlag, "m", false, "")
	flag.BoolVar(&patchFlag, "p", false, "")
	flag.BoolVar(&rcFlag, "rc", false, "")
	flag.BoolVar(&noRcFlag, "no-rc", false, "")
	flag.BoolVar(&initFlag, "init", false, "")
	flag.BoolVar(&verboseFlag, "v", false, "")
	flag.BoolVar(&forceFlag, "f", false, "")
	flag.BoolVar(&forceFlag, "force", false, "")
	flag.BoolVar(&noConfirmFlag, "y", false, "")
	flag.BoolVar(&noConfirmFlag, "no-confirm", false, "")
	flag.BoolVar(&pushFlag, "P", false, "")
	flag.BoolVar(&pushFlag, "push", false, "")

	flag.Usage = func() { fmt.Println(usage) }
	flag.Parse()

	if verboseFlag {
		log.Verbose = true
	}

	qtyTagTypeFlags := 0
	if majorFlag {
		qtyTagTypeFlags++
	}
	if minorFlag {
		qtyTagTypeFlags++
	}
	if patchFlag {
		qtyTagTypeFlags++
	}

	if qtyTagTypeFlags == 0 {
		log.Error("Invalid usage: must specify one of -M / -m / -p")
		os.Exit(1)
	}

	if qtyTagTypeFlags > 1 {
		log.Error("Invalid usage: cannot specify more than one of -M / -m / -p")
		os.Exit(1)
	}

	if qtyTagTypeFlags > 0 && initFlag {
		log.Error("Invalid usage: cannot specify -init with any other version flags")
		os.Exit(1)
	}

	if rcFlag && noRcFlag {
		log.Error("Invalid usage: cannot specify both -rc and -no-rc")
		os.Exit(1)
	}

	// validate the repo state

	repoClean, err := git.RepoIsClean()
	if err != nil {
		log.Error("%v", err)
		os.Exit(1)
	}
	if !repoClean {
		if forceFlag {
			log.Warn("Repo is not clean, but continuing because force flag is specified")
		} else {
			log.Error("Repo is not clean")
			os.Exit(1)
		}
	}

	currentBranch, err := git.CurrentBranch()
	if err != nil {
		log.Error("%v", err)
		os.Exit(1)
	}
	if !slices.Contains(defaultBranches, currentBranch) {
		if forceFlag {
			log.Warn("Current branch (%s) doesn't look like a default branch, but continuing because force flag is specified", currentBranch)
		} else {
			log.Error("Current branch (%s) doesn't look like a default branch", currentBranch)
			os.Exit(1)
		}
	}

	// get the current version

	description, err := git.Describe()
	if err != nil && !initFlag {
		log.Error("%v", err)
		os.Exit(1)
	}

	var currentVer semver.SemVer
	if initFlag {
		if description != "" {
			log.Error("Refusing to create initial version when other versions already exist")
			os.Exit(1)
		} else {
			// set current version to v0.0.0 and trigger a minor version bump
			currentVer = semver.SemVer{Prefix: "v"}
			minorFlag = true
		}
	} else {
		currentVer, err = semver.FromString(description)
		if err != nil {
			log.Error("%v", err)
			os.Exit(1)
		}
	}

	if currentVer.CommitDistance == 0 && !initFlag {
		if forceFlag {
			log.Warn("There have been no commits since the last tag, but continuing because force flag is specified")
		} else {
			log.Info("There have been no commits since the last tag")
			os.Exit(0)
		}
	}

	// finally, decide the new version and tag it

	newVer := currentVer.Bump(majorFlag, minorFlag, patchFlag, rcFlag, noRcFlag)

	log.Info("Prev:  %s", currentVer.String())
	log.Info("New:   %s", newVer.String())

	err = git.CreateTag(&newVer, noConfirmFlag)
	if err != nil {
		log.Error("%v", err)
		os.Exit(1)
	}
	log.Info("Created tag")

	if pushFlag {
		err = git.PushTags()
		if err != nil {
			log.Error("%v", err)
			os.Exit(1)
		}
		log.Info("Pushed tags")
	}
}
