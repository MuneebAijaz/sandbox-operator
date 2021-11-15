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
	finalizerUtil "github.com/stakater/operator-utils/util/finalizer"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	userFinalizer string = "http://tenantoperator.stakater.com/tenant"
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

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcilerUtil.DoNotRequeue()
		}
		// Error reading channel, requeue
		return reconcilerUtil.RequeueWithError(err)
	}

	/*
		//&devtasksv1.User.GetDeletionTimestamp()
		fmt.Println(user1.DeletionTimestamp)

		if user1.GetDeletionTimestamp() != nil {
			r.Log.Info("Deletion timestamp found for channel " + req.Name)
			if finalizerUtil.HasFinalizer(user1, userFinalizer) {
				return r.finalizeChannel(ctx, req, user1)
			}
			// Finalizer doesn't exist so clean up is already done
			return reconcilerUtil.DoNotRequeue()
		}

		// Add finalizer if it doesn't exist
		if !finalizerUtil.HasFinalizer(user1, userFinalizer) {
			r.Log.Info("Adding finalizer for channel " + req.Name)

			finalizerUtil.AddFinalizer(user1, userFinalizer)

			err := r.Client.Update(ctx, user1)
			if err != nil {
				return reconcilerUtil.ManageError(r.Client, user1, err, true)
			}
		}

	*/
	//log.Log.Info("throwing values", user1.Spec.Name)

	//namespaceName := req.NamespacedName

	r.Log.Info("namespace", metav1.NamespaceAll)

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

func (r *UserReconciler) finalizeChannel(ctx context.Context, req ctrl.Request, user1 *devtasksv1.User) (ctrl.Result, error) {
	if user1 == nil {
		return reconcilerUtil.DoNotRequeue()
	}

	sandboxName := user1.Name
	log := r.Log.WithValues("Name", sandboxName)

	err := r.Client.Get(ctx, req.NamespacedName, user1)

	if err != nil && err.Error() != "channel_not_found" && err.Error() != "already_archived" {
		return reconcilerUtil.ManageError(r.Client, user1, err, false)
	}

	finalizerUtil.DeleteFinalizer(user1, userFinalizer)
	log.V(1).Info("Finalizer removed for channel")

	err = r.Client.Update(context.Background(), user1)
	if err != nil {
		return reconcilerUtil.ManageError(r.Client, user1, err, false)
	}

	return reconcilerUtil.DoNotRequeue()

}

//err = r.Get(ctx, req.NamespacedName, sandbox1)
//log.Log.Info("sandbox1", sandbox1)

func (r *UserReconciler) updateUser(ctx context.Context, user *devtasksv1.User, req ctrl.Request) (ctrl.Result, error) {

	name := user.Spec.Name
	//deptname := user.Spec.DeptName
	//orgName := user.Spec.OrgName
	//sandboxCount := user.Spec.SandBoxCount

	sandbox1 := &devtasksv1.Sandbox{}
	err := r.Get(ctx, req.NamespacedName, sandbox1)

	sName := user.Status.Name
	//sSBCount := user.Status.SandBoxCount

	sbName := sandbox1.Spec.Name
	sbUserName := strings.Split(sbName, "-")[1]

	if sName != name {
		if sbUserName == sName {

		}
	}

	if err == nil {
		return ctrl.Result{}, err
	}

	log := r.Log.WithValues("name", name)

	log.Info("Updating user details")
	return reconcilerUtil.ManageSuccess(r.Client, user)
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&devtasksv1.User{}).
		Complete(r)
}
