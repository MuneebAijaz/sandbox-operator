/*
Copyright 2021.

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
	"strings"

	devtasksv1 "github.com/MuneebAijaz/sandbox-operator/api/v1"
	"github.com/go-logr/logr"

	//,finalizerUtil "github.com/stakater/operator-utils/util/finalizer"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// UserReconciler reconciles a User object
type UserReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=devtasks.sandboxop.com,resources=users,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=devtasks.sandboxop.com,resources=users/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=devtasks.sandboxop.com,resources=users/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the User object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile

func (r *UserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	user1 := &devtasksv1.User{}
	err := r.Get(ctx, req.NamespacedName, user1)

	log1 := r.Log.WithValues("get user", user1.Name)

	if err != nil {
		log1.Info("User CR does not exists")
	}

	//log.Log.Info("throwing values", user1.Spec.Name)

	for i := 1; i <= user1.Spec.SandBoxCount; i++ {
		sandbox1 := &devtasksv1.Sandbox{

			TypeMeta: metav1.TypeMeta{
				APIVersion: "devtasks.sandboxop.com/v1",
				Kind:       "Sandbox",
			},

			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s%s%s%d", "sb-", strings.ToLower(user1.Spec.Name), "-", i),
				Namespace: req.Namespace,
			},

			Spec: devtasksv1.SandboxSpec{
				Name: fmt.Sprintf("%s%s%s%d", "SB-", user1.Spec.Name, "-", i),
				Type: "T1",
			},
		}
		err = r.Client.Create(ctx, sandbox1)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				r.Log.Info("already exists")
			}
			return ctrl.Result{}, err
		}

	}
	return ctrl.Result{}, nil
}

/*

func (r *UserReconciler) finalizeChannel(req ctrl.Request, sandbox *devtasksv1.Sandbox) (ctrl.Result, error) {
	if sandbox == nil {
		return reconcilerUtil.DoNotRequeue()
	}

	sandboxName := sandbox.Name
	log := r.Log.WithValues("sandbox name", sandboxName)

	err := r.Client.Get(ctx,req.NamespacedName, sandbox)

	if err != nil && err.Error() != "channel_not_found" && err.Error() != "already_archived" {
		return reconcilerUtil.ManageError(r.Client, channel, err, false)
	}

	finalizerUtil.DeleteFinalizer(channel, channelFinalizer)
	log.V(1).Info("Finalizer removed for channel")

	err = r.Client.Update(context.Background(), channel)
	if err != nil {
		return reconcilerUtil.ManageError(r.Client, channel, err, false)
	}

	return reconcilerUtil.DoNotRequeue()
}

*/

//err = r.Get(ctx, req.NamespacedName, sandbox1)
//log.Log.Info("sandbox1", sandbox1)

func (r *UserReconciler) updateUser(ctx context.Context, user *devtasksv1.User) (ctrl.Result, error) {

	name := user.Spec.Name
	log := r.Log.WithValues("channelID", name)

	log.Info("Updating channel details")

	/*

		name := channel.Spec.Name
		users := channel.Spec.Users
		topic := channel.Spec.Topic
		description := channel.Spec.Description

		_, err := r.SlackService.RenameChannel(channelID, name)
		if err != nil {
			log.Error(err, "Error renaming channel")
			return reconcilerUtil.ManageError(r.Client, channel, err, false)
		}

		_, err = r.SlackService.SetTopic(channelID, topic)
		if err != nil {
			log.Error(err, "Error setting channel topic")
			return reconcilerUtil.ManageError(r.Client, channel, err, false)
		}

		_, err = r.SlackService.SetDescription(channelID, description)
		if err != nil {
			log.Error(err, "Error setting channel description")
			return reconcilerUtil.ManageError(r.Client, channel, err, false)
		}

		err = r.SlackService.InviteUsers(channelID, users)
		if err != nil {
			log.Error(err, "Error inviting users to channel")
			return reconcilerUtil.ManageError(r.Client, channel, err, false)
		}

		err = r.SlackService.RemoveUsers(channelID, users)
		if err != nil {
			log.Error(err, "Error removing users from the channel")
			return reconcilerUtil.ManageError(r.Client, channel, err, false)
		}
	*/
	return reconcilerUtil.ManageSuccess(r.Client, user)
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&devtasksv1.User{}).
		Complete(r)
}
