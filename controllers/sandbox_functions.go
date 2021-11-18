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
	finalizerUtil "github.com/stakater/operator-utils/util/finalizer"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// SandboxReconciler reconciles a Sandbox object
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

func (r *SandboxReconciler) handleCreate(ctx context.Context, req ctrl.Request, sandbox1 *devtasksv1.Sandbox) (ctrl.Result, error) {
	//namespace1 := req.Namespace
	//namespaces := corev1.NamespaceAll

	namespaceName := &v1.NamespaceList{}

	log.Log.Info("hello namespace", namespaceName)

	return ctrl.Result{}, nil

	//err := r.Client.Delete(ctx, namespace1)
}

func (r *SandboxReconciler) handleUpdate(req ctrl.Request, ctx context.Context, sandbox1 *devtasksv1.Sandbox) (ctrl.Result, error) {

	sandbox1 = &devtasksv1.Sandbox{}
	err := r.Get(ctx, req.NamespacedName, sandbox1)
	if err != nil {
		r.Log.Info("Can not get sandbox resource, doesnt exist")
	}

	ns := corev1.NamespaceList{}

	userName := sandbox1.Spec.Name

	for _, b := range ns.Items {
		if strings.Contains(b.Name, userName) {
			nsSpec := &devtasksv1.Sandbox{
				Spec: devtasksv1.SandboxSpec{
					Name: fmt.Sprintf("%s%s%s%d", "SB-", userName, "-", 1),
				},
			}
			err = r.Client.Update(ctx, nsSpec)
			if err != nil {
				log.Log.WithValues("Can not update namespace")
			}
		}
	}

	return ctrl.Result{}, err
}

func (r *SandboxReconciler) finalizeSandbox(ctx context.Context, req ctrl.Request, sandbox1 *devtasksv1.Sandbox) (ctrl.Result, error) {
	if sandbox1 == nil {
		return reconcilerUtil.DoNotRequeue()
	}

	log := r.Log.WithValues("Name", sandbox1)

	err := r.Client.Get(ctx, req.NamespacedName, sandbox1)

	if err != nil {
		return reconcilerUtil.ManageError(r.Client, sandbox1, err, false)
	}

	//sandbox1 := &devtasksv1.Sandbox{}
	//uName := sandbox1.Spec.Name
	//uSBCount := sandbox1.Spec.SandBoxCount

	r.deleteNamespace(req, ctx, sandbox1)

	finalizerUtil.DeleteFinalizer(sandbox1, userFinalizer)
	log.V(1).Info("Finalizer removed for sandbox1")

	err = r.Client.Update(context.Background(), sandbox1)
	if err != nil {
		return reconcilerUtil.ManageError(r.Client, sandbox1, err, false)
	}
	return reconcilerUtil.DoNotRequeue()
}

func (r *SandboxReconciler) updateSandbox(ctx context.Context, sandbox1 *devtasksv1.Sandbox, req ctrl.Request) (ctrl.Result, error) {

	sandbox1 = &devtasksv1.Sandbox{}

	if sandbox1.Spec.Name != sandbox1.Status.Name {
		r.Log.Info("Can not change the user name")
		sandbox1.Spec.Name = sandbox1.Status.Name
		r.Client.Update(ctx, sandbox1)
	}

	sandbox1 = &devtasksv1.Sandbox{}
	err := r.Get(ctx, req.NamespacedName, sandbox1)

	if err != nil {
		return ctrl.Result{}, err
	}

	r.Log.Info("Updating user details")
	return reconcilerUtil.ManageSuccess(r.Client, sandbox1)
}
