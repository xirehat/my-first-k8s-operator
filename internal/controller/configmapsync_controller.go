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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "github.com/xirehat/my-first-k8s-operator/api/v1"
)

// ConfigMapSyncReconciler reconciles a ConfigMapSync object
type ConfigMapSyncReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.xirehat.com,resources=configmapsyncs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.xirehat.com,resources=configmapsyncs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.xirehat.com,resources=configmapsyncs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConfigMapSync object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *ConfigMapSyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

    ctx := context.Background()
    log := r.Log.WithValues("configmapsync", req.NamespacedName)
// Fetch the ConfigMapSync instance
    configMapSync := &appsv1.ConfigMapSync{}
    if err := r.Get(ctx, req.NamespacedName, configMapSync); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }
    // Fetch the source ConfigMap
    sourceConfigMap := &corev1.ConfigMap{}
    sourceConfigMapName := types.NamespacedName{
        Namespace: configMapSync.Spec.SourceNamespace,
        Name:      configMapSync.Spec.ConfigMapName,
    }
    if err := r.Get(ctx, sourceConfigMapName, sourceConfigMap); err != nil {
        return ctrl.Result{}, err
    }
    // Create or Update the destination ConfigMap in the target namespace
    destinationConfigMap := &corev1.ConfigMap{}
    destinationConfigMapName := types.NamespacedName{
        Namespace: configMapSync.Spec.DestinationNamespace,
        Name:      configMapSync.Spec.ConfigMapName,
    }
    if err := r.Get(ctx, destinationConfigMapName, destinationConfigMap); err != nil {
        if errors.IsNotFound(err) {
            log.Info("Creating ConfigMap in destination namespace", "Namespace", configMapSync.Spec.DestinationNamespace)
            destinationConfigMap = &corev1.ConfigMap{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      configMapSync.Spec.ConfigMapName,
                    Namespace: configMapSync.Spec.DestinationNamespace,
                },
                Data: sourceConfigMap.Data, // Copy data from source to destination
            }
            if err := r.Create(ctx, destinationConfigMap); err != nil {
                return ctrl.Result{}, err
            }
        } else {
            return ctrl.Result{}, err
        }
    } else {
        log.Info("Updating ConfigMap in destination namespace", "Namespace", configMapSync.Spec.DestinationNamespace)
        destinationConfigMap.Data = sourceConfigMap.Data // Update data from source to destination
        if err := r.Update(ctx, destinationConfigMap); err != nil {
            return ctrl.Result{}, err
        }
    }
	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigMapSyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.ConfigMapSync{}).
		Named("configmapsync").
		Complete(r)
}
