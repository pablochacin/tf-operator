package cmdrunner

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

)

func TestCmdRunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Command Runner Suite")
}

var _ = Describe("Command Runner", func() {
	var (
        result  *CmdResult
        err     error
	)


	Context("Run Job with default enviroment", func() {

        BeforeEach(func() {
		    env := &CmdEnvironment{}
	        result, err = env.Run("echo", "-n", "testing echo")
	    })

		It("Should not fail", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result).ToNot(BeNil())
		})

        It("Should return non error exit code", func(){
            Expect(result.ExitCode).To(Equal(0))
        })

		It("should return stdout", func() {
			Expect(result.Output).To(Equal("testing echo"))
		})
	})

    Context("Capture stderr", func() {
        BeforeEach(func(){
            env := &CmdEnvironment{}
            result, err = env.Run("sh", "-c", "echo -n >&2 'testing echo'")
        })

		It("Should not fail", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result).NotTo(BeNil())
		})

        It("Should capture stderr", func(){
            Expect(result.Output).To(Equal("testing echo"))
        })
    })

    Context("Run job with environment variables", func(){
	    BeforeEach(func() {
		    env :=  &CmdEnvironment{
                        InheritEnv: false,
                        Env: map[string]string{"FOO": "BAR"},
                    }
            result, err = env.Run("env")
	    })

        It("Should not fail", func(){
            Expect(err).To(BeNil())
            Expect(result).NotTo(BeNil())
        })

        It("Should return env variables", func() {
            Expect(result.Output).To(Equal("FOO=\"BAR\"\n"))
        })
    })

    Context("Capture return code", func(){
        BeforeEach(func(){
            env := &CmdEnvironment{}
            result, err = env.Run("/bin/false")
        })

        It("Should not fail", func(){
            Expect(err).To(BeNil())
            Expect(result).NotTo(BeNil())
        })

        It("Should return exit code", func() {
            Expect(result.ExitCode).To(Equal(1))
        })
     })

    Context("Changes the working dir", func(){
        BeforeEach(func(){
            env := &CmdEnvironment{WorkDir: "/tmp"}
            result, err = env.Run("pwd")
        })

        It("Should not fail", func(){
            Expect(err).To(BeNil())
            Expect(result).NotTo(BeNil())
        })

        It("Should run in working directory", func() {
            Expect(result.Output).To(Equal("/tmp\n"))
        })
 
    })
 })
