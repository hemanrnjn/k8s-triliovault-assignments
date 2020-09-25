//reading pods from incluster config via client-go
package main

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"time"
)

func main() {
	config,err:=rest.InClusterConfig()
	if err!=nil {
		panic(err.Error())
	}

	newclient,err:=kubernetes.NewForConfig(config)
	if err!=nil {
		panic(err.Error())
	}


	for {
		podName,err:=newclient.CoreV1().Pods("vishu").Get("nginx-pod",metav1.GetOptions{})
		if err!=nil {
			panic(err.Error())
		}
		fmt.Println(podName.Name)
		time.Sleep(10 * time.Second)
	}
}
