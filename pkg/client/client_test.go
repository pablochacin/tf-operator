package client

import (
    "io/ioutil"
    "os"
    "path/filepath"
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

    main_tf = `
variable "greetee" {
  type = string
}

output "greetings" {
    value = "Hello ${var.greetee}"
}
`

    terraform_tfvars = `
greetee = "World"
`

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

// createTfWorkDir creates a temporary dir with the terraform configuration files
// the files map has the path and content for each file. 
func createTfWorkDir(files map[string]string) (string, error) {
   tfDir, err := ioutil.TempDir("", "terraform")
   if err != nil {
       return "", err
   }

    // for each file, create the parent path and write the content
    for filePath, content := range(files) {
        basePath, fileName := filepath.Split(filePath)
        absPath := filepath.Join(tfDir, basePath)
        err = os.MkdirAll(absPath, os.ModePerm)
        if err != nil {
           return "", err
        }
        // it may be only a dir must be created, not a file. Ignore content
        if fileName != "" {
            err = ioutil.WriteFile(filepath.Join(absPath, fileName), []byte(content), 0666)
            if err != nil {
                return "", err
            }
        }
    }

    return tfDir, nil
}

var _ = Describe("Client", func() {
	var (
        stack       *tfo.Stack
		err         error
        rc          ctl.Client
        tfDir       string
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

        Context("tf files cannot be accessed", func() {
            var tfFiles map[string]string

            JustBeforeEach(func() {
                rc = newFakeClient()
                c, _ := NewFromRuntimeClient(rc)
                tfDir, err =  createTfWorkDir(tfFiles)
                tfvars := filepath.Join(tfDir,  "terraform.tfvars")
                tfconf := filepath.Join(tfDir, "tfconfig")
                stack, err = c.CreateStack(stackName, namespace, tfconf, tfvars)
            })

            AfterEach(func(){
                os.RemoveAll(tfDir)
            })

            Context("Config directory doesn't exists", func(){
                BeforeEach(func(){
                    // create tfvars but not tfconfig
                    tfFiles = map[string]string{
                        "terraform.tfvars": terraform_tfvars,
                    }
                })

                It("Should Return an error", func() {
                    Expect(err).To(HaveOccurred())
                    Expect(Is(err,ErrorReasonFileCanNotBeAccessed)).To(BeTrue())
                    Expect(stack).To(BeNil())
                })
           })

           Context("tfvars file doesn't exists", func(){
                BeforeEach(func(){
                    // create tfconf files but not tfvars
                    tfFiles = map[string]string{
                        "tfconfig/main.tf": main_tf,
                    }
                })

                It("Should Return an error", func() {
                    Expect(err).To(HaveOccurred())
                    Expect(Is(err,ErrorReasonFileCanNotBeAccessed)).To(BeTrue())
                    Expect(stack).To(BeNil())
                })
           })

           Context("tfconfig dir is empty", func(){
                BeforeEach(func(){
                    // create empty tfconf dir  
                    tfFiles = map[string]string{
                        "terraform.tfvars": "",
                        "tfconfig/": "",
                    }
                })

                It("Should Return an error", func() {
                    Expect(err).To(HaveOccurred())
                    Expect(Is(err,ErrorReasonFileCanNotBeAccessed)).To(BeTrue())
                    Expect(stack).To(BeNil())
                })
           })

           Context("tfvars file is empty", func(){
                BeforeEach(func(){
                    // create tfconf but not tfvars
                    tfFiles = map[string]string{
                        "terraform.tfvars": "",
                        "tfconfig/main.tf": main_tf,
                    }
                })

                It("Should Return an error", func() {
                    Expect(err).To(HaveOccurred())
                    Expect(Is(err,ErrorReasonInvalidFileContent)).To(BeTrue())
                    Expect(stack).To(BeNil())
                })
           })

           Context("tfconfig file is empty", func(){
                BeforeEach(func(){
                    // create tfconf but not tfvars
                    tfFiles = map[string]string{
                        "terraform.tfvars": terraform_tfvars,
                        "tfconfig/main.tf": "",
                    }
                })

                It("Should Return an error", func() {
                    Expect(err).To(HaveOccurred())
                    Expect(Is(err,ErrorReasonInvalidFileContent)).To(BeTrue())
                    Expect(stack).To(BeNil())
                })
           })
        })
    })
})
