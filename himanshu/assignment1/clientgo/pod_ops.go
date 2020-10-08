package clientgo

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

func PodOps() {
	clientset := getClientSet(true)
	podsClient := clientset.CoreV1().Pods("himanshu")

	fileBytes, err := ioutil.ReadFile("../../himanshu/assignment1/sample/pod.yaml")
	if err != nil {
		fileBytes, err = ioutil.ReadFile("sample/pod.yaml")
		if err != nil {
			panic(err.Error())
		}
	}

	var podSpec apiv1.Pod
	dec := k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(string(fileBytes))), 1000)

	if err := dec.Decode(&podSpec); err != nil {
		panic(err)
	}

	// Create Pod
	fmt.Println("Creating pod...")
	result, err := podsClient.Create(context.TODO(), &podSpec, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created pod %q.\n", result.GetObjectMeta().GetName())

	// Get Pod
	fmt.Printf("Getting pod in namespace %q:\n", "himanshu")
	result, getErr := podsClient.Get(context.TODO(), "demo-pod", metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("Failed to get latest version of Pod: %v", getErr))
	}
	fmt.Printf("Latest Pod: %s \n", result.Name)

	// Updating Pod
	fmt.Printf("Updating pod in namespace %q:\n", "himanshu")
	result.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
	_, updateErr := podsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
	if updateErr != nil {
		panic(updateErr)
	}
	fmt.Println("Updated pod...")

	// Verifying Update
	fmt.Println("Verifying Update...")
	result, getErr = podsClient.Get(context.TODO(), "demo-pod", metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("Failed to get latest version of Pod: %v", getErr))
	}
	if result.Spec.Containers[0].Image == "nginx:1.13" {
		fmt.Println("Verified Successfully")
	} else {
		panic(fmt.Errorf("Verification failed. Image found %s, expected: nginx:1.13",
			result.Spec.Containers[0].Image))
	}

	// Delete Pods
	fmt.Printf("Deleting pod in namespace %q:\n", "himanshu")
	deletePolicy := metav1.DeletePropagationForeground
	deleteErr := podsClient.Delete(context.TODO(), "demo-pod", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if deleteErr != nil {
		panic(deleteErr)
	}
	fmt.Println("Deleted pod...")
}
