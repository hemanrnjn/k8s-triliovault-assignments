package main

import (
	"fmt"
	"github.com/hemanrnjn/k8s-triliovault-assignments/himanshu/assignment1/clientgo"
	"github.com/hemanrnjn/k8s-triliovault-assignments/himanshu/assignment1/controller"
	"time"
)

func main() {
	fmt.Println("Executing client-go operations...")
	clientgo.DeployOps()
	clientgo.PodOps()
	clientgo.SecretOps()
	fmt.Println("client-go operations completed")

	fmt.Println("Sleeping for 1 min to finish termination of previous resources..")
	time.Sleep(time.Minute * 1)

	fmt.Println("Executing controller-runtime operations...")
	controller.DeployOps()
	controller.PodOps()
	controller.SecretOps()
	fmt.Println("controller-runtime operations completed")
}
