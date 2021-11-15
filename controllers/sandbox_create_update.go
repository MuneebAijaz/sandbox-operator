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
	//	finalizerUtil "https://github.com/stakater/operator-utils/tree/master/util/finalizer"
	//	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
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

func (r *SandboxReconciler) handleCreate(ctx context.Context, user *devtasksv1.User, req ctrl.Request, namespace *v1.Namespace) (ctrl.Result, error) {
	//namespace1 := req.Namespace
	//namespaces := corev1.NamespaceAll

	namespaceName := &v1.NamespaceList{}

	log.Log.Info("hello namespace", namespaceName)

	return ctrl.Result{}, nil

	//err := r.Client.Delete(ctx, namespace1)
}

func (r *SandboxReconciler) handleUpdate(req ctrl.Request, ctx context.Context, sandbox1 *devtasksv1.Sandbox) (ctrl.Result, error) {
	//namespace1 := req.Namespace
	//namespaces := corev1.NamespaceAll

	user1 := &devtasksv1.User{}
	err := r.Get(ctx, req.NamespacedName, user1)
	if err != nil {
		r.Log.Info("Can not get sandbox resource, doesnt exist")
	}

	ns := corev1.NamespaceList{}

	userName := user1.Spec.Name

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
