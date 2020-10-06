package clientgo

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func int32Ptr(i int32) *int32 {
	return &i
}

func DeployOps(clientset kubernetes.Clientset) {
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	bytes, err := ioutil.ReadFile("../sample/deployment.yaml")
	if err != nil {
		panic(err.Error())
	}

	var deploymentSpec appsv1.Deployment
	err = yaml.Unmarshal(bytes, &deploymentSpec)
	if err != nil {
		panic(err.Error())
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(&deploymentSpec)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	// List Deployments
	fmt.Printf("Listing deployment in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}

	// Updating Deployments
	fmt.Printf("Updating deployment in namespace %q:\n", apiv1.NamespaceDefault)
	result, getErr := deploymentsClient.Get("demo-deployment", metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
	}

	result.Spec.Replicas = int32Ptr(1)                           // reduce replica count
	result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
	_, updateErr := deploymentsClient.Update(result)
	if updateErr != nil {
		panic(err)
	}
	fmt.Println("Updated deployment...")

	// Delete Deployments
	fmt.Printf("Deleting deployment in namespace %q:\n", apiv1.NamespaceDefault)
	deletePolicy := metav1.DeletePropagationForeground
	deleteErr := deploymentsClient.Delete("demo-deploy", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if deleteErr != nil {
		panic(deleteErr)
	}
	fmt.Println("Deleted deployment...")
}
