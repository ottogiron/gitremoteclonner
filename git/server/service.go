package server

import (
	"strings"

	"os/exec"

	"sync"

	"io"

	"github.com/pkg/errors"
)

const (
	successStatus = 0
	failStatus    = 1
)

//GitService operations in the git server
type GitService interface {
	Execute(dir string, command string, output io.Writer) error
}

var _ GitService = (*gitService)(nil)

type gitService struct {
	mu sync.Mutex
}

//NewGitService returns a new GitService
func NewGitService() GitService {
	return &gitService{}
}

func (g *gitService) Execute(dir string, command string, output io.Writer) error {

	command = strings.Replace(command, "\\", "", -1)

	commandTokens := strings.Fields(command)

	if len(command) == 0 {
		return errors.Errorf("Invalid empty command %s", command)
	}

	if commandTokens[0] != "git" {
		return errors.Errorf("Invalid non git command: %s", command)
	}

	cmd := exec.Command("sh", "-c", command)

	cmd.Dir = dir
	cmd.Stdout = output
	cmd.Stderr = output
	err := cmd.Run()

	if err != nil {
		return errors.Wrapf(err, "Failed to run command %s", command)
	}
	return nil
}
