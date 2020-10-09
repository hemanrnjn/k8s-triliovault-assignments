package clientgo

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getLocalConfig() *rest.Config {
	//var kubeconfig *string
	//if home := homedir.HomeDir(); home != "" {
	//	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	//} else {
	//	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	//}
	//flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", "/home/himanshu/.kube/config")
	if err != nil {
		panic(err)
	}
	return config
}

func getInClusterConfig() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	return config
}

func getClientSet(inCluster bool) *kubernetes.Clientset {
	var config *rest.Config
	if inCluster {
		config = getInClusterConfig()
	} else {
		config = getLocalConfig()
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientset
}

func cleanUpClientGoResources(resource string) {
	fmt.Println("Cleaning up client-go resources..")
	if resource == "deploy" {
		client := getClientSet(true).AppsV1().Deployments("himanshu")
		deletePolicy := metav1.DeletePropagationForeground
		deleteErr := client.Delete(context.TODO(), "demo-deploy", metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		if deleteErr != nil {
			panic(deleteErr)
		}
	} else if resource == "pod" {
		client := getClientSet(true).CoreV1().Pods("himanshu")
		deletePolicy := metav1.DeletePropagationForeground
		deleteErr := client.Delete(context.TODO(), "demo-pod", metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		if deleteErr != nil {
			panic(deleteErr)
		}
	} else {
		client := getClientSet(true).CoreV1().Secrets("himanshu")
		deletePolicy := metav1.DeletePropagationForeground
		deleteErr := client.Delete(context.TODO(), "demo-secret", metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		if deleteErr != nil {
			panic(deleteErr)
		}
	}
	fmt.Println("Cleaned up client-go resources..")
}
