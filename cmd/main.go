package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"

	"github.com/markormesher/semver-tagger/internal/git"
	"github.com/markormesher/semver-tagger/internal/log"
	"github.com/markormesher/semver-tagger/internal/semver"
)

var usage = `
Usage: semver-tagger [-a|-M|-m|-p] [options]

-a                Automatically determine the tag type (default behaviour)
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

var patchCommitPattern = *regexp.MustCompile(`^(Merge pull request ')?(chore|fix|ci)(\([\w \-]+\))?:`)

func bold(str string) string {
	return "\033[1m" + str + "\033[0m"
}

func main() {
	// config

	var autoFlag, majorFlag, minorFlag, patchFlag, rcFlag, noRcFlag, initFlag, verboseFlag, forceFlag, noConfirmFlag, pushFlag bool
	flag.BoolVar(&autoFlag, "a", false, "")
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
	if autoFlag {
		qtyTagTypeFlags++
	}
	if majorFlag {
		qtyTagTypeFlags++
	}
	if minorFlag {
		qtyTagTypeFlags++
	}
	if patchFlag {
		qtyTagTypeFlags++
	}

	if qtyTagTypeFlags == 0 && !rcFlag && !initFlag {
		autoFlag = true
		log.Info("No tag type specified; assuming -a for automatic")
	}

	if qtyTagTypeFlags > 1 {
		log.Error("Invalid usage: cannot specify more than one of -a / -M / -m / -p")
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

	// determine the tag type if we're in auto mode

	if autoFlag {
		if currentVer.Rc > 0 {
			log.Debug("Latest tag is an RC; creating a RC tag")
			rcFlag = true
		} else {
			commitMessages, err := git.CommitsSinceLastTag()
			if err != nil {
				log.Error("%v", err)
				os.Exit(1)
			}

			allPatchCommits := true
			for _, msg := range commitMessages {
				if !patchCommitPattern.MatchString(msg) {
					allPatchCommits = false
					break
				}
			}

			if allPatchCommits {
				log.Debug("All commits since the latest tag are patch commits; creating a patch tag")
				patchFlag = true
			} else {
				log.Debug("Found one or more non-patch commits since the latest tag; creating a minor tag")
				minorFlag = true
			}
		}
	}

	// finally, decide the new version and tag it

	newVer := currentVer.Bump(majorFlag, minorFlag, patchFlag, rcFlag, noRcFlag)

	log.Info("Prev:  %s", bold(currentVer.String()))
	log.Info("New:   %s", bold(newVer.String()))

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
