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

-a                Detect the tag type based on commit messages since the last tag (default behaviour)
-M                Bump the major version number
-m                Bump the minor version number
-p                Bump the patch version number
-rc               Create a release candidate tag
-init             Create the initial tag (v0.1.0)

-v, --verbose     Verbose logging
-f, --force       Create a tag even no matter what
-y, --no-confirm  Create the tag without confirming
`

// TODO: actually detect this from the repo (using refs/remotes/origin/HEAD doesn't work on local-only repos)
var defaultBranches = []string{"main", "master", "develop"}

func main() {
	// config
	autoFlag := flag.Bool("a", false, "")
	majorFlag := flag.Bool("M", false, "")
	minorFlag := flag.Bool("m", false, "")
	patchFlag := flag.Bool("p", false, "")
	rcFlag := flag.Bool("rc", false, "")
	initFlag := flag.Bool("init", false, "")
	verboseFlag := flag.Bool("v", false, "")
	forceFlag := flag.Bool("f", false, "")
	forceFlag = flag.Bool("force", false, "")
	noConfirmFlag := flag.Bool("y", false, "")
	noConfirmFlag = flag.Bool("no-confirm", false, "")

	flag.Usage = func() { fmt.Println(usage) }
	flag.Parse()

	if *verboseFlag {
		log.Verbose = true
	}

	// ensure only one version type flag was passed
	qtyTagTypeFlags := 0
	if *autoFlag {
		qtyTagTypeFlags++
	}
	if *majorFlag {
		qtyTagTypeFlags++
	}
	if *minorFlag {
		qtyTagTypeFlags++
	}
	if *patchFlag {
		qtyTagTypeFlags++
	}
	if qtyTagTypeFlags > 1 {
		log.Error("Invalid usage: cannot pass more than one of -a / -M / -m / -p")
		os.Exit(1)
	} else if qtyTagTypeFlags == 0 {
		log.Debug("No tag set; assuming -a")
		*autoFlag = true
	}

	// validate the repo state

	repoClean, err := git.RepoIsClean()
	if err != nil {
		log.Error("%v", err)
		os.Exit(1)
	}
	if !repoClean {
		if *forceFlag {
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
		if *forceFlag {
			log.Warn("Current branch (%s) doesn't look like a default branch, but continuing because force flag is specified", currentBranch)
		} else {
			log.Error("Current branch (%s) doesn't look like a default branch", currentBranch)
			os.Exit(1)
		}
	}

	// get the current version

	description, err := git.Describe()
	if err != nil {
		if *initFlag {
			log.Info("Creating the initial tag")
			err = git.CreateTag(&semver.SemVer{Prefix: "v", Major: 0, Minor: 1, Patch: 0}, *noConfirmFlag)
			if err != nil {
				log.Error("%v", err)
				os.Exit(1)
			}
		} else {
			log.Error("%v", err)
			os.Exit(1)
		}
	}

	currentVer, err := semver.FromString(description)
	if err != nil {
		log.Error("%v", err)
		os.Exit(1)
	}

	if currentVer.CommitDistance == 0 {
		if *forceFlag {
			log.Warn("There have been no commits since the last tag, but continuing because force flag is specified")
		} else {
			log.Info("There have been no commits since the last tag")
			os.Exit(0)
		}
	}

	// finally, decide the new version and tag it

	newVer := currentVer

	if *autoFlag {
		// TODO: build this
		log.Error("Auto-tag isn't built yet. Sorry")
		os.Exit(1)
	}

	switch {
	case *majorFlag:
		newVer.Major++
		newVer.Minor = 0
		newVer.Patch = 0

	case *minorFlag:
		newVer.Minor++
		newVer.Patch = 0

	case *patchFlag:
		newVer.Patch++
	}

	if *rcFlag {
		newVer.Rc++
	}

	log.Info("Prev:  %s", currentVer.String())
	log.Info("New:   %s", newVer.String())

	err = git.CreateTag(&newVer, *noConfirmFlag)
	if err != nil {
		log.Error("%v", err)
		os.Exit(1)
	}
}
