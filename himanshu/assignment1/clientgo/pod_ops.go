package clientgo

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func PodOps(clientset kubernetes.Clientset) {

	podsClient := clientset.CoreV1().Pods(apiv1.NamespaceDefault)

	bytes, err := ioutil.ReadFile("../sample/pod.yaml")
	if err != nil {
		panic(err.Error())
	}

	var podSpec apiv1.Pod
	err = yaml.Unmarshal(bytes, &podSpec)
	if err != nil {
		panic(err.Error())
	}

	// Create Pod
	fmt.Println("Creating pod...")
	result, err := podsClient.Create(&podSpec)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created pod %q.\n", result.GetObjectMeta().GetName())

	// List Pods
	fmt.Printf("Listing pods in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := podsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s \n", d.Name)
	}

	// Updating Pod
	fmt.Printf("Updating pod in namespace %q:\n", apiv1.NamespaceDefault)
	result, getErr := podsClient.Get("demo-pod", metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("Failed to get latest version of Pod: %v", getErr))
	}

	result.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
	_, updateErr := podsClient.Update(result)
	if updateErr != nil {
		panic(err)
	}
	fmt.Println("Updated pod...")

	// Delete Pods
	fmt.Printf("Deleting pod in namespace %q:\n", apiv1.NamespaceDefault)
	deletePolicy := metav1.DeletePropagationForeground
	deleteErr := podsClient.Delete("demo-pod", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if deleteErr != nil {
		panic(deleteErr)
	}
	fmt.Println("Deleted pod...")
}
