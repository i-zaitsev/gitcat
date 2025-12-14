package gitclone

import (
	"os/exec"

	"github.com/i-zaitsev/gitcat/pkg/gitpath"
)

// Clone clones a git repository via SSH or HTTPS.
// For SSH, it is assumed that the SSH key is properly configured.
// For HTTPS, only public repositories are supported (no authentication).
func Clone(repoUrl *gitpath.GitPath, localDir string) error {
	cmd := exec.Command("git", "clone", repoUrl.Path, localDir)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// Status checks the git status of a repository.
func Status(repoDir string) (string, error) {
	cmd := exec.Command("git", "status")
	cmd.Dir = repoDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
