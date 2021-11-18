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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

func (r *SandboxReconciler) createNamespace(ctx context.Context, req ctrl.Request, sandbox1 *devtasksv1.Sandbox) (ctrl.Result, error) {

	sandbox1 = &devtasksv1.Sandbox{}

	nsSpec := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s%s", "ns-", strings.ToLower(sandbox1.Spec.Name)),
		},
	}
	err := r.Client.Create(ctx, nsSpec)
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *SandboxReconciler) updateNamespace(req ctrl.Request, ctx context.Context, sandbox *devtasksv1.Sandbox) (ctrl.Result, error) {
	//namespace1 := req.Namespace
	//namespaces := corev1.NamespaceAll

	ns := corev1.NamespaceList{}
	sbName := sandbox.Spec.Name

	for _, b := range ns.Items {
		if strings.Contains(b.Name, sbName) {
			nsSpec := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: fmt.Sprintf("%s%s", "ns-", strings.ToLower(sbName)),
				},
			}
			err := r.Client.Update(ctx, nsSpec)
			if err != nil {
				log.Log.WithValues("Can not update namespace")
				return ctrl.Result{}, err
			}
		}
	}
	return ctrl.Result{}, nil
}

func (r *SandboxReconciler) deleteNamespace(req ctrl.Request, ctx context.Context, sandbox *devtasksv1.Sandbox) (ctrl.Result, error) {
	//namespace1 := req.Namespace
	//namespaces := corev1.NamespaceAll

	ns := &corev1.Namespace{}
	//sbName := sandbox.Spec.Name

	nsName := r.Client.Get(ctx, types.NamespacedName{Name: "default"}, ns)
	fmt.Println("logging nslist ", nsName)

	//r.Log.Info("client list by client.get", nsName)

	nslist := r.Client.List(ctx, &corev1.NamespaceList{})
	fmt.Println("logging nslist ", nslist)
	//r.Log.WithValues("logging nslist", "namespace")
	//r.Log.Info("client list by client.list ", nslist)

	/*
			nsList, err := corev1.NamespaceList.Items
		        Namespaces().
		        List(context.Background(), metav1.ListOptions{})
			nsName = &corev1.NamespaceList()
			for _, b := range ns.Items {
				if strings.Contains(b.Name, sbName) {
					nsSpec := &corev1.Namespace{
						ObjectMeta: metav1.ObjectMeta{
							Name: fmt.Sprintf("%s%s", "ns-", strings.ToLower(sbName)),
						},
					}
					err := r.Client.Update(ctx, nsSpec)
					if err != nil {
						log.Log.WithValues("Can not update namespace")
						return ctrl.Result{}, err

					}
				}
			}
	*/

	return ctrl.Result{}, nil
}
