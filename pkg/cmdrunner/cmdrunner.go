package cmdrunner

import (
    "fmt"
	"os"
	"os/exec"
)

// CmdEnvironment defines the environment options for executing a command
type CmdEnvironment struct {
    InheritEnv bool
    WorkDir string
    Env     map[string]string
}

// CmdResult contains the result of executing a command. Output contains the
// combined stdout and stderr outputs.
type CmdResult struct {
	ExitCode int
	Output   string
}

// Run runs a shell cmd in an execution environment
func (e *CmdEnvironment)Run(shellCmd string, args ...string) (*CmdResult, error) {
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
        for k,v := range(e.Env) {
            env = append(env, fmt.Sprintf("%s=\"%s\"", k, v))
        }
    }
    cmd.Env = env

    output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	result := &CmdResult{
                ExitCode: cmd.ProcessState.ExitCode(),
                Output: string(output),
              }

    return result, nil
}
