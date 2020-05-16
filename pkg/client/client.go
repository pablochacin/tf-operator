package client

import (
    "context"
    "fmt"
    "io/ioutil"
    "os"

	tfo "github.com/pablochacin/tf-operator/api/v1alpha1"
    corev1 "k8s.io/api/core/v1"
    apierr "k8s.io/apimachinery/pkg/api/errors"
    runtime "sigs.k8s.io/controller-runtime/pkg/client"
)

// Client exposes a high-level interface for tf-operator actions
type Client interface {
	// GetStack returns a stack with the name in the given namespace
	GetStack(stackName string, namespace string) (*tfo.Stack, error)

    // CreateStack creates a stack from local tf files
    CreateStack(name string, namespace string, tfconf string, tfvars string) (*tfo.Stack, error)
}

// tfoClient Client implementation
type client struct {
    rc runtime.Client
}

// GetStack returns an existing stack or an error
func (c *client)GetStack(stackName string, namespace string) (*tfo.Stack, error) {
	return nil, nil
}

// NewClientFromKubeconfig creates a Client from a kubeconfig
func NewFromKubeconfig(kubeconfig string) (Client, error) {
	return &client{}, nil
}

// NewClientFromRuntimeClient create a Client from a runtime client
func NewFromRuntimeClient(rc runtime.Client) (Client, error) {
    return &client{rc: rc}, nil
}

// CreateStack creates a stack from local tf files
func (c *client)CreateStack(name string, namespace string, tfconf string, tfvars string) (*tfo.Stack, error){

    // check stack doesn't exits
    err :=  c.rc.Get(
        context.TODO(),
        runtime.ObjectKey{Name: name, Namespace: namespace},
        &tfo.Stack{},
    )
    if !apierr.IsNotFound(err) {
        errDesc := fmt.Sprintf("stack %s already exists in namespace %s", name, namespace)
        return nil, NewTFOError(errDesc, ErrorReasonAlreadyExists)
    }

    // create configmap for config
    _, err = createConfigMap(name+"-"+tfconf,namespace, tfconf)
    if err != nil {
        return  nil, err
    }

    // create secret for tfvars
    _, err = createSecret(name+"-tfvars", namespace, tfvars)
    if err != nil {
        return  nil, err
    }

    return nil, nil
}

// createConfigMap create a ConfigMap from a the files in a directory
func createConfigMap(name string, namespace string, dirPath string) (*corev1.ConfigMap, error) {
    fileList, err := ioutil.ReadDir(dirPath)
    if err != nil {
        desc := fmt.Sprintf("error accessing directory %s: %v", dirPath, err)
        return nil, NewTFOError(desc, ErrorReasonFileCanNotBeAccessed)
    }

    if len(fileList) == 0 {
        desc := fmt.Sprintf("directory %s is empty", dirPath)
        return nil, NewTFOError(desc, ErrorReasonFileCanNotBeAccessed)
    }

    return nil, nil
}

// createSecret create a Secret from a file
func createSecret(name string, namespace string, filePath string) (*corev1.Secret, error) {
    _, err := os.Stat(filePath)
    if err != nil {
        desc := fmt.Sprintf("error accessing file %s: %v", filePath, err)
        return nil, NewTFOError(desc, ErrorReasonFileCanNotBeAccessed)
    }

    return nil, nil
}
