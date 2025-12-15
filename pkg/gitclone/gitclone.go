package gitclone

import (
	"os/exec"

	"github.com/i-zaitsev/gitcat/pkg/gitpath"
	"github.com/i-zaitsev/gitcat/pkg/log"
)

// Clone clones a git repository via SSH or HTTPS.
// For SSH, it is assumed that the SSH key is properly configured.
// For HTTPS, only public repositories are supported (no authentication).
func Clone(repoUrl *gitpath.GitPath, localDir string) error {
	log.Debug("cloning repository", "url", repoUrl.Path, "dir", localDir)
	cmd := exec.Command("git", "clone", repoUrl.Path, localDir)
	if err := cmd.Run(); err != nil {
		log.Error("git clone failed", "error", err, "url", repoUrl.Path)
		return err
	}
	log.Debug("repository cloned successfully", "dir", localDir)
	return nil
}

// Status checks the git status of a repository.
func Status(repoDir string) (string, error) {
	log.Debug("checking git status", "dir", repoDir)
	cmd := exec.Command("git", "status")
	cmd.Dir = repoDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("git status failed", "error", err, "dir", repoDir)
		return "", err
	}
	return string(output), nil
}
