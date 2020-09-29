# Assignment 2

* There was already an API (v1). So for multi-version API, I created new API using:
`kubebuilder create api --group batch --version v2 --kind CronJob`. Although I didn't make any new controller.

* When I build a new API, a new types.go gets created in v2 version and with this, a conversion.go file is created that is required to manage both the versions from the controller (in both versions). It contains `ConvertTo()` and `ConvertFrom()` functions. Also, due to this, we need to create a marker `+kubebuilder:storageversion` in first version types.go that contains the information of both the versions.

* Contoller Runtime helps in this conversion and I used validating and defaulting webhooks for this. However, this webhook implementation is done by controller-runtime only. A hub and spoke model is used for this. So, in this implementation, the hub is v1 and spoke is v2. A hub is the main version which first gets converted into from other version and all the remaining versions are known as spokes. A Hub is an interface that stores a `runtime.Object` and `Hub()`. So `Hub()` is now in types.go of v1.

* Controller-gen is mainly for the autogeneration of manifests (like CRDs and Deepcopy objects). It has its own CLI interface and uses generators to work (like `+webhook`, `crd`, `+rbac`, `+optional` etc.)

* Previously the v1 API had fields like `Concurrency Policy`, `Schedule`, `JobTemplate` etc. Now I made new fields taking `CronField` type and using `CronSchedule` struct in which all types of cron schedules are mentioned. (* * * * * )


## Markers

Markers are just comments of single line that represent some addtional information of our implementation (like fields) and used by `controller-runtime` and `make` to make changes in CRDs for additional recommedations or restrictions. It is started with a + symbol with some small configurations (like optional) added.
Sample Syntax: `+path:to:marker:arg1=val,arg2=val2`
The arguments can be of strings, ints, bools, slices, or maps data type.

Marker Types:
1. Empty: These are similar to boolean flags on CLI (--). Eg. `+kubebuilder:validation:Optional`
2. Anonymous: These take only one value as argument. Eg. `+kubebuilder:validation:Minimum=0`
3. Multi-Option: These take multiple arguments. Ordering is not compulsory. Eg. "+kubebuilder:printcolumn:name="Data Size",type=string,JSONPath=`.status.size`"


## Subresources

Subresources are special kind of endpoints that have a suffix appended to their path of the normal resource. For example, the pod standard HTTP path is `/api/v1/namespace/namespace/pods/name`. Pods have a number of subresources, such as /logs, /portforward, /exec, and /status.

There are mainly two types of subresources : `/scale` and `/status`. Both are not enabled by default. We can enable the `status` subresource by `// +kubebuilder:subresource:status`. When I write this genarator (marker) within the types.go file near `status` struct and perform `make install`, the CRD gets updated with these fields:
`subresources:
    status: {}`

Once I create this cron job status instance, an endpoint depicting the status or the current state of the resource is created. I can curl this endpoint to check the status. Only controller is allowed to change its contents. The 3 types of requests (verbs) `PUT/POST/PATCH` ignores all the changes in status subresource and does everything without interfering status.

**While I can change the fields and values in `spec` of types.go OR .yaml of CR, if I change the contents of `status`, it will throw and error and won't change. This is because the `PUT` & `POST` will avoid accidentally overwriting the status in read-modiy-write scenarios. If it overwrites the status, it will corrupt the data that actually represented the status of that resource.**

Thus in the CRD, the `status` struct is empty by default. However, to get the contents that get stored in the `status`, nothing else needs to be done. The `Get()` method of `client` package returns the entire object that holds the status of the actual resource.


Used Links:
* https://book.kubebuilder.io/multiversion-tutorial/tutorial.html
* https://book.kubebuilder.io/reference/controller-gen.html
* https://book.kubebuilder.io/reference/markers.html
* https://www.oreilly.com/library/view/programming-kubernetes/9781492047094/ch04.html
* https://book-v1.book.kubebuilder.io/basics/status_subresource.html
* https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#status-subresource
* https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#spec-and-status