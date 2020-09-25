# Assignment 1

## Client-Go:

Client-Go is a powerful tool and a library of many pre-built functions that allow us to interact to Kubernetes API not in a general way (i.e. using `kubectl`) but defining our own functions to use and access resources on the Kubernetes Cluster. It is already being used with Kubebuilder and Operator-SDK. The main function of this is to create clients and then these clients interact with Kubernetes. This has two types of Config: *Incluster Config* and *Outcluster Config*. Although it is recommended to use the same `client-go` version that is being used with `k8s.io/api` and `k8s.io/apimachinery`.

### Made using InCluster Config:

Incluster Config is way of interacting with Kubernetes Cluster by interacting to the API server directly inside the cluster, which means that the code will be deployed inside the cluster (via pods) and it will access the resources. This is particularly useful when we don't want to mention the `kubeconfig` file inside the code so its not hard-coded and directly picks up the service account token from `var/run/secrets/kubernetes.io/serviceaccount`.

However, to distinguish the work done by this config, a new clusterrolebinding is created with default service account using: `kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default`.

The `client-go` gives us functions like `rest.InClusterConfig()` to start a config and `kubernetes.NewForConfig()` to create our own custom client.

Steps Performed:
1. Made incluster-config.go file
2. Made Dockerfile
3. Installed all the dependencies that were required with incluster-config.go file (like client-go, googleapis-gnostic)
4. GOOS=linux go build -o ./app . (exporting the env variable and building go code)
5. docker build -t myimage . (building dockerfile for pulling it afterwards in pods)
6. Pushed the docker image to docker hub. (using docker login)
7. Made a deployment file using the docker image name. (docker.io/username/imagename)
8. kubectl create -n vishu -f deployment.yaml (creating the deployment)
9. Checked logs for the deployed pod which printed the number of pods.

Used Links:
* https://github.com/kubernetes/client-go/issues/741#issuecomment-603440039
* https://github.com/kubernetes/client-go/tree/master/examples/in-cluster-client-configuration
* https://medium.com/swlh/clientset-module-for-in-cluster-and-out-cluster-3f0d80af79ed
* https://github.com/kubernetes/client-go/tree/master/examples/create-update-delete-deployment

However, the pod inside the code runs picks up the `KUBERNETES_SERVICE_PORT` and `KUBERNETES_SERVICE_HOST` since it requires that to run client-go and connect to API server of Kubernetes (just like connecting an app to a database).

### Made using OutCluster Config:

This is way easier than using InCluster Config. But it needs to be hard-coded and should be used with the `--kubeconfig` and/or `--master` flag while executing to tell the `client-go` where to find the Kubernetes Cluster. Even though its much easier to implement, it can't be ported to another cluster without making significant changes.

* We have to provide the entire path `kubeconfig` file. The `flag` package needs to parse it too.
* `clientcmd.BuildConfigFromFlags()` helps in creating a config from the flags that have been retrieved from `kubeconfig` path.

Steps Performed:
1. Made outcluster-config-service.go file
2. ./app --kubeconfig /root/.kube/config" (running program)

Used-links:
* https://github.com/kubernetes/client-go/tree/master/examples/out-of-cluster-client-configuration


### What I Learnt:

* Learnt that resources names must always be in lowercase.
* Incluster config and Outcluster config can be effective in specific cases and gives more freedom in interacting with Kubernetes.
* Using these configs and client-go, the code can be more structured and give better performance.
* Client-Go already has many predefined functions, packages, structs etc. for every built-in resource of Kubernetes.
* There are two options while updating any resource via client-go:
  1. Change the template spec of the resource in a function and use update function. This is like patching the resource (rollout). But it will overwrite any previous changes.
  2. We can use `retry` package that arrives with client-go and changes the resource specs until the error goes away. This is useful and also preserves changes that were made from other clients.


## Controller-Runtime

Controller-runtime is a tool and library for Kubernetes Controllers that provides interfaces and functions for reconciling over our custom Kubernetes Objects and perform CRUD operations on them. This is used by both Kubebuilder and Operator-SDK. Mainly it deals with giving interfaces for clients, managers, builders and reconcilers that are used for handling operations for the controller (webhooks, api-conversion and logging too). It has several packages inside `/pkg` that help to establish this.

* The builder builds an application controller (that we get with `kubeadm init`) and returns a `manager.Manager` to start it.
* The client has `client.Client` interface for performing CRUD operations for specfic resources.
* The manager has `manager.Manager` interface for maintaining the controller lifecycle and client.
* The reconciler has `reconcile.Reconciler` interface that handles the reconciling process and updates the resources according to the logic written and metrics derived from API server.

Although in the given Kubebuilder example, the client and managers are replaced with `ctrl`. Also, the Create, Read (list) and delete operations are done on cronjobs.

Used-links:
* https://sdk.operatorframework.io/docs/building-operators/golang/references/client/
* https://godoc.org/github.com/kubernetes-sigs/controller-runtime
* https://book.kubebuilder.io/reference/controller-gen.html