package jobs

import (
	"math/rand"
    "strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// CharSet defines the alphanumeric set for random string generation
	charSet = "0123456789abcdefghijklmnopqrstuvwxyz"

	// the image used in the Job for performing the copies
	jobImage = "busybox:latest"

	// command to execute in the job
	jobCommand = "tfoctl"

	// name of the tf config volume in the job spec
	tfconfigVolName = "tfconf"

	// path for mounting the tf conf
	tfconfigPath = "/var/lib/tfoperator/tfconfig"

	// name of the tfvars volume in the job spec
	tfvarsVolName = "tfvars"

	// path for mounting the tfvars
	tfvarsPath = "/var/lib/tfoperator"

	// name of the tfvol in the job spec
	tfstateVolName = "tfstate"

	// path for mounting the tfstate
	tfstatePath = "/var/lib/tfoperator"
)

var (
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

	// Template for the Job
	jobTemplate = batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "unset", // this will be set
			Labels: map[string]string{},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes:       []corev1.Volume{},
					Containers: []corev1.Container{
						{
							Name:            "unset", // this will be set
							Image:           jobImage,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command:         []string{},
							Args:            []string{},
							VolumeMounts:    []corev1.VolumeMount{},
						},
					},
				},
			},
		},
	}
)

type JobConfig struct {
	Command   string   // command to execute
	Args      []string // options to the command
	Namespace string   // Stack's namespace
	Stack     string   // Stack name
	TfConfig  string   // TfConfig ConfigMap name
	Tfvars    string   // tfvars Secret name
	Tfstate   string   // tfstate Secret name
}

// buildJob returns a Job for running a command
// and mounting volumes from secrets and config maps
func BuildJob(cfg *JobConfig) (*batchv1.Job, error) {
	job := jobTemplate.DeepCopy()

	job.Name = strings.ToLower(cfg.Stack) +
               "-" + strings.ToLower(cfg.Command) +
               "-" + randomString(6)
	job.Namespace = cfg.Namespace

	jobPodSpec := &job.Spec.Template.Spec
	jobCont0 := &jobPodSpec.Containers[0]

	jobCont0.Name = "run-" + jobCommand
	jobCont0.Command = []string{jobCommand, cfg.Command}
	jobCont0.Args = append(cfg.Args, "--stack", cfg.Stack, "--namespace", cfg.Namespace)

	labels := map[string]string{
		"stack.tf-operator.io": cfg.Stack,
	}
	for k, v := range labels {
		job.ObjectMeta.Labels[k] = v
	}

	err := volumeFromSecret(jobPodSpec, tfvarsVolName, tfvarsPath, cfg.Tfvars)
	if err != nil {
	}

	err = volumeFromConfigMap(jobPodSpec, tfconfigVolName, tfconfigPath, cfg.TfConfig)
	if err != nil {
	}

    if cfg.Tfstate != "" {
        err = volumeFromSecret(jobPodSpec, tfstateVolName, tfstatePath, cfg.Tfstate)
	    if err != nil {
	    }
    }

	return job, nil
}

// volumeFromSecret mounts a volume from a secret in container 0 of a Job
func volumeFromSecret(podSpec *corev1.PodSpec, volName string, volPath string, secret string) error {
	secretVolume := corev1.Volume{
		Name: volName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secret,
			},
		},
	}
	podSpec.Volumes = append(podSpec.Volumes, secretVolume)

	secretVolumeMount := corev1.VolumeMount{
		Name:      volName,
		MountPath: volPath,
		ReadOnly:  true,
	}

	jobCont0 := &podSpec.Containers[0]
	jobCont0.VolumeMounts = append(jobCont0.VolumeMounts, secretVolumeMount)

	return nil
}

// volumeFromConfigmap mounts a volume from a configmap in container 0 of a Job
func volumeFromConfigMap(podSpec *corev1.PodSpec, volName string, volPath string, cfgMap string) error {
	configMapVolume := corev1.Volume{
		Name: volName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cfgMap,
				},
			},
		},
	}
	podSpec.Volumes = append(podSpec.Volumes, configMapVolume)

	configMapVolumeMount := corev1.VolumeMount{
		Name:      volName,
		MountPath: volPath,
		ReadOnly:  true,
	}

	jobCont0 := &podSpec.Containers[0]
	jobCont0.VolumeMounts = append(jobCont0.VolumeMounts, configMapVolumeMount)

	return nil
}

// RandomString returns a random alphanumeric string.
func randomString(n int) string {
	result := make([]byte, n)
	for i := range result {
		result[i] = charSet[rnd.Intn(len(charSet))]
	}
	return string(result)
}
