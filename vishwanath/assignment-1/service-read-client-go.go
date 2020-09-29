//Program is run using : "./app --kubeconfig /root/.kube/config"
package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var kubernetesConfig *string
	if homeDir:=homeDir(); homeDir!="" {
		kubernetesConfig=flag.String("kubeconfig",filepath.Join(homeDir,".kube","config"),"/root/.kube/config")
	} else {
		kubernetesConfig=flag.String("kubeconfig","","/root/.kube/config")
	}
	flag.Parse()

	config,err:=clientcmd.BuildConfigFromFlags("",*kubernetesConfig)
	if err!=nil {
		panic(err.Error())
	}

	newclient,err:=kubernetes.NewForConfig(config)
	if err!=nil {
		panic(err.Error())
	}

	for {
		svcName,err:=newclient.CoreV1().Services("vishu").Get(context.TODO(),"my-service",metav1.GetOptions{})
		if err!=nil {
			panic(err.Error())
		}
		fmt.Println(svcName.Name)
		time.Sleep(10 * time.Second)
	}
}

//Already exporting ROOTPATH as "/root"
func homeDir() string {
	//if h := os.Getenv("ROOTPATH"); h != "" {
	//	return h
	//}
	return os.Getenv("ROOTPATH")
}