package client

import (
	tfo "github.com/pablochacin/tf-operator/api/v1alpha1"
)

// Client exposes a high-level interface for tf-operator actions
type Client interface {

	// GetStack returns a stack with the name in the given namespace
	GetStack(stackName string, namespace string) (*tfo.Stack, error)
}

// client implementation
type client struct {
}

// GetStack returns an existing stack or an error
func (c client) GetStack(stackName string, namespace string) (*tfo.Stack, error) {
	return nil, nil
}

// NewClient creates a client from a kubeconfig
func NewClient(kubeconfig string) (Client, error) {
	return client{}, nil
}
