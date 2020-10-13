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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:validation:Enum=Mongo;Cassandra;PSQL
// DBType defines the type db instance of an DBClusterSpec
type DBType string

const (
	// Mongo is mongo db type of DBCluster
	Mongo DBType = "Mongo"

	// Cassandra is cassandra db type of DBCluster
	Cassandra DBType = "Cassandra"

	// MySQL is MySQL db type of DBCluster
	PSQL DBType = "PSQL"
)

// +kubebuilder:validation:type=string
// Status specifies the status of WorkloadJob operating on
type Status string

const (
	// InProgress means the process is under execution
	InProgress Status = "InProgress"

	// Failed means the process is unsuccessful due to an error.
	Failed Status = "Failed"

	// Available means the resources blocked for the process execution are now available.
	Active Status = "Active"
)

// DBClusterSpec defines the desired state of DBCluster
type DBClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Enum=Mongo;Cassandra;PSQL
	// Type is the type of database for DBCluster setup.
	Type DBType `json:"type"`

	// +kubebuilder:validation:Optional
	Replicas int64 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Optional
	TerminationGracePeriod int64 `json:"terminationGracePeriod,omitempty"`
}

// DBClusterStatus defines the observed state of DBCluster
type DBClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=InProgress;Active;Failed
	// +nullable:true
	// Status is the status of the cluster creation operation.
	Status Status `json:"status,omitempty"`
}

// DBCluster is the Schema for the dbclusters API
// +kubebuilder:object:root=true
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type DBCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DBClusterSpec   `json:"spec,omitempty"`
	Status DBClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DBClusterList contains a list of DBCluster
type DBClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DBCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DBCluster{}, &DBClusterList{})
}
