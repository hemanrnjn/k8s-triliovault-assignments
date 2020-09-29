# Assignment 3

## Custom Controllers

A controller is a component and a small system provided by Kubernetes that handles events over the resources/workloads and reconciles over them to make some decisions according to the given code. However, any in-built controller will reconile over the in-built resources.

A custom controller is built to reconcile over certain in-built and/or custom resources and perform some operation based on the code written inside the controller. It has mainly defined two resource structs that contain what to be checked and what needs to be maintained. Those structs are: `spec{}` and `status{}`. 
* `spec{}` holds fields that define what the desired condition is for the resource.
* `status{}` holds fields that define what the current condition is of the resource.
A controller tries to find out what the current status of resource is and tries to match it with what mentioned in spec.

The controller is built and operated by two go files:
* main.go — this is the entry point for the controller as well as where everything gets started.
* controller.go — It holds the controller struct and methods, and where all of the work is done with controller loop as `Reconcile()`

The controller works on the concept of working on events. So it checks on events and performs actions accordingly. An event is the combination of an action (create, update, or delete) and a resource key (typically in the format of namespace/name). While checking the events, the controller mainly devotes itself to the queue. A queue is a sequential store that the Kubernetes (informer) handles to decide which resource is next to be operated.
However, if the error is returned from `Reconcile()`, then it means that `Reconcile()` will be requeued (using logging). There are also two components of a controller: Informer/SharedInformer and Workqueue.

So, the event flow from informer goes like this:
`Controller Manager, API server, ETCD -> Informer -> Queue -> Controller`
This works vice-versa too.

### Creating a new controller

I now build a new controller while creating a new API with this command:
`kubebuilder create api --group mygroup.k8s.io" --version v1beta1 --kind MyKind`
`Create Resource [y/n]
y
Create Controller [y/n]
y`

Using a `CreateOrUpdate()` method from controller-runtime, it will build or update a resource as mentioned in the controller.go file. `CreateOrUpdate()` first calls `Get()` on the object. If the object does not exist, `Create()` will be called. If it does exist, `Update()` will be called. Just before calling either `Create()` or `Update()`, the mutate callback will be called. Mutate callback is where are the changes to callback must be performed.

The `OwnerReferences` field makes sure that whatever is passed to that field, that reference takes the responsibility of handling the resource. Also, due to presence of this field, in case if the custom resource gets deleted then the deployment will get deleted too and later be taken out via garbage-collector. However, in this program, I have passed Controller Reference in this field of resource. Even if owner gets deleted, the resource will get deleted too.

* In the controller, I tracked the events of the custom resource and managed to clean up all the previous deployments that were made for this kind. It similarly records events (via `r.Recorder.Eventf()`) and deletes the old deployments using `r.Client.Delete()` Also, the `r.Client.Update()` updates the replicas (pods) of the deployment. Did CRUD operations on them.
* I added the owner reference to the custom resource via the `buildDeployment()` which contains the template of the deployment with a struct field `OwnerReferences` that takes `NewControllerRef()`. In this function, the controller reference is passed as owner reference. Also, there is `Owns()` inside `NewControllerManagedBy(mgr)` that shows which resource it is handling.
* `FieldIndexer` is an interface that knows how to index over a particular field of resource manifest that can later be used by field selector. To use this, a `GetFieldIndexer()` is used within `SetupWithManager()` which indexes over the `deploymentOwnerKey` of resource and checks if the owner matches or not.
* Predicate is a `controller-runtime` package that contains functions that are used by controllers to filter events before they are handled by the event handler functions. These functions are built by handling the events from `handler.EventHandler`. Predicates are defined in `SetupWithManager()` to create, update and delete events. Every predicate returns either a `true` or `false` telling that the operation can be done or not. After defining the our functions that take `Event`, we can use predicate to run them or not using a struct:
  `p := predicate.Funcs{
		CreateFunc: createFunction,
		DeleteFunc: deleteFunction,
		UpdateFunc: updateFunction,
	}`
and mentioning `WithEventFilter(p)` in `NewControllerManagedBy(mgr)`.
* A combination of `kind`,`group` and `version` is known as GVK. We can define it within a struct:
`type Gvk struct {
	Group   string `json:"group,omitempty" yaml:"group,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	Kind    string `json:"kind,omitempty" yaml:"kind,omitempty"`
}`
And use it like:
`func FromKind(k string) Gvk {
	return Gvk{
		Kind: k,
	}
}`

## Shared Informer

Shared Informer is a caching layer that's in-built and present before ETCD. So, any request that server through API-server will first hit tot he shared infromer and if it doesn't find hat it wants or can't perform operation there, then request will move on to ETCD. Whenever there is an update to shared informer, events get generated. Controller reconciles on these events. It can also be regarded as a queue shared across cluster.

Every resource or CR that we create gets queued into shared informer. Each controller acts as broker which will pick up objects from the queue based on type of resource and state of resource. If any controller goes down for some reason, the queue remains intact. If the cluster goes down, queue can be rebuilt. Controller acquires mutation lock on a queue object as soon as it starts processing it. So in case we are running multiple replicas of our controller, there is a guarantee that the object will be processed only once, making our transaction idempotent.

Controller Flow using Shared Informer:
`User -modifies-> any resource -gets notified via informers to-> Controller -checks status and reconciles-> changes made as required`
The SharedInformer can't track where each controller is up to (because it's shared), so the controller must provide its own queuing and retrying mechanism (if required). Hence, most Resource Event Handlers simply place items onto a per-consumer workqueue.

SharedInformer watches for changes on the current state of Kubernetes objects and sends events to Workqueue where events are then popped up by worker(s) to process. It is inside `cache` package. Also, it contains a single local cache for all the resources for all the controllers. It maintains a local cache exposed by GetStore() and by GetIndexer(). So, multiple controllers can act over resources simultaneously without duplication of data.

SharedInformer only creates a single watch on the upstream server, regardless of how many downstream consumers are reading events from the informer. It has already provided hooks to receive notifications of adding, updating and deleting a particular resource.

We can use it using: 
* `import "k8s.io/client-go/informers"`
* `factory := informers.NewSharedInformerFactory(clientset, 0)`
* `informer := factory.Core().V1().Pods().Informer()`

It has a `sharedIndexInformer{}` struct which contains indexer and controller fields. An `Indexer` interface is also present indexing and sending/receiving the indexes.

### Reason for using Shared Informer

The controller goes to the API server and ETCD every time for checking and updating resources. Now, this becomes time consuming too quickly in case if another request is also being handled by API server. So, when any resource is already persisted in ETCD and no change is made into it and any request follows the condition of checking the resource, then Shared Informer quickly sends that info to the request short the timing process. Its like a linkage to the actual storage (ETCD).  
 
Used Links:
* https://trstringer.com/extending-k8s-custom-controllers/
* https://engineering.pivotal.io/post/gp4k-kubebuilder-lessons/
* https://book.kubebuilder.io/multiversion-tutorial/api-changes.html
* https://github.com/jetstack/kubebuilder-sample-controller
* https://github.com/kubernetes-sigs/kustomize/blob/master/api/resid/gvk.go
* https://engineering.bitnami.com/articles/a-deep-dive-into-kubernetes-controllers.html
* https://medium.com/@muhammet.arslan/write-your-own-kubernetes-controller-with-informers-9920e8ab6f84
* https://gianarb.it/blog/kubernetes-shared-informer
* https://github.com/kubernetes/client-go/blob/master/tools/cache/shared_informer.go