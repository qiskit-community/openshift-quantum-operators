/*
Copyright 2021.

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
	apps "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// QiskitPlaygroundSpec defines the desired state of QiskitPlayground
type QiskitPlaygroundSpec struct {
	// +kubebuilder:default:="jupyter/scipy-notebook:latest"
	Image string `json:"image,omitempty"`
	// +kubebuilder:default:="IfNotPresent"
	ImagePullPolicy apiv1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// +optional
	PVC string `json:"pvc,omitempty"`
	// +kubebuilder:default:=false
	LoadBalancer bool `json:"loadbalancer,omitempty" description:"Define if load balancer service type is supported. By default false"`
	// +optional
	Resources *apiv1.ResourceRequirements `json:"resources,omitempty"`
}

// QiskitPlaygroundStatus defines the observed state of QiskitPlayground
type QiskitPlaygroundStatus struct {
	// +optional
	Condition *apps.DeploymentCondition `json:"condition,omitempty"`
	//    Reason             *string          `json:"reason,omitempty" description:"one-word CamelCase reason for the condition's last transition"`
	// +optional
	//    Message            *string          `json:"message,omitempty" description:"human-readable message indicating details about last transition"`
	// +optional
	//    LastHeartbeatTime  *metav1.Time     `json:"lastHeartbeatTime,omitempty" description:"last time we got an update on a given condition"`
	// +optional
	//    LastTransitionTime *metav1.Time     `json:"lastTransitionTime,omitempty" description:"last time the condition transit from one status to another"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// QiskitPlayground is the Schema for the qiskitplaygrounds API
type QiskitPlayground struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QiskitPlaygroundSpec   `json:"spec,omitempty"`
	Status QiskitPlaygroundStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// QiskitPlaygroundList contains a list of QiskitPlayground
type QiskitPlaygroundList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QiskitPlayground `json:"items"`
}

func init() {
	SchemeBuilder.Register(&QiskitPlayground{}, &QiskitPlaygroundList{})
}
