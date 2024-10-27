package git

import (
	"fmt"
	"os/exec"

	"github.com/markormesher/semver-tagger/internal/log"
)

func execCmd(cmd string) (string, error) {
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Debug("Error running command: %v", string(out))
		return "", err
	}

	return string(out), nil
}

func RepoIsClean() (bool, error) {
	status, err := execCmd("git status --porcelain")
	if err != nil {
		return false, fmt.Errorf("error checking whether workspace is clean: %w", err)
	}

	return status == "", nil
}

func DefaultBranch() (string, error) {
	branch, err := execCmd("basename $(git symbolic-ref --short refs/remotes/origin/HEAD)")
	if err != nil {
		return "", fmt.Errorf("error checking default branch: %w", err)
	}

	return branch, nil
}

func CurrentBranch() (string, error) {
	branch, err := execCmd("git rev-parse --abbrev-ref HEAD")
	if err != nil {
		return "", fmt.Errorf("error checking current branch: %w", err)
	}

	return branch, nil
}
