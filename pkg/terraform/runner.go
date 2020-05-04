package terraform

import (
	"path"

	"github.com/pablochacin/tf-operator/pkg/cmdrunner"
)

type TfRunner interface {
	Apply() error
}

// TfWorkspace defines the working environment for the Terraform Runner
type TfWorkspace struct {
	runner   cmdrunner.Runner
	tfvars   string
	tfconfig string
	tfstate  string
	workDir  string
}

// NewWithCmdRunner builds a TfWorkspace with a given command runner
func NewWithCmdRunner(runner cmdrunner.Runner, tfvars string, tfconfig string, tfstate string, workDir string) *TfWorkspace {
	runner.SetWorkDir(workDir)
	return &TfWorkspace{
		runner:   runner,
		tfvars:   tfvars,
		tfconfig: tfconfig,
		tfstate:  tfstate,
		workDir:  workDir,
	}
}

// New builds a TfWorkspace with a default command runner
func New(tfvars string, tfconfig string, tfstate string, workDir string) *TfWorkspace {
	return NewWithCmdRunner(cmdrunner.New(), tfvars, tfconfig, tfstate, workDir)
}

// Init initializes terraform
func (w *TfWorkspace)Init() error {
	args := []string{"init",
		" -input=false",
	}

	_, err := w.runner.Run("terraform", args...)
	if err != nil {
		return err
	}

	return nil

}

// Apply applies terraform plan
func (w *TfWorkspace) Apply() error {
	args := []string{"apply",
		"-input=false",
		"-auto-aprove",
		"-var-file", w.tfvars,
		"-state", w.tfstate,
		"-state-out", path.Join(w.workDir, "terraform.tfstate"),
	}
	_, err := w.runner.Run("terraform", args...)
	if err != nil {
		return err
	}

	return nil
}
