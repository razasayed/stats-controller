/*
Copyright 2024.

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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	monitoringv1 "github.com/razasayed/stats-controller/api/v1"
)

// StatsReconciler reconciles a Stats object
type StatsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=monitoring.io,resources=stats,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.io,resources=stats/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=monitoring.io,resources=stats/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Stats object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
func (r *StatsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	log.FromContext(ctx).Info("Reconciling Stats...")
	// Fetch the Stats instance
	stats := &monitoringv1.Stats{}
	if err := r.Get(ctx, req.NamespacedName, stats); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Count the number of running Pods
	podList := &corev1.PodList{}
	if err := r.List(ctx, podList); err != nil {
		log.FromContext(ctx).Error(err, "Failed to list pods")
		return ctrl.Result{}, err
	}
	runningPods := 0
	for _, pod := range podList.Items {
		if pod.Status.Phase == corev1.PodRunning {
			runningPods++
		}
	}

	// Count the number of Deployments
	deploymentList := &appsv1.DeploymentList{}
	if err := r.List(ctx, deploymentList); err != nil {
		log.FromContext(ctx).Error(err, "Failed to list deployments")
		return ctrl.Result{}, err
	}
	deployments := len(deploymentList.Items)
	deploymentNames := []string{}
	for _, deployment := range deploymentList.Items {
		deploymentNames = append(deploymentNames, deployment.Name)
	}

	// Count the number of DaemonSets
	daemonSetList := &appsv1.DaemonSetList{}
	if err := r.List(ctx, daemonSetList); err != nil {
		log.FromContext(ctx).Error(err, "Failed to list daemonsets")
		return ctrl.Result{}, err
	}
	daemonSets := len(daemonSetList.Items)
	daemonSetNames := []string{}
	for _, daemonSet := range daemonSetList.Items {
		daemonSetNames = append(daemonSetNames, daemonSet.Name)
	}

	// Count the number of StatefulSets
	statefulSetList := &appsv1.StatefulSetList{}
	if err := r.List(ctx, statefulSetList); err != nil {
		log.FromContext(ctx).Error(err, "Failed to list statefulsets")
		return ctrl.Result{}, err
	}
	statefulSets := len(statefulSetList.Items)
	statefulSetNames := []string{}
	for _, statefulSet := range statefulSetList.Items {
		statefulSetNames = append(statefulSetNames, statefulSet.Name)
	}

	// Count the number of ReplicaSets
	replicaSetList := &appsv1.ReplicaSetList{}
	if err := r.List(ctx, replicaSetList); err != nil {
		log.FromContext(ctx).Error(err, "Failed to list replicasets")
		return ctrl.Result{}, err
	}
	replicaSets := len(replicaSetList.Items)
	replicaSetNames := []string{}
	for _, replicaSet := range replicaSetList.Items {
		replicaSetNames = append(replicaSetNames, replicaSet.Name)
	}

	// Update the status
	stats.Status.RunningPods = runningPods
	stats.Status.Deployments = deployments
	stats.Status.DeploymentNames = deploymentNames
	stats.Status.DaemonSets = daemonSets
	stats.Status.DaemonSetNames = daemonSetNames
	stats.Status.StatefulSets = statefulSets
	stats.Status.StatefulSetNames = statefulSetNames
	stats.Status.ReplicaSets = replicaSets
	stats.Status.ReplicaSetNames = replicaSetNames

	if err := r.Status().Update(ctx, stats); err != nil {
		log.FromContext(ctx).Error(err, "Failed to update Stats status")
		return ctrl.Result{}, err
	}
	log.FromContext(ctx).Info("Reconciling Stats Done !")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StatsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1.Stats{}).
		Complete(r)
}
