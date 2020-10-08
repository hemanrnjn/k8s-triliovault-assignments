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

func SecretOps() {
	cl, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		fmt.Println("failed to create client")
		os.Exit(1)
	}

	fileBytes, err := ioutil.ReadFile("../../himanshu/assignment1/sample/secret.yaml")
	if err != nil {
		fileBytes, err = ioutil.ReadFile("sample/secret.yaml")
		if err != nil {
			panic(err.Error())
		}
	}

	// Create Secret
	fmt.Println("Creating secret...")
	var secretSpec corev1.Secret
	dec := k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(string(fileBytes))), 1000)

	if err := dec.Decode(&secretSpec); err != nil {
		panic(err)
	}

	err = cl.Create(context.Background(), &secretSpec)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created secret...")

	// Get Secret
	fmt.Println("Getting secret...")
	var secret corev1.Secret
	if err = cl.Get(context.Background(), client.ObjectKey{
		Namespace: "himanshu",
		Name:      "demo-secret",
	}, &secret); err != nil {
		cleanUpControllerResources("secret", secretSpec)
		panic(fmt.Errorf("Failed to get Secret: %v", err))
	}
	fmt.Printf("Latest Secret: %s \n", secret.Name)

	// Update Secret
	fmt.Println("Updating secret...")
	secret.Data["username"] = []byte("bXktc2VjcmV0LWFwcA==") // change username
	err = cl.Update(context.Background(), &secret)
	if err != nil {
		cleanUpControllerResources("secret", secret)
		panic(err)
	}
	fmt.Println("Updated secret...")

	// Verifying Update
	fmt.Println("Verifying Update...")
	if err = cl.Get(context.Background(), client.ObjectKey{
		Namespace: "himanshu",
		Name:      "demo-secret",
	}, &secret); err != nil {
		cleanUpControllerResources("secret", secret)
		panic(fmt.Errorf("Failed to get Secret: %v", err))
	}
	if bytes.Compare(secret.Data["username"], []byte("bXktc2VjcmV0LWFwcA==")) == 0 {
		fmt.Println("Verified Successfully")
	} else {
		cleanUpControllerResources("secret", secret)
		panic(fmt.Errorf("Verification failed. Secret found: %b, expected %b", secret.Data["username"],
			[]byte("bXktc2VjcmV0LWFwcA==")))
	}


	// Delete Secret
	fmt.Println("Deleting secret...")
	err = cl.Delete(context.Background(), &secret)
	if err != nil {
		panic(err)
	}
	fmt.Println("Deleted secret...")
}
