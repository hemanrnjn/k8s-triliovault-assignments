package clientgo

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func SecretOps(clientset kubernetes.Clientset) {

	secretClient := clientset.CoreV1().Secrets("himanshu")

	bytes, err := ioutil.ReadFile("../sample/secret.yaml")
	if err != nil {
		panic(err.Error())
	}

	var secretSpec apiv1.Secret
	err = yaml.Unmarshal(bytes, &secretSpec)
	if err != nil {
		panic(err.Error())
	}

	// Create Secret
	fmt.Println("Creating secret...")
	result, err := secretClient.Create(&secretSpec)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created secret %q.\n", result.GetObjectMeta().GetName())

	// List Secret
	fmt.Printf("Listing secrets in namespace %q:\n", "himanshu")
	list, err := secretClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s \n", d.Name)
	}

	// Updating Secret
	fmt.Printf("Updating secret in namespace %q:\n", "himanshu")
	result, getErr := secretClient.Get("demo-secret", metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("Failed to get latest version of Secret: %v", getErr))
	}

	result.Data["username"] = []byte("bXktc2VjcmV0LWFwcA==") // change username
	_, updateErr := secretClient.Update(result)
	if updateErr != nil {
		panic(err)
	}
	fmt.Println("Updated secret...")

	// Delete Secret
	fmt.Printf("Deleting secret in namespace %q:\n", "himanshu")
	deletePolicy := metav1.DeletePropagationForeground
	deleteErr := secretClient.Delete("demo-secret", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if deleteErr != nil {
		panic(deleteErr)
	}
	fmt.Println("Deleted secret...")
}
	