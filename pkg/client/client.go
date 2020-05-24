package client

import (
    "context"
    "fmt"
    "io/ioutil"
    "os"
    "path"

	tfo "github.com/pablochacin/tf-operator/api/v1alpha1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    apierr "k8s.io/apimachinery/pkg/api/errors"
    apirtm "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/client-go/tools/clientcmd"
    ctlclient "sigs.k8s.io/controller-runtime/pkg/client"
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
    rc ctlclient.Client
}

// GetStack returns an existing stack or an error
func (c *client)GetStack(stackName string, namespace string) (*tfo.Stack, error) {
	return nil, nil
}

// NewClientFromKubeconfig creates a Client from a kubeconfig
func NewFromKubeconfig(kubeconfig string) (Client, error) {
	sch := apirtm.NewScheme()
	corev1.AddToScheme(sch)
	tfo.AddToScheme(sch)
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	kubeClient, err := ctlclient.New(
		cfg,
		ctlclient.Options{Scheme: sch},
	)
	if err != nil {
		return nil, err
	}

	return &client{rc: kubeClient}, nil
}

// NewClientFromRuntimeClient create a Client from a runtime client
func NewFromRuntimeClient(rc ctlclient.Client) (Client, error) {
    return &client{rc: rc}, nil
}

// CreateStack creates a stack from local tf files
func (c *client)CreateStack(name string, namespace string, tfconf string, tfvars string) (*tfo.Stack, error){

    // check stack doesn't exits
    err :=  c.rc.Get(
        context.TODO(),
        ctlclient.ObjectKey{Name: name, Namespace: namespace},
        &tfo.Stack{},
    )
    if err == nil {
        errDesc := fmt.Sprintf("stack %s already exists in namespace %s", name, namespace)
        return nil, NewTFOError(errDesc, ErrorReasonAlreadyExists)
    }

    // ignore not found error, as it is expected
    if !apierr.IsNotFound(err) {
        errDesc := fmt.Sprintf("runtime error creating stack: %s", err)
        return nil, NewTFOError(errDesc, ErrorReasonRuntimeError)
    }

    // create configmap for config
    tfconfMap, err := createConfigMap(name+"-tfconf",namespace, tfconf)
    if err != nil {
        return  nil, err
    }
    err = c.rc.Create(context.TODO(), tfconfMap)
    if err != nil {
        return  nil, err
    }

    // create secret for tfvars
    tfvarsSecret, err := createSecret(name+"-tfvars", namespace, tfvars)
    if err != nil {
        return  nil, err
    }
    err = c.rc.Create(context.TODO(), tfvarsSecret)
    if err != nil {
        return  nil, err
    }


    //create Stack
    stack := &tfo.Stack{
        ObjectMeta: metav1.ObjectMeta{
            Name: name,
            Namespace: namespace,
        },
        Spec: tfo.StackSpec{
            TfConfig: corev1.LocalObjectReference{Name: tfconfMap.Name},
            TfVars:   corev1.LocalObjectReference{Name: tfvarsSecret.Name},
        },
    }

    err = c.rc.Create(context.TODO(), stack)
    if err != nil {
        return  nil, err
    }

    return stack, nil

}

// createConfigMap create a ConfigMap from a the files in a directory
func createConfigMap(name string, namespace string, dirPath string) (*corev1.ConfigMap, error) {

    configMap := &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name: name,
            Namespace: namespace,
        },
        Data: map[string]string{},
    }

    fileList, err := ioutil.ReadDir(dirPath)
    if err != nil {
        desc := fmt.Sprintf("error accessing directory %s: %v", dirPath, err)
        return nil, NewTFOError(desc, ErrorReasonFileCanNotBeAccessed)
    }

    if len(fileList) == 0 {
        desc := fmt.Sprintf("directory %s is empty", dirPath)
        return nil, NewTFOError(desc, ErrorReasonFileCanNotBeAccessed)
    }

    for _, file := range fileList {
        filePath := path.Join(dirPath, file.Name())
        if file.Size() == 0 {
            desc := fmt.Sprintf("file %s is empty", filePath)
            return nil, NewTFOError(desc, ErrorReasonInvalidFileContent)
        }
        data, err := ioutil.ReadFile(filePath)
        if err != nil {
            desc := fmt.Sprintf("error reading file %s: %v", filePath, err)
            return nil, NewTFOError(desc, ErrorReasonFileCanNotBeAccessed)
        }
        configMap.Data[file.Name()] = string(data)

    }

    return configMap, nil
}

// createSecret create a Secret from a file
func createSecret(name string, namespace string, filePath string) (*corev1.Secret, error) {
    secret := &corev1.Secret{
        ObjectMeta: metav1.ObjectMeta{
            Name: name,
            Namespace: namespace,
        },
        Data: map[string][]byte{},
    }

    fileInfo, err := os.Stat(filePath)
    if err != nil {
        desc := fmt.Sprintf("error accessing file %s: %v", filePath, err)
        return nil, NewTFOError(desc, ErrorReasonFileCanNotBeAccessed)
    }

    if fileInfo.Size() == 0 {
        desc := fmt.Sprintf("file %s is empty", filePath)
        return nil, NewTFOError(desc, ErrorReasonInvalidFileContent)
    }

    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        desc := fmt.Sprintf("error reading file %s: %v", filePath, err)
        return nil, NewTFOError(desc, ErrorReasonFileCanNotBeAccessed)
    }
    secret.Data[fileInfo.Name()] = data

    return secret, nil
}
