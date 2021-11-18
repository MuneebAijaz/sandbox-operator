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
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	//"sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	userFinalizer string = "tenantoperator.stakater.com/tenant"
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
	//_ = log.FromContext(ctx)
	log := r.Log.WithValues("user", req.NamespacedName)

	user1 := &devtasksv1.User{}
	err := r.Get(ctx, req.NamespacedName, user1)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcilerUtil.DoNotRequeue()
		}
		// Error reading user1, requeue
		return reconcilerUtil.RequeueWithError(err)
	}

	if user1.Status.Name == "" {
		user1.Status.Name = user1.Spec.Name
		user1.Status.SandBoxCount = user1.Spec.SandBoxCount
		err = r.Status().Update(ctx, user1)
		if err != nil {
			r.Log.Info("Unable to update user status")
			return reconcilerUtil.ManageError(r.Client, user1, err, false)
		}
		r.createSandboxes(ctx, user1)
	} else if user1.Spec.Name != user1.Status.Name || user1.Spec.SandBoxCount != user1.Status.SandBoxCount {
		r.Log.Info("inside the update user function")
		r.updateUser(ctx, user1, req)
	}

	// user1 is marked for deletion
	if user1.GetDeletionTimestamp() != nil {
		log.Info("Deletion timestamp found for user1 " + req.Name)
		if finalizerUtil.HasFinalizer(user1, userFinalizer) {
			return r.finalizeUser(ctx, req, user1)
		}
		// Finalizer doesn't exist so clean up is already done
		return reconcilerUtil.DoNotRequeue()
	}

	// Add finalizer if it doesn't exist
	if !finalizerUtil.HasFinalizer(user1, userFinalizer) {
		log.Info("Adding finalizer for user1 " + req.Name)
		finalizerUtil.AddFinalizer(user1, userFinalizer)
		err := r.Client.Update(ctx, user1)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, user1, err, true)
		}
	}
	return ctrl.Result{}, nil
}

func (r *UserReconciler) handleDelete(ctx context.Context, user1 *devtasksv1.User, req ctrl.Request) (ctrl.Result, error) {
	sandbox1 := &devtasksv1.Sandbox{}
	uName := user1.Spec.Name
	//sbList := &devtasksv1.SandboxList{}

	//	time.Sleep(5 * time.Second)

	r.Log.Info("inside the handleDelete function")

	for i := 1; i <= user1.Spec.SandBoxCount; i++ {

		err := r.Client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: fmt.Sprintf("%s%s%s%d", "sb-", strings.ToLower(user1.Spec.Name), "-", i)}, sandbox1)
		if err != nil {
			r.Log.Info("can not get sandbox resources in handleDelete")
			return reconcilerUtil.ManageError(r.Client, sandbox1, err, false)

		}
		r.Log.Info("scope inside for loop to delete sandboxes")
		if strings.Contains(sandbox1.Spec.Name, uName) {
			r.Client.Delete(ctx, sandbox1)
		}
	}
	return ctrl.Result{}, nil
}

func (r *UserReconciler) createSandboxes(ctx context.Context, user1 *devtasksv1.User) (ctrl.Result, error) {

	for i := 1; i <= user1.Spec.SandBoxCount; i++ {
		sandbox1 := &devtasksv1.Sandbox{

			TypeMeta: metav1.TypeMeta{
				APIVersion: "devtasks.sandboxop.com/v1",
				Kind:       "Sandbox",
			},

			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s%s%s%d", "sb-", strings.ToLower(user1.Spec.Name), "-", i),
				Namespace: user1.Namespace,
			},

			Spec: devtasksv1.SandboxSpec{
				Name: fmt.Sprintf("%s%s%s%d", "SB-", user1.Spec.Name, "-", i),
				Type: "T1",
			},
		}

		err := r.Client.Create(ctx, sandbox1)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				r.Log.Info("already exists")
			}
			return reconcilerUtil.ManageError(r.Client, sandbox1, err, true)
		}
	}

	return reconcilerUtil.DoNotRequeue()
}

func (r *UserReconciler) createMoreSandboxes(ctx context.Context, user1 *devtasksv1.User) (ctrl.Result, error) {

	for i := user1.Status.SandBoxCount; i < user1.Spec.SandBoxCount; i++ {

		sandbox1 := &devtasksv1.Sandbox{

			TypeMeta: metav1.TypeMeta{
				APIVersion: "devtasks.sandboxop.com/v1",
				Kind:       "Sandbox",
			},

			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s%s%s%d", "sb-", strings.ToLower(user1.Spec.Name), "-", i+1),
				Namespace: user1.Namespace,
			},

			Spec: devtasksv1.SandboxSpec{
				Name: fmt.Sprintf("%s%s%s%d", "SB-", user1.Spec.Name, "-", i),
				Type: "T1",
			},
		}
		err := r.Client.Create(ctx, sandbox1)
		if err != nil {
			r.Log.Info("error from more sandbox creation")
			if errors.IsAlreadyExists(err) {
				r.Log.Info("already exists")
			}
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *UserReconciler) finalizeUser(ctx context.Context, req ctrl.Request, user1 *devtasksv1.User) (ctrl.Result, error) {
	if user1 == nil {
		return reconcilerUtil.DoNotRequeue()
	}

	log := r.Log.WithValues("Name", user1)

	err := r.Client.Get(ctx, req.NamespacedName, user1)

	if err != nil {
		return reconcilerUtil.ManageError(r.Client, user1, err, false)
	}

	//sandbox1 := &devtasksv1.Sandbox{}
	//uName := user1.Spec.Name
	//uSBCount := user1.Spec.SandBoxCount

	r.handleDelete(ctx, user1, req)

	finalizerUtil.DeleteFinalizer(user1, userFinalizer)
	log.V(1).Info("Finalizer removed for user1")

	err = r.Client.Update(context.Background(), user1)
	if err != nil {
		return reconcilerUtil.ManageError(r.Client, user1, err, false)
	}
	return reconcilerUtil.DoNotRequeue()
}

func (r *UserReconciler) updateUser(ctx context.Context, user1 *devtasksv1.User, req ctrl.Request) (ctrl.Result, error) {

	r.Log.Info("inside the updateUser function")

	sandbox1 := &devtasksv1.Sandbox{}

	if user1.Spec.Name != user1.Status.Name {
		r.Log.Info("Can not change the user name")
		user1.Spec.Name = user1.Status.Name
		r.Client.Update(ctx, user1)
	}

	if user1.Spec.SandBoxCount < user1.Status.SandBoxCount {
		r.Log.Info("Updated user count can not be less than the current one. Resetting it to previous value")
		user1.Spec.SandBoxCount = user1.Status.SandBoxCount
		r.Client.Update(ctx, user1)

	} else if user1.Spec.SandBoxCount > user1.Status.SandBoxCount {
		r.Log.Info("spec count is more than status count")

		r.createMoreSandboxes(ctx, user1)
		user1.Spec.SandBoxCount = user1.Status.SandBoxCount
		r.Client.Update(ctx, user1)
	}

	sandbox1 = &devtasksv1.Sandbox{}
	err := r.Get(ctx, req.NamespacedName, sandbox1)

	if err != nil {
		return ctrl.Result{}, err
	}

	r.Log.Info("Updating user details")
	return reconcilerUtil.ManageSuccess(r.Client, user1)
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&devtasksv1.User{}).
		Complete(r)
}
