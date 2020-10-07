package clientgo

import (
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
