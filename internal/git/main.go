package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/markormesher/semver-tagger/internal/log"
)

func execCmd(cmd string) (string, error) {
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Debug("Error running command: %v", string(out))
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func RepoIsClean() (bool, error) {
	status, err := execCmd("git status --porcelain")
	if err != nil {
		return false, fmt.Errorf("error checking whether workspace is clean: %w", err)
	}

	return status == "", nil
}

func CurrentBranch() (string, error) {
	branch, err := execCmd("git rev-parse --abbrev-ref HEAD")
	if err != nil {
		return "", fmt.Errorf("error checking current branch: %w", err)
	}

	return branch, nil
}
