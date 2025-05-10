/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	animalsv1 "github.com/amir-aharon/goat-operator/api/v1"
)

// GoatReconciler reconciles a Goat object
type GoatReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=animals.baaaa.io,resources=goats,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=animals.baaaa.io,resources=goats/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=animals.baaaa.io,resources=goats/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Goat object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *GoatReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var goat animalsv1.Goat
	if err := r.Get(ctx, req.NamespacedName, &goat); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var podList v1.PodList
	if err := r.List(ctx, &podList, client.InNamespace(req.Namespace)); err != nil {
		return ctrl.Result{}, err
	}

	for _, pod := range podList.Items {
		if pod.Name == goat.Spec.PodName {
			fmt.Println("Baaaaaa! matched pod", pod.Name)
		}
	}

	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

func (r *GoatReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&animalsv1.Goat{}).
		Watches(
			&v1.Pod{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
				var goats animalsv1.GoatList

				// List all Goat CRs in the same namespace as the pod
				if err := r.List(ctx, &goats, client.InNamespace(obj.GetNamespace())); err != nil {
					return nil
				}

				var reqs []reconcile.Request
				for _, g := range goats.Items {
					reqs = append(reqs, reconcile.Request{
						NamespacedName: types.NamespacedName{
							Namespace: g.Namespace,
							Name:      g.Name,
						},
					})
				}
				return reqs
			}),
		).
		Complete(r)
}
