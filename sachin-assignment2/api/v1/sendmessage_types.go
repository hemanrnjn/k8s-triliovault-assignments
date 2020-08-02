/*
Copyright 2020 Sachin.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SendMessageSpec defines the desired state of SendMessage
type SendMessageSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Message where you need to add your message. Edit SendMessage_types.go to remove/update
	// +kubebuilder:validation:Optional
	Message string `json:"message,omitempty"`

	// MessageCarrier where you will define the MessageCarrier like WhatsApp or Telegram. Edit SendMessage_types.go to remove/update
	// +kubebuilder:validation:Optional
	MessageCarrier string `json:"messageCarrier,omitempty"`

	// CarrierToken where you will give the access token for calling api. Edit SendMessage_types.go to remove/update
	// +kubebuilder:validation:Optional
	CarrierToken string `json:"carrierToken,omitempty"`
}

// SendMessageStatus defines the observed state of SendMessage
type SendMessageStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Status where status will be stored of the message. Edit SendMessage_types.go to remove/update
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum:Sent;Finished
	Status string `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SendMessage is the Schema for the sendmessages API
// +kubebuilder:subresource:status
type SendMessage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SendMessageSpec   `json:"spec,omitempty"`
	Status SendMessageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SendMessageList contains a list of SendMessage
type SendMessageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SendMessage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SendMessage{}, &SendMessageList{})
}
