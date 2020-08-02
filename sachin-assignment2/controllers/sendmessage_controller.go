/*
Copyright 2020 Sachin.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	sendmessagev1 "sendmessage/io/m/api/v1"
)

// SendMessageReconciler reconciles a SendMessage object
type SendMessageReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=sendmessage.sendmessage.io,resources=sendmessages,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=sendmessage.sendmessage.io,resources=sendmessages/status,verbs=get;update;patch

func (r *SendMessageReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("sendmessage", req.NamespacedName)

	sendMessage := &sendmessagev1.SendMessage{}
	if err := r.Get(ctx, req.NamespacedName, sendMessage); err != nil{
		log.Error(err, "unable to fetch sendMessages")

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if sendMessage.Status.Status == ""{
		if sendMessage.Spec.MessageCarrier == "whatsapp"{
			// Sending Message by whatsapp
			fmt.Println("Sending message by whatsapp")
			sendMessage.Status.Status = "Sent"

		}else if sendMessage.Spec.MessageCarrier == "telegram"{
			// Sending Message by Telegram
			fmt.Println("Sending message by telegram")
			sendMessage.Status.Status = "Sent"
		}
	}else{
		// Delete the resource, as the message has been already been delivered
		sendMessage.Status.Status = "Finished"
	}

	err := r.Status().Update(ctx, sendMessage)
	if err != nil{
		log.Error(err, "Unable to update the resource instance")
	}

	// Deleting the Resource if it has been finished.
	if sendMessage.Status.Status == "Finished"{
		err := r.Delete(ctx, sendMessage)
		if err != nil{
			log.Error(err, "Unable to delete the resource")
		}
	}

	return ctrl.Result{}, nil
}

func (r *SendMessageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&sendmessagev1.SendMessage{}).
		Complete(r)
}
