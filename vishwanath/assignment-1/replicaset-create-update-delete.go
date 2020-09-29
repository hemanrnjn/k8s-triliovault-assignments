//Performing CRUD Operations on a replica-set via client-go
//Program is run using : "./app --kubeconfig /root/.kube/config"
//While running the program, try to watch the resources on another terminal
package main

import (
	"context"
	"flag"
	"fmt"
	"k8s.io/client-go/util/retry"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func main() {

	//Fetching the kubeconfig file
	var kubernetesConfig *string
	if homeDir := homeDir(); homeDir != "" {
		kubernetesConfig = flag.String("kubeconfig", filepath.Join(homeDir, ".kube", "config"), "/root/.kube/config")
	} else {
		kubernetesConfig = flag.String("kubeconfig", "", "/root/.kube/config")
	}
	flag.Parse()

	//Building config
	config, err := clientcmd.BuildConfigFromFlags("", *kubernetesConfig)
	if err != nil {
		panic(err.Error())
	}

	//Building client from config
	newclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//Setting namespace for usage and starting interaction with Replica-Set
	const NamespaceName string = "vishu"
	replicaSetClient := newclient.AppsV1().ReplicaSets(NamespaceName)

	//Providing Template Spec of replica-set
	replicaSet := &appsv1.ReplicaSet{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "vishu-replicaset",
			Namespace: "vishu",
		},
		Spec: appsv1.ReplicaSetSpec{
			Replicas: int32Ptr(3),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "vishu-replicaset",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "vishu-replicaset",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "replica-set-nginx",
							Image: "nginx:1.18",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	//Creating the replica-set
	fmt.Println("Creating replica-set-nginx ...")
	rsResult, err := replicaSetClient.Create(context.TODO(), replicaSet, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created replica-set %q.\n", rsResult.GetObjectMeta().GetName())

	//Listing the replica-set
	fmt.Println("Showing Replica-Set: ")
	replicaList, err := replicaSetClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, replicas := range replicaList.Items {
		fmt.Println("This is replica-set: %s and these are the number of Pods: %d", replicas.Name, *replicas.Spec.Replicas) //de-referencing is necessary because the original data-type is pointer only
	}

	//Updating the replica-set
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := replicaSetClient.Get(context.TODO(), "vishu-replicaset", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Replica-Set: %v", getErr))
		}

		result.Spec.Replicas = int32Ptr(5)                           // increasing replicas
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.19" // updating nginx
		_, updateErr := replicaSetClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Replica-Set Successfully Updated")

	//Deleting Replica-Set
	fmt.Println("Deleting Replica-Set ... ")
	deleteReplica := metav1.DeletePropagationForeground
	if deleteErr := replicaSetClient.Delete(context.TODO(), "vishu-replicaset", metav1.DeleteOptions{
		PropagationPolicy: &deleteReplica,
	}); deleteErr != nil {
		panic(err.Error())
	}
	fmt.Println("ReplicaSet deleted")
}

func homeDir() string {
	return os.Getenv("ROOTPATH")
}
func int32Ptr(i int32) *int32 { return &i }
