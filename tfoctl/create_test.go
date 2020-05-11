package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	tfo "github.com/pablochacin/tf-operator/api/v1alpha1"
)

func TestCreate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tfoctl create Suite")
}

// fakeClient mocks the Client interface for testing
type fakeClient struct {
	// if not nil, error to return in the invocation
	err error

	// if err not set, stack to return
	stack *tfo.Stack
}

// GetStack return a stack or an error set in the fakeClient struct
func (c *fakeClient) GetStack(stackName string, namespace string) (*tfo.Stack, error) {

	if c.err != nil {
		return nil, c.err
	}

	return c.stack, nil
}

// withError sets the error to return on
func (c *fakeClient) withError(err error) {
	c.err = err
}

var _ = Describe("create", func() {
	var (
		stackName = "stack-name"
		opts      *createOpts
		err       error
	)

	Context("stack already exists", func() {
		BeforeEach(func() {
			opts = &createOpts{
				stack:     stackName,
				namespace: "default",
				// Return stack TODO: complete stack's fields
				client: &fakeClient{
					stack: &tfo.Stack{},
				},
			}

			err = opts.run()
		})

		It("Should Return an error", func() {
			Expect(err).To(HaveOccurred())
		})
	})
})
