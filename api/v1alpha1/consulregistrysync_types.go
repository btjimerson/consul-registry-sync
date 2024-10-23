/*
Copyright 2024.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ConsulRegistrySyncSpec defines the desired state of ConsulRegistrySync
type ConsulRegistrySyncSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The address of the Consul service; i.e. http://localhost:8500
	ConsulAddress string `json:"consuladdress,omitempty"`

	// The path to the kube config
	KubeConfig string `json:"kubeconfig,omitempty"`

	// The namespace to create service entries in
	ServiceEntryNamespace string `json:"serviceentrynamespace,omitempty"`
}

// ConsulRegistrySyncStatus defines the observed state of ConsulRegistrySync
type ConsulRegistrySyncStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ConsulRegistrySync is the Schema for the consulregistrysyncs API
type ConsulRegistrySync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConsulRegistrySyncSpec   `json:"spec,omitempty"`
	Status ConsulRegistrySyncStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConsulRegistrySyncList contains a list of ConsulRegistrySync
type ConsulRegistrySyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConsulRegistrySync `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConsulRegistrySync{}, &ConsulRegistrySyncList{})
}
