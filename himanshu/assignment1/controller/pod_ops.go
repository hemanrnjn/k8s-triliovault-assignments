package controller

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	corev1 "k8s.io/api/core/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func PodOps() {
	cl, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		fmt.Println("failed to create client")
		os.Exit(1)
	}

	fileBytes, err := ioutil.ReadFile("../../himanshu/assignment1/sample/pod.yaml")
	if err != nil {
		fileBytes, err = ioutil.ReadFile("sample/pod.yaml")
		if err != nil {
			panic(err.Error())
		}
	}

	// Create Pod
	fmt.Println("Creating pod...")
	var podSpec corev1.Pod
	dec := k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(string(fileBytes))), 1000)

	if err := dec.Decode(&podSpec); err != nil {
		panic(err)
	}

	err = cl.Create(context.Background(), &podSpec)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created pod...")

	// Get Pod
	fmt.Println("Listing pods...")
	var pod corev1.Pod
	if err = cl.Get(context.Background(), client.ObjectKey{
		Namespace: "himanshu",
		Name:      "demo-pod",
	}, &pod); err != nil {
		panic(fmt.Errorf("Failed to get latest version of Pod: %v", err))
	}
	fmt.Printf("Latest Pod: %s \n", pod.Name)

	// Update Pod
	fmt.Println("Updating pod...")
	pod.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
	err = cl.Update(context.Background(), &pod)
	if err != nil {
		panic(err)
	}
	fmt.Println("Updated pod...")

	// Verifying Update
	fmt.Println("Verifying Update...")
	if err = cl.Get(context.Background(), client.ObjectKey{
		Namespace: "himanshu",
		Name:      "demo-pod",
	}, &pod); err != nil {
		panic(fmt.Errorf("Failed to get latest version of Pod: %v", err))
	}
	if pod.Spec.Containers[0].Image == "nginx:1.13" {
		fmt.Println("Verified Successfully")
	} else {
		panic(fmt.Errorf("Verification failed. Image found %s, expected: nginx:1.13",
			pod.Spec.Containers[0].Image))
	}

	// Delete Pod
	fmt.Println("Deleting pod...")
	err = cl.Delete(context.Background(), &pod)
	if err != nil {
		panic(err)
	}
	fmt.Println("Deleted pod...")
}
