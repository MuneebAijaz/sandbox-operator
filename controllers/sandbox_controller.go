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

	devtasksv1 "github.com/MuneebAijaz/sandbox-operator/api/v1"
	"github.com/go-logr/logr"
	finalizerUtil "github.com/stakater/operator-utils/util/finalizer"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	NamespaceFinalizer string = "tenantoperator.stakater.com/namespace"
)

// SandboxReconciler reconciles a Sandbox object
type SandboxReconciler struct {
	Log logr.Logger
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=devtasks.sandboxop.com,resources=sandboxes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=devtasks.sandboxop.com,resources=sandboxes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=devtasks.sandboxop.com,resources=sandboxes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Sandbox object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *SandboxReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("sandbox", req.NamespacedName)

	sandbox1 := &devtasksv1.Sandbox{}
	err := r.Get(ctx, req.NamespacedName, sandbox1)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcilerUtil.DoNotRequeue()
		}
		// Error reading sandbox1, requeue
		return reconcilerUtil.RequeueWithError(err)
	}

	// sandbox1 is marked for deletion
	if sandbox1.GetDeletionTimestamp() != nil {
		log.Info("Deletion timestamp found for sandbox1 " + req.Name)
		if finalizerUtil.HasFinalizer(sandbox1, userFinalizer) {
			return r.finalizeSandbox(ctx, req, sandbox1)
		}
		// Finalizer doesn't exist so clean up is already done
		return reconcilerUtil.DoNotRequeue()
	}

	if sandbox1.Status.Name == "" {
		sandbox1.Status.Name = sandbox1.Spec.Name
		err = r.Status().Update(ctx, sandbox1)
		if err != nil {
			r.Log.Info("Unable to update user status")
			return reconcilerUtil.ManageError(r.Client, sandbox1, err, false)
		}
		r.createNamespace(ctx, req, sandbox1)

	} else if sandbox1.Spec.Name != sandbox1.Status.Name {
		r.updateSandbox(ctx, sandbox1, req)
	}

	// Add finalizer if it doesn't exist
	if !finalizerUtil.HasFinalizer(sandbox1, userFinalizer) {
		log.Info("Adding finalizer for sandbox1 " + req.Name)

		finalizerUtil.AddFinalizer(sandbox1, userFinalizer)

		err := r.Client.Update(ctx, sandbox1)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, sandbox1, err, true)
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SandboxReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&devtasksv1.Sandbox{}).
		Complete(r)
}
