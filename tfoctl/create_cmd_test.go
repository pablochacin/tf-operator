package main

import (
	"bytes"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/spf13/cobra"
)

func TestCreateCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tfoctl create Suite")
}

const (
	stackName = "stack-name"
)

// dummyRunE dummy function that prevents the command to be executed
func dummyRunE(c *cobra.Command, args []string) error {
	return nil
}

var _ = Describe("run create command", func() {
	var (
		cmd    *cobra.Command
		output *bytes.Buffer
	)

	BeforeEach(func() {
		cmd = newCreateCmd()
		output = new(bytes.Buffer)
		cmd.SetOutput(output)
		cmd.RunE = dummyRunE
	})

	Context("Create with default values", func() {
		BeforeEach(func() {
			cmd.SetArgs([]string{"-s", stackName})
			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should have default Values", func() {
			defaults := map[string]string{
				"config":    "./",
				"namespace": "default",
				"map":       "",
				"state":     "",
				"vars":      "terraform.tfvars",
			}
			for flagName, value := range defaults {
				flag := cmd.Flags().Lookup(flagName)
				Expect(flag).ShouldNot(BeNil())
				Expect(flag.Value.String()).Should(Equal(value))
			}
		})

		It("Should have stack name set", func() {
			flag := cmd.Flags().Lookup("stack")
			Expect(flag).ShouldNot(BeNil())
			Expect(flag.Changed).To(BeTrue())
			Expect(flag.Value.String()).Should(Equal(stackName))
		})
	})

	Context("Create with invalid values", func() {
		var (
			err error
		)

		Context("Created without required parameters", func() {
			BeforeEach(func() {
				cmd.SetArgs([]string{})
				err = cmd.Execute()
			})

			It("Should fail", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("Created with conflicting parameters", func() {
			BeforeEach(func() {
				cmd.SetArgs([]string{"-s", "my stack", "-m", "mymap", "-c", "myConfDir"})
				err = cmd.Execute()
			})

			It("Should fail", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
