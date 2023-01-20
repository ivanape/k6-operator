/*
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
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/go-logr/logr"
	"github.com/grafana/k6-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// K6Reconciler reconciles a K6 object
type K6Reconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile takes a K6 object and takes the appropriate action in the cluster
// +kubebuilder:rbac:groups=k6.io,resources=k6s,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=k6.io,resources=k6s/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
func (r *K6Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("namespace", req.Namespace, "name", req.Name)

	// Fetch the CRD
	k6 := &v1alpha1.K6{}
	err := r.Get(ctx, req.NamespacedName, k6)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Request deleted. Nothing to reconcile.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Could not fetch request")
		return ctrl.Result{Requeue: true}, err
	}

	log.Info(fmt.Sprintf("Reconcile(); stage = %s", k6.Status.Stage))

	switch k6.Status.Stage {
	case "", "initialization":
		return InitializeJobs(ctx, log, k6, r)
	case "initialized":
		return CreateJobs(ctx, log, k6, r)
	case "created":
		return StartJobs(ctx, log, k6, r)
	case "started":
		// wait for test to finish and then mark as finished
		return FinishJobs(ctx, log, k6, r)
	case "error", "finished":
		// delete if configured
		if k6.Spec.Cleanup == "post" {
			log.Info("Cleaning up all resources")
			r.Delete(ctx, k6)
		}
		// notify if configured
		return ctrl.Result{}, nil
	}

	err = fmt.Errorf("invalid status")
	log.Error(err, "Invalid status for the k6 resource.")
	return ctrl.Result{}, err
}

// SetupWithManager sets up a managed controller that will reconcile all events for the K6 CRD
func (r *K6Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.K6{}).
		Owns(&batchv1.Job{}).
		Watches(&source.Kind{Type: &v1.Pod{}},
			&handler.EnqueueRequestsFromMapFunc{ToRequests: &K6PodsWatchMap{log: mgr.GetLogger()}},
			builder.WithPredicates(filterK6Pods())).
		Complete(r)
}
