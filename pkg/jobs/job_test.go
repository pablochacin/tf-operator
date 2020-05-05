package jobs

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestJobs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jobs Suite")
}

func getVolumeSources(volumes []corev1.Volume) []string {
	sources := []string{}
	for _, v := range volumes {
		source := ""
		if v.VolumeSource.Secret != nil {
			source = v.VolumeSource.Secret.SecretName
		} else if v.VolumeSource.ConfigMap != nil {
			source = v.VolumeSource.ConfigMap.Name
		}
		sources = append(sources, source)
	}
	return sources
}

var _ = Describe("Apply Job Builder", func() {
	Context("Create Job with valid Config", func() {
		var (
			cfg = &JobConfig{
				Command:   "apply",
				Args:      []string{""},
				Namespace: "TestNS",
				Stack:     "TestStack",
				TfConfig:  "TestConfig",
				Tfvars:    "TestVars",
				Tfstate:   "TestState",
			}
			applyJob  *batchv1.Job
			spec      corev1.PodSpec
			err       error
			container corev1.Container
		)

		BeforeEach(func() {
			applyJob, err = BuildJob(cfg)
			spec = applyJob.Spec.Template.Spec

            // It should have not failed
			Expect(err).ShouldNot(HaveOccurred())
			Expect(applyJob).ToNot(BeNil())

            // It should have one container
		    Expect(spec.Containers).To(HaveLen(1))

			container = spec.Containers[0]
		})

		It("Should have the apply command set", func() {
			Expect(container.Command).Should(ContainElements(cfg.Command))
		})

		It("Should have the arguments set", func() {
			Expect(container.Args).Should(ContainElements(cfg.Args))
		})

		It("Should have volume mounts with secrets and configmap", func() {
			// check secreats and ConfigMaps are mounted in container
			sourceNames := []string{cfg.TfConfig, cfg.Tfvars, cfg.Tfstate}
			Expect(getVolumeSources(spec.Volumes)).To(ContainElements(sourceNames))
		})
	})
})
