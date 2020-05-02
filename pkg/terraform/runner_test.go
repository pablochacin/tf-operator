package terraform

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pablochacin/tf-operator/pkg/cmdrunner"
)

func TestCmdRunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Terraform Runner Suite")
}

type MockRunner struct {
	shellCmd string
	args     []string
    workDir  string
    env      map[string]string
}

func (r *MockRunner) Run(cmd string, args ...string) (*cmdrunner.CmdResult, error) {
	r.shellCmd = cmd
	r.args = args

	return nil, nil
}

func (r *MockRunner) SetWorkDir(path string) error {
    r.workDir = path
	return nil
}

func (r *MockRunner) SetInheritEnv(inherit bool) {
}

func (r *MockRunner) SetEnv(env map[string]string) {
    r.env = env
}

func (r *MockRunner) AddEnv(varibale string, value string) {
}

func NewMockRunner() *MockRunner {
	return &MockRunner{}
}

var _ = Describe("Terraform Runner", func() {
	var (
		err        error
		mockRunner *MockRunner
	)

    Context("Set command runner environment", func(){
        BeforeEach(func(){
            mockRunner = NewMockRunner()
			_ = NewWithCmdRunner(mockRunner, "/path/to/tfvars", "/path/to/tfconfig", "/path/to/tfstate", "/path/to/workDir")
        })

        It("Shoudl set the working directory", func(){
            Expect(mockRunner.workDir).To(Equal("/path/to/workDir"))
        })
    })

	Context("Run Init", func() {
		BeforeEach(func() {
			mockRunner = NewMockRunner()
			tfRunner := NewWithCmdRunner(mockRunner, "/path/to/tfvars", "/path/to/tfconfig", "/path/to/tfstate", "/path/to/workDir")
			err = tfRunner.Init()
		})

		It("Should not fail", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Should call terraform init", func() {
			Expect(mockRunner.shellCmd).To(Equal("terraform"))
			Expect(mockRunner.args).To(ContainElement("init"))
		})

	})

	Context("Run Apply", func() {
		BeforeEach(func() {
			mockRunner = NewMockRunner()
			tfRunner := NewWithCmdRunner(mockRunner, "/path/to/tfvars", "/path/to/tfconfig", "/path/to/tfstate", "/path/to/workDir")
			err = tfRunner.Apply()
		})

		It("Should not fail", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Should call terraform apply", func() {
			Expect(mockRunner.shellCmd).To(Equal("terraform"))
			Expect(mockRunner.args).To(ContainElement("apply"))
		})

	})
})
