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
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"

	//	finalizerUtil "https://github.com/stakater/operator-utils/tree/master/util/finalizer"
	//	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func (r *SandboxReconciler) createNamespace(ctx context.Context, req ctrl.Request, sandbox1 *devtasksv1.Sandbox) (ctrl.Result, error) {

	err := r.Get(ctx, req.NamespacedName, sandbox1)
	if err != nil {
		r.Log.Info("can not get sandbox resource in createNamespace function")
	}

	nsSpec := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s%s", "ns-", strings.ToLower(sandbox1.Spec.Name)),
		},
	}

	log := r.Log.WithValues("Logging for namespace", nsSpec)

	err = client.Client.Create(r.Client, ctx, nsSpec)

	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Info("already exists")
		} else {
			r.Log.Info("can not create new sandboxes")
		}
		return reconcilerUtil.ManageError(r.Client, sandbox1, err, true)
	}
	return ctrl.Result{}, nil
}

func (r *SandboxReconciler) deleteNamespace(req ctrl.Request, ctx context.Context, sandbox1 *devtasksv1.Sandbox) (ctrl.Result, error) {

	ns := &corev1.Namespace{}

	log := r.Log.WithValues("Logging for namespace", ns)

	err := r.Client.Get(ctx, types.NamespacedName{Name: fmt.Sprintf("%s%s", "ns-", strings.ToLower(sandbox1.Spec.Name))}, ns)
	if err != nil {
		log.Info("unable to get the specified namespace in deleteNamespace function")
		return ctrl.Result{}, err
	}

	log.Info("getting namespace ns")

	err = r.Client.Delete(ctx, ns)
	if err != nil {
		log.Info("unable to delete the specified namespace in deleteNamespace function")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
