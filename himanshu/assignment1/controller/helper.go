package controller

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func cleanUpControllerResources(resource string, object interface{}) {
	fmt.Println("Cleaning up controller resources..")
	cl, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		panic(err)
	}
	if resource == "deploy" {
		deploy := object.(appsv1.Deployment)
		err := cl.Delete(context.Background(), &deploy)
		if err != nil {
			panic(err)
		}
	} else if resource == "pod" {
		pod := object.(corev1.Pod)
		err := cl.Delete(context.Background(), &pod)
		if err != nil {
			panic(err)
		}
	} else {
		secret := object.(corev1.Secret)
		err := cl.Delete(context.Background(), &secret)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Cleaned up controller resources..")
}
