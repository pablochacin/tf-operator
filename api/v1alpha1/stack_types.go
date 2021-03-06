/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StackSpec defines the desired state of Stack
type StackSpec struct {
	// Reference to the config map with the configuration file(s)
	TfConfig corev1.LocalObjectReference `json:"tfconfig"`

	// Reference to secrect with tfvars
	TfVars corev1.LocalObjectReference `json:"tfvars"`
}

// StackStatus defines the observed state of Stack
type StackStatus struct {

	// Reference to secrect with tf state
	TfState corev1.LocalObjectReference `json:"tfstate"`

	// base64 encoded tfout
	TfOutput string `json:"tfout"`
}

// +kubebuilder:object:root=true

// Stack is the Schema for the stacks API
type Stack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StackSpec   `json:"spec,omitempty"`
	Status StackStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StackList contains a list of Stack
type StackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Stack `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Stack{}, &StackList{})
}
