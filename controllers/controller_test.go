package controllers

import (
    "context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	tfo "github.com/pablochacin/tf-operator/api/v1alpha1"
    batchv1 "k8s.io/api/batch/v1"
    corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    rmt "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)


var (
    stackName   = "stack-name"
    namespace   = "test"
    tfconfigMap = map[string]string{
        "main.tf":  `
variable "greetee" {
  type = string
}

output "greetings" {
    value = "Hello ${var.greetee}"
}
`,
    }
    terraformTfvars = `
greetee = "World"
`

)

// createStack creates a Stack from a tfconfig ConfigMap and a tfvars Secret
func createStack(name string, ns string, tfconfig *corev1.ConfigMap, tfvars *corev1.Secret) *tfo.Stack {
    return &tfo.Stack{
        ObjectMeta: metav1.ObjectMeta{
            Name: name,
            Namespace: ns,
        },
        Spec: tfo.StackSpec{
            TfConfig: corev1.LocalObjectReference{Name: tfconfig.Name},
            TfVars:   corev1.LocalObjectReference{Name: tfvars.Name},
        },
    }
}

// createTfConfigMap creates a configMap from a Map with the names and content
// of tf files
func createTfConfigMap(name string, ns string, content map[string]string) *corev1.ConfigMap {
    return &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name: name,
            Namespace: ns,
        },
        Data: content,
    }
}

// createTfvarsSecret creates a Secret from the content of a tfvars
func createTfvarsSecret(name string, ns string, content string) *corev1.Secret {
    return &corev1.Secret{
        ObjectMeta: metav1.ObjectMeta{
            Name: name,
            Namespace: ns,
        },
        Data: map[string][]byte{
            "terraform.tfvars": []byte(content),
        },
    }
}

var _ = Describe("Controller", func() {
	var (
        stack      *tfo.Stack
        request    ctrl.Request
        result     ctrl.Result
        reconciler *StackReconciler
		err        error
        initObjs = []rmt.Object{}
    )

    Context("Reconciliate Stack", func(){

        JustBeforeEach(func() {
            for _, obj  := range(initObjs) {
                err := k8sClient.Create(context.TODO(), obj)
                Expect(err).NotTo(HaveOccurred())
            }

            reconciler = &StackReconciler {
                Client: k8sClient,
                Log:    ctrl.Log.WithName("controllers").WithName("Stack"),
            }

            result, err = reconciler.Reconcile(request)
        })

        AfterEach(func(){
            // clear objects from cluster
            for _, obj := range(initObjs) {
                k8sClient.Delete(context.TODO(), obj)
            }
            initObjs = []rmt.Object{}
         })

        Context("stack created", func() {
            BeforeEach(func() {
                tfvars := createTfvarsSecret(stackName, namespace, terraformTfvars)
                tfconfig := createTfConfigMap(stackName, namespace, tfconfigMap)
                stack = createStack(stackName, namespace, tfconfig, tfvars)
                initObjs = append(initObjs,stack, tfvars, tfconfig)
                request = ctrl.Request{
                    types.NamespacedName{
                        Name: stack.Name,
                        Namespace: stack.Namespace,
                    },
                }
            })


            It("Should not Return an error", func() {
                Expect(err).NotTo(HaveOccurred())
                Expect(result).NotTo(BeNil())
                jobs := &batchv1.JobList{}
                Expect(jobs.Size()).To(Equal(int(1)))
            })
        })
    })
})
