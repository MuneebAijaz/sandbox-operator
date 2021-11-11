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

	devtasksv1 "github.com/MuneebAijaz/sandbox-operator/api/v1"
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

	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Log.Info("already exists")
		}
	}

	//log.Log.Info("throwing values", user1.Spec.Name)

	for i := 0; i < user1.Spec.SandBoxCount; i++ {
		sandbox1 := &devtasksv1.Sandbox{

			TypeMeta: metav1.TypeMeta{
				APIVersion: "devtasks.sandboxop.com/v1",
				Kind:       "Sandbox",
			},

			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s%d", "sandbox-user-", user1.Spec.SandBoxCount),
				Namespace: req.Namespace,
			},

			Spec: devtasksv1.SandboxSpec{
				Name: fmt.Sprintf("%s%s%s%d", "SB-", user1.Spec.Name, "-", user1.Spec.SandBoxCount),
				Type: "T1",
			},
		}
		err = r.Client.Create(ctx, sandbox1)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				log.Log.Info("already exists")
			}
			return ctrl.Result{}, err
		}

	}

	//err = r.Get(ctx, req.NamespacedName, sandbox1)
	//log.Log.Info("sandbox1", sandbox1)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			//return reconcilerUtil.DoNotRequeue()
		}
		// Error reading channel, requeue
		//return reconcilerUtil.RequeueWithError(err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&devtasksv1.User{}).
		Complete(r)
}
