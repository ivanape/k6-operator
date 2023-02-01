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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// K6SuiteSpec defines the desired state of K6Suite
type K6SuiteSpec struct {
	K6TestCases []K6Spec `json:"k6TestCases,omitempty"`
}

// K6SuiteStatus defines the observed state of K6Suite
type K6SuiteStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Stage            Stage              `json:"stage,omitempty"`
	K6TestConditions []K6SuiteCondition `json:"k6TestConditions,omitempty"`
}

// K6SuiteCondition contains condition information for an Issuer.
type K6SuiteCondition struct {
	// Type of the condition, known values are ('Ready').
	Stage Stage `json:"stage,omitempty"`

	// LastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	// +optional
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`

	// Reason is a brief machine readable explanation for the condition's last
	// transition.
	// +optional
	Reason string `json:"reason,omitempty"`

	// Message is a human readable description of the details of the last
	// transition, complementing reason.
	// +optional
	Message string `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// K6Suite is the Schema for the k6suites API
type K6Suite struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   K6SuiteSpec   `json:"spec,omitempty"`
	Status K6SuiteStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// K6SuiteList contains a list of K6Suite
type K6SuiteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []K6Suite `json:"items"`
}

func init() {
	SchemeBuilder.Register(&K6Suite{}, &K6SuiteList{})
}
