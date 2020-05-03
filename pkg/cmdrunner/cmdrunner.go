package cmdrunner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// Runner a command runner
type Runner interface {
	Run(shellCmd string, args ...string) (*CmdResult, error)
	SetWorkDir(path string) error
	SetInheritEnv(inherit bool)
	SetEnv(env map[string]string)
	AddEnv(varibale string, value string)
}

// CmdResult contains the result of executing a command. Output contains the
// combined stdout and stderr outputs.
type CmdResult struct {
	ExitCode int
	Output   string
}

// cmdEnvironment defines the environment options for executing a command
type cmdEnvironment struct {
	InheritEnv bool
	WorkDir    string
	Env        map[string]string
}

func New() Runner {
	return &cmdEnvironment{}
}

// Run runs a shell cmd in an execution environment
func (e *cmdEnvironment) Run(shellCmd string, args ...string) (*CmdResult, error) {
	cmd := exec.Command(shellCmd, args...)

	if e.WorkDir != "" {
		cmd.Dir = e.WorkDir
	}

	// setup environement
	env := []string{}
	if e.InheritEnv {
		env = append(env, os.Environ()...)
	}
	if e.Env != nil {
		for k, v := range e.Env {
			env = append(env, fmt.Sprintf("%s=\"%s\"", k, v))
		}
	}
	cmd.Env = env

	output, err := cmd.CombinedOutput()
	if err != nil {
		if errors.Is(err, &exec.ExitError{}) {
			return nil, err
		}
	}

	result := &CmdResult{
		ExitCode: cmd.ProcessState.ExitCode(),
		Output:   string(output),
	}

	return result, nil
}

func (e *cmdEnvironment) SetEnv(env map[string]string) {
	e.Env = env
}

func (e *cmdEnvironment) AddEnv(variable string, value string) {
	e.Env[variable] = value
}

func (e *cmdEnvironment) SetWorkDir(path string) error {
	// TODO: check path is valid
	e.WorkDir = path
	return nil
}

func (e *cmdEnvironment) SetInheritEnv(inherit bool) {
	e.InheritEnv = inherit
}
