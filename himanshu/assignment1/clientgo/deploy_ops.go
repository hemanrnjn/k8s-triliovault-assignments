package clientgo

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

func int32Ptr(i int32) *int32 {
	return &i
}

func DeployOps() {
	clientset := getClientSet(true)
	deploymentsClient := clientset.AppsV1().Deployments("himanshu")

	fileBytes, err := ioutil.ReadFile("../../himanshu/assignment1/sample/deployment.yaml")
	if err != nil {
		fileBytes, err = ioutil.ReadFile("sample/deployment.yaml")
		if err != nil {
			panic(err.Error())
		}
	}

	var deploymentSpec appsv1.Deployment
	dec := k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(string(fileBytes))), 1000)

	if err := dec.Decode(&deploymentSpec); err != nil {
		panic(err)
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), &deploymentSpec, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	// Get Deployment
	fmt.Printf("Getting deployment in namespace %q:\n", "himanshu")
	result, getErr := deploymentsClient.Get(context.TODO(), "demo-deploy", metav1.GetOptions{})
	if getErr != nil {
		cleanUpClientGoResources("deploy")
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
	}
	fmt.Printf("Latest Deployment: %s with (%d replicas)\n", result.Name, *result.Spec.Replicas)

	// Updating Deployments
	result.Spec.Replicas = int32Ptr(2)                           // reduce replica count
	result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
	_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
	if updateErr != nil {
		cleanUpClientGoResources("deploy")
		panic(updateErr)
	}
	fmt.Println("Updated deployment...")

	// Verifying Update
	fmt.Println("Verifying Update...")
	result, getErr = deploymentsClient.Get(context.TODO(), "demo-deploy", metav1.GetOptions{})
	if getErr != nil {
		cleanUpClientGoResources("deploy")
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
	}
	if *result.Spec.Replicas == 2 && result.Spec.Template.Spec.Containers[0].Image == "nginx:1.13" {
		fmt.Println("Verified Successfully")
	} else {
		cleanUpClientGoResources("deploy")
		panic(fmt.Errorf("Verification failed. Replicas found: %d, expected 2 and Image found %s, expected: nginx:1.13",
			*result.Spec.Replicas, result.Spec.Template.Spec.Containers[0].Image))
	}

	// Delete Deployments
	fmt.Printf("Deleting deployment in namespace %q:\n", "himanshu")
	deletePolicy := metav1.DeletePropagationForeground
	deleteErr := deploymentsClient.Delete(context.TODO(), "demo-deploy", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if deleteErr != nil {
		panic(deleteErr)
	}
	fmt.Println("Deleted deployment...")
}
