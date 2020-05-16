package client

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	tfo "github.com/pablochacin/tf-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    rmt "k8s.io/apimachinery/pkg/runtime"
    ctl "sigs.k8s.io/controller-runtime/pkg/client"
    fake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)


const (
    stackName   = "stack-name"
    namespace   = "test"
)



func TestCreate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tf-Operator Client Suite")
}


// newFakeClient Builds a fake controller runtime initialized with tfo scheme and
// optional list of initial objects
func newFakeClient(objs...rmt.Object) ctl.Client {
    sch := rmt.NewScheme()
    tfo.AddToScheme(sch)
    rc := fake.NewFakeClientWithScheme(sch, objs...)
    return rc
}

var _ = Describe("Client", func() {
	var (
        stack       *tfo.Stack
		err         error
        rc          ctl.Client
)

    Context("Create Stack", func(){
        Context("stack already exists", func() {
            BeforeEach(func() {
                rc = newFakeClient(
                    &tfo.Stack{
                        ObjectMeta: metav1.ObjectMeta{
                            Name: stackName,
                            Namespace: namespace,
                        },
                })
                c, _ := NewFromRuntimeClient(rc)
                stack, err = c.CreateStack(stackName, namespace,"","")
            })

            It("Should Return an error", func() {
                Expect(err).To(HaveOccurred())
                Expect(Is(err,ErrorReasonAlreadyExists)).To(BeTrue())
                Expect(stack).To(BeNil())
            })
        })
    })
})
