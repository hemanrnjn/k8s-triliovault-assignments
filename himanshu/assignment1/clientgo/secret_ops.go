package clientgo

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

func SecretOps() {
	clientset := getClientSet(false)
	secretClient := clientset.CoreV1().Secrets("himanshu")

	fileBytes, err := ioutil.ReadFile("../../himanshu/assignment1/sample/secret.yaml")
	if err != nil {
		panic(err.Error())
	}

	var secretSpec corev1.Secret
	dec := k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(string(fileBytes))), 1000)

	if err := dec.Decode(&secretSpec); err != nil {
		panic(err)
	}

	// Create Secret
	fmt.Println("Creating secret...")
	result, err := secretClient.Create(context.TODO(), &secretSpec, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created secret %q.\n", result.GetObjectMeta().GetName())

	// Get Secret
	fmt.Printf("Getting secret in namespace %q:\n", "himanshu")
	result, getErr := secretClient.Get(context.TODO(), "demo-secret", metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("Failed to get latest version of Secret: %v", getErr))
	}
	fmt.Printf("Latest Secret: %s \n", result.Name)

	// Updating Secret
	fmt.Printf("Updating secret in namespace %q:\n", "himanshu")
	result.Data["username"] = []byte("bXktc2VjcmV0LWFwcA==") // change username
	_, updateErr := secretClient.Update(context.TODO(), result, metav1.UpdateOptions{})
	if updateErr != nil {
		panic(updateErr)
	}
	fmt.Println("Updated secret...")

	// Verifying Update
	fmt.Println("Verifying Update...")
	result, getErr = secretClient.Get(context.TODO(), "demo-secret", metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("Failed to get latest version of Secret: %v", getErr))
	}
	if bytes.Compare(result.Data["username"], []byte("bXktc2VjcmV0LWFwcA==")) == 0 {
		fmt.Println("Verified Successfully")
	} else {
		panic(fmt.Errorf("Verification failed. Secret found: %b, expected %b", result.Data["username"],
			[]byte("bXktc2VjcmV0LWFwcA==")))
	}

	// Delete Secret
	fmt.Printf("Deleting secret in namespace %q:\n", "himanshu")
	deletePolicy := metav1.DeletePropagationForeground
	deleteErr := secretClient.Delete(context.TODO(), "demo-secret", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if deleteErr != nil {
		panic(deleteErr)
	}
	fmt.Println("Deleted secret...")
}
