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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	tfv1alpha1 "github.com/pablochacin/tf-operator/api/v1alpha1"
    "github.com/pablochacin/tf-operator/pkg/jobs"
)

// StackReconciler reconciles a Stack object
type StackReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=tf.tf-operator.io,resources=stacks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tf.tf-operator.io,resources=stacks/status,verbs=get;update;patch

func (r *StackReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("stack", req.NamespacedName)

    var stack = tfv1alpha1.Stack{}
    err := r.Get(ctx, req.NamespacedName, &stack)
    if err != nil {
        log.Error(err, "unable to fetch Stack")
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // is stack been deleted?
	if !stack.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, stack)
	}

    return r.reconcileUpdate(ctx, stack)

	return ctrl.Result{}, nil
}

func (r *StackReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tfv1alpha1.Stack{}).
		Complete(r)
}

// reconcileDelete handles Stack delete
func (r *StackReconciler)reconcileDelete(ctx context.Context, stack tfv1alpha1.Stack) (ctrl.Result, error) {
    return ctrl.Result{}, nil
}

// reconcileUpdate handles stack creation and updates
func (r *StackReconciler)reconcileUpdate(ctx context.Context, stack tfv1alpha1.Stack) (ctrl.Result, error) {
   jobCfg := &jobs.JobConfig{
        Command:    "apply",
        Namespace:  stack.Namespace,
        Stack:      stack.Name,
        TfConfig:   stack.Spec.TfConfig.Name,
        Tfvars:     stack.Spec.TfVars.Name,
        Tfstate:    stack.Status.TfState.Name,
    }

    job, err := jobs.BuildJob(jobCfg)
    if err != nil {
        return ctrl.Result{}, err
    }

    err = r.Create(ctx, job)
    if err != nil {
        return ctrl.Result{}, err
    }

    return ctrl.Result{}, nil
}
