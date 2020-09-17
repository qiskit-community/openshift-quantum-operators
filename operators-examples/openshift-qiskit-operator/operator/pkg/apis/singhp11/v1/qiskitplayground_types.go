package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// QiskitPlaygroundSpec defines the desired state of QiskitPlayground
type QiskitPlaygroundSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// QiskitPlaygroundStatus defines the observed state of QiskitPlayground
type QiskitPlaygroundStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QiskitPlayground is the Schema for the qiskitplaygrounds API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=qiskitplaygrounds,scope=Namespaced
type QiskitPlayground struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QiskitPlaygroundSpec   `json:"spec,omitempty"`
	Status QiskitPlaygroundStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QiskitPlaygroundList contains a list of QiskitPlayground
type QiskitPlaygroundList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QiskitPlayground `json:"items"`
}

func init() {
	SchemeBuilder.Register(&QiskitPlayground{}, &QiskitPlaygroundList{})
}
