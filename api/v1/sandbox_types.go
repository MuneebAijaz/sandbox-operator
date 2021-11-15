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

/*
To provision a SandBox env on a given cluster on onboarding of a new User.
By default for every user, one Sandbox env would be created, but users can request more(with some MAX limit) SandBox env.
On deletion of a particular User, all the Sandbox env owned by the user should be deleted.

Specifications

Kind: User
Name: string (Required=True)
OrgName: string (Required=True)
DeptName: string (Required=False, ENUM={Engineering, QA, Operations})
SandBoxCount: int (Required=False, Default=1, Limit=5)


Kind: SandBox
Name: string(Required=True) # Format: SB-{UserName}-{Count}
Type: string (Required=False, ENUM={T1} // We can add more types in future

*/

// +kubebuilder:validation:Required

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

/*
To provision a SandBox env on a given cluster on onboarding of a new User.
By default for every user, one Sandbox env would be created, but users can request more(with some MAX limit) SandBox env.
On deletion of a particular User, all the Sandbox env owned by the user should be deleted.

Specifications

Kind: User
Name: string (Required=True)
OrgName: string (Required=True)
DeptName: string (Required=False, ENUM={Engineering, QA, Operations})
SandBoxCount: int (Required=False, Default=1, Limit=5)
Kind: SandBox
Name: string(Required=True) # Format: SB-{UserName}-{Count}
Type: string (Required=False, ENUM={T1} // We can add more types in future

*/

// SandboxSpec defines the desired state of Sandbox
type SandboxSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Sandbox. Edit sandbox_types.go to remove/update
	Name string `json:"name"`

	// +optional
	// +kubebuilder:validation:Enum=T1
	Type string `json:"type"`
}

// SandboxStatus defines the observed state of Sandbox
type SandboxStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Name string `json:"name"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Sandbox is the Schema for the sandboxes API
type Sandbox struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SandboxSpec   `json:"spec,omitempty"`
	Status SandboxStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SandboxList contains a list of Sandbox
type SandboxList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Sandbox `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Sandbox{}, &SandboxList{})
}
