package controller

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func int32Ptr(i int32) *int32 {
	return &i
}

func DeployOps() {
	cl, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		fmt.Println("failed to create client")
		os.Exit(1)
	}

	fileBytes, err := ioutil.ReadFile("../../himanshu/assignment1/sample/deployment.yaml")
	if err != nil {
		fileBytes, err = ioutil.ReadFile("sample/deployment.yaml")
		if err != nil {
			panic(err.Error())
		}
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	var deploymentSpec appsv1.Deployment
	dec := k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(string(fileBytes))), 1000)
	if err := dec.Decode(&deploymentSpec); err != nil {
		panic(err)
	}

	err = cl.Create(context.Background(), &deploymentSpec)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created deployment...")

	// Get Deployment
	fmt.Println("Getting deployment...")
	var deploy appsv1.Deployment
	if err = cl.Get(context.Background(), client.ObjectKey{
		Namespace: "himanshu",
		Name:      "demo-deploy",
	}, &deploy); err != nil {
		cleanUpControllerResources("deploy", deploymentSpec)
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", err))
	}
	fmt.Printf("Latest Deployment: %s with (%d replicas)\n", deploy.Name, *deploy.Spec.Replicas)

	// Update Deployment
	fmt.Println("Updating deployment...")
	deploy.Spec.Replicas = int32Ptr(2)
	deploy.Spec.Template.Spec.Containers[0].Image = "nginx:1.13"
	err = cl.Update(context.Background(), &deploy)
	if err != nil {
		cleanUpControllerResources("deploy", deploy)
		panic(err)
	}
	fmt.Println("Updated deployment...")

	// Verifying Update
	fmt.Println("Verifying Update...")
	if err = cl.Get(context.Background(), client.ObjectKey{
		Namespace: "himanshu",
		Name:      "demo-deploy",
	}, &deploy); err != nil {
		cleanUpControllerResources("deploy", deploy)
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", err))
	}
	if *deploy.Spec.Replicas == 2 && deploy.Spec.Template.Spec.Containers[0].Image == "nginx:1.13" {
		fmt.Println("Verified Successfully")
	} else {
		cleanUpControllerResources("deploy", deploy)
		panic(fmt.Errorf("Verification failed. Replicas found: %d, expected 2 and Image found %s, expected: nginx:1.13",
			*deploy.Spec.Replicas, deploy.Spec.Template.Spec.Containers[0].Image))
	}

	// Delete Deployment
	fmt.Println("Deleting deployment...")
	err = cl.Delete(context.Background(), &deploy)
	if err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment...")
}
