# Assignement1 

This is the first assignment and then the other one. Explore Client-go library from k8s to interact with its resources
Requirements:
1. It will be a simple binary which will initialize the clients
2. create any 5 resources of the k8s
3. Need to perform some crud operations on them
4. try to perform the operations in default namespace and different namespaces
5. Then perform the operations on the objects by getting them using some label-selector
6. Also make changes in spec as well as status and learn about the subresources
7. Also check controller-runtime clients after all the above and try out the same operations with it



First of all, i have created a kubeconfg with help of inbuilt functions. Then initialize Clientset.

Then with help of clientset, i have created first resource of a deployment in default namesapce. Then after successfully created deployments, i have list out all the deployments. Where we can see the deployment which we have created previously can be seen here in the list. 
After verifying, that the deployment has been successfully built i updated the spec of the deployment and after that i have deleted that resource.


Then, i created a namespace resource which we can verify in the list which i have displayed. 

Once, we have verified that the namespace has been created then i have again created a deployment in the same namespace which i have created previously. 
That thing, we can verifying the listing of deployments in the same namespace. After which i have updated the same deployment. 


After which, i have created a NodePort service, and then i list out the services and then we can see that the service is visible in the list. 
After these steps, i deleted the deployments, deleted Service and deleted namespace. 

