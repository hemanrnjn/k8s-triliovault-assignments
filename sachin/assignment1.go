package main

import (
	"bufio"
	"context"
	//"flag"
	"fmt"
	"os"
	//"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
	"k8s.io/apimachinery/pkg/util/intstr"
	config_cr "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	/*var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()*/
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	//config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Create First Resource
	fmt.Println("Creating first resource in default namespace...")

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "demo",},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "demo",},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	// List Deployments
	prompt()
	fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}

	// Update Deployment
	prompt()
	fmt.Println("Updating deployment...")

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(context.TODO(), "demo-deployment", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		result.Spec.Replicas = int32Ptr(1)                           // reduce replica count
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated deployment...")

	// List Deployments
	prompt()
	fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
	list, err = deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}

	// Delete Deployment
	prompt()
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), "demo-deployment", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")


	// Creating Second Resource
	fmt.Println("Creating Second Resource")
	namespaceClient := clientset.CoreV1().Namespaces()

	// Create Namespace
	fmt.Println("Creating namespace...")

	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-namespace",
		},
	}
	namespaceresult, err := namespaceClient.Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created namespace %q.\n", namespaceresult.GetObjectMeta().GetName())

	// List Namespace
	prompt()
	fmt.Println("Listing Namespaces")
	namespacelist, err := namespaceClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil{
		panic(err)
	}
	for _, namespaceL := range namespacelist.Items{
		fmt.Println(namespaceL.Name)
	}

	// Create Third Resource
	fmt.Printf("Creating deployments in namespace %q:\n", "test-namespace")

	deploymentsClient = clientset.AppsV1().Deployments("test-namespace")

	deployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "demo",},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "demo",},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err = deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	// List Deployments
	prompt()
	fmt.Printf("Listing deployments in namespace %q:\n", "test-namespace")
	list, err = deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}

	// Update Deployment
	prompt()
	fmt.Println("Updating deployment...")

	retryErr = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(context.TODO(), "demo-deployment", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		result.Spec.Replicas = int32Ptr(5)                           // increase replica count
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated deployment...")


	// Create Third Resource
	fmt.Printf("Creating service in namespace %q:\n", "test-namespace")

	serviceNpClient := clientset.CoreV1().Services("test-namespace")

	serviceNp :=  &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "service-np",
		},
		Spec: apiv1.ServiceSpec{
			//Selector: &metav1.LabelSelector{
			//	MatchLabels: map[string]string{"app": "demo",},
			//},
			Ports: []apiv1.ServicePort{
				{
					Port: 80,
					TargetPort: intstr.FromInt(80),
					NodePort: 31000,
				},
			},
			Type: apiv1.ServiceTypeNodePort,
		},
	}

	// Create Node Port Ip
	fmt.Println("Creating service")
	serviceResult, err := serviceNpClient.Create(context.TODO(), serviceNp, metav1.CreateOptions{})
	if err != nil{
		panic(err)
	}
	fmt.Printf("Created service %q.\n", serviceResult.GetObjectMeta().GetName())

	// List Deployments
	prompt()
	labelSelector := fmt.Sprintf("app=%s", "demo")
	fmt.Println(labelSelector)
	fmt.Printf("Listing deployments using label selector in namespace %q:\n", "test-namespace")
	list, err = deploymentsClient.List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector,})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}

	// Delete Deployment
	prompt()
	fmt.Println("Deleting deployment...")
	deletePolicy = metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), "demo-deployment", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")

	// Delete Service
	prompt()
	fmt.Println("Deletig Service...")
	deleteServiceNpPolicy := metav1.DeletePropagationForeground
	if err := serviceNpClient.Delete(context.TODO(), "service-np", metav1.DeleteOptions{
		PropagationPolicy: &deleteServiceNpPolicy,
	}); err != nil{
		panic(err)
	}
	fmt.Println("Deleted Service")

	// Delete Namespace
	prompt()
	fmt.Println("Deletig Namespace...")
	deleteNameSpacePolicy := metav1.DeletePropagationForeground
	if err := namespaceClient.Delete(context.TODO(), "test-namespace", metav1.DeleteOptions{
		PropagationPolicy: &deleteNameSpacePolicy,
	}); err != nil{
		panic(err)
	}
	fmt.Println("Deleted Namespace")

	// Using controller-runtime client
	prompt()
	fmt.Println("Client-go work has been finished here, will start with controller-runtime client")
	fmt.Println("Listing out all the pods in sachin namespace")
	cl, err := client.New(config_cr.GetConfigOrDie(), client.Options{})
	if err != nil{
		panic(err)
	}
	podList := &apiv1.PodList{}
	err = cl.List(context.Background(), podList, client.InNamespace("sachin"))
	if err != nil{
		panic(err)
	}

	for _, pdl := range podList.Items {
                fmt.Printf(" * %s \n", pdl.Name)
        }

	// Creating a deployment using controller-runtime client
	prompt()
	fmt.Println("Creating Deployment using controller-runtime client")
	deployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
			Namespace: "sachin",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "demo",},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "demo",},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	err = cl.Create(context.Background(), deployment)
	if err != nil{
		panic(err)
	}

	fmt.Println("Deployment created using controller-runtime in namespace sachin")

	// listing out all the deployments in sachin namespace
	prompt()
	fmt.Println("Listing out all the deployments in sachin namespace.")

	deploymentList := &appsv1.DeploymentList{}
	err = cl.List(context.Background(), deploymentList, client.InNamespace("sachin"))
	if err != nil{
		panic(err)
	}

	for _, dpl := range deploymentList.Items {
                fmt.Printf(" * %s \n", dpl.Name)
        }

	// Deleting the last created deployment from sachin namespace
	prompt()
	fmt.Println("Deleting the last created deployment from sachin namespace")
	err = cl.Delete(context.Background(), deployment)
	if err != nil{
		panic(err)
	}

	// listing out all the deployments in sachin namespace
	prompt()
	fmt.Println("Listing out all the deployments in sachin namespace.")

	deploymentList = &appsv1.DeploymentList{}
	err = cl.List(context.Background(), deploymentList, client.InNamespace("sachin"))
	if err != nil{
		panic(err)
	}

	for _, dpl := range deploymentList.Items {
                fmt.Printf(" * %s \n", dpl.Name)
        }


}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}

func int32Ptr(i int32) *int32 { return &i }

