package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/markormesher/semver-tagger/internal/git"
	"github.com/markormesher/semver-tagger/internal/log"
)

var usage = `
Usage: semver-tagger [-a|-M|-m|-p] [options]

-a                Detect the tag type based on commit messages since the last tag (default behaviour)
-M                Tag a new major version
-m                Tag a new minor version
-p                Tag a new patch version

-f, --force       Create a tag even if the repo is not clean or is on a non-default branch
-v, --verbose     Verbose logging
`

func main() {
	// config
	autoFlag := flag.Bool("a", false, "")
	majorFlag := flag.Bool("M", false, "")
	minorFlag := flag.Bool("m", false, "")
	patchFlag := flag.Bool("p", false, "")
	verboseFlag := flag.Bool("v", false, "")
	forceFlag := flag.Bool("f", false, "")
	forceFlag = flag.Bool("force", false, "")

	flag.Usage = func() { fmt.Println(usage) }
	flag.Parse()

	// validate config
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

	if *verboseFlag {
		log.Verbose = true
	}

	repoClean, err := git.RepoIsClean()
	if err != nil {
		log.Error("Failed to check whether repo is clean: %v", err)
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

	defaultBranch, err := git.DefaultBranch()
	if err != nil {
		log.Error("Failed to check default branch: %v", err)
		os.Exit(1)
	}
	currentBranch, err := git.CurrentBranch()
	if err != nil {
		log.Error("Failed to check current branch: %v", err)
		os.Exit(1)
	}
	if defaultBranch != currentBranch {
		if *forceFlag {
			log.Warn("Repo is not on the default branch, but continuing because force flag is specified")
		} else {
			log.Error("Repo is not on the default branch")
			os.Exit(1)
		}
	}

}
