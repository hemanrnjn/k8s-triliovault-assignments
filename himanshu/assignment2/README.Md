## Subresources

Subresources are special HTTP endpoints, using a suffix appended to the HTTP path of the normal resource.
For example, the pod standard HTTP path is `/api/v1/namespace/namespace/pods/name`. 
Pods have a number of subresources, such as `/logs`, `/portforward`,`/exec`, and `/status`. The corresponding subresource HTTP paths are:

```
/api/v1/namespace/namespace/pods/name/logs

/api/v1/namespace/namespace/pods/name/portforward

/api/v1/namespace/namespace/pods/name/exec

/api/v1/namespace/namespace/pods/name/status
```


The subresource endpoints use a different protocol than the main resource endpoint.

Custom resources support two subresources:
`/scale` and `/status`. Both are opt-in—that is, they must be explicitly enabled in the CRD.

#### Status subresource

The `/status` subresource is used to split the user-provided specification of a CR instance from the controller-provided status. The main motivation for this is privilege separation:

```
The user usually should not write status fields.

The controller should not write specification fields.
```

The RBAC mechanism for access control does not allow rules at that level of detail. Those rules are always per resource.
The `/status` subresource solves this by providing two endpoints that are resources on their own.
Each can be controlled with RBAC rules independently. This is often called a spec-status split.

Resources (including custom resources) that have a `/status` subresource have changed semantics, also for the main resource endpoint:

```
They ignore changes to the status on the main HTTP endpoint during create (the status is just dropped during a create) and updates.

Likewise, the /status subresource endpoint ignores changes outside of the status of the payload. A create operation on the /status endpoint is not possible.

Whenever something outside of metadata and outside of status changes (this especially means changes in the spec),
the main resource endpoint will increase the metadata.generation value. This can be used as a trigger for a controller indicating that the user desire has changed.
```

Note that usually both spec and status are sent in update requests, but technically you could leave out the respective other part in a request payload.

Also note that the `/status` endpoint will ignore everything outside of the status, including metadata changes like labels or annotations.

The spec-status split of a custom resource is enabled as follows:

```
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
spec:
  subresources:
    status: {}
```

Note here that the status field in that YAML fragment is assigned the empty object. This is the way to set a field that has no other properties.
 
Just writing:
```
subresources:
  status:
```

will result in a validation error because in YAML the result is a null value for status.

Source [O'Reilly Book: Programming Kubernetes by Michael Hausenblas, Stefan Schimanski](https://www.oreilly.com/library/view/programming-kubernetes/9781492047094/ch04.html)