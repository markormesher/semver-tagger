package git

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/markormesher/semver-tagger/internal/log"
	"github.com/markormesher/semver-tagger/internal/semver"
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

func Describe() (string, error) {
	describe, err := execCmd("git describe --tags")
	if err != nil {
		return "", fmt.Errorf("error describing repo: %w", err)
	}

	return describe, nil
}

func CreateTag(tag *semver.SemVer, noConfirm bool) error {
	if !noConfirm {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Create tag? [Yn] ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" && input != "y" && input != "Y" {
			log.Info("Aborting")
			os.Exit(0)
		}
	}

	_, err := execCmd(fmt.Sprintf("git tag -m '%s' '%s'", tag.String(), tag.String()))
	if err != nil {
		return fmt.Errorf("error creating tag: %w", err)
	}

	return nil
}

func PushTags() error {
	_, err := execCmd("git push --tags")
	if err != nil {
		return fmt.Errorf("error pushing tags: %w", err)
	}

	return nil
}
