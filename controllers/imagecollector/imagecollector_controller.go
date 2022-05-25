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

package imagecollector

import (
	"context"
	"flag"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	eraserv1alpha1 "github.com/Azure/eraser/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	collectorImage = flag.String("collector-image", "ghcr.io/azure/collector:latest", "collector image")
	log            = logf.Log.WithName("controller").WithValues("process", "imagecollector-controller")
)

// ImageCollectorReconciler reconciles a ImageCollector object
type Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler.
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("imagecollector-controller", mgr, controller.Options{
		Reconciler: r,
	})
	if err != nil {
		return err
	}

	// Watch for changes to ImageJobs with ImageCollector owner (created by this collector controller)
	err = c.Watch(&source.Kind{Type: &eraserv1alpha1.ImageJob{}}, &handler.EnqueueRequestForOwner{OwnerType: &eraserv1alpha1.ImageCollector{}, IsController: true})
	if err != nil {
		return err
	}

	// GenericEvent in order to trigger first reconcile
	ch := make(chan event.GenericEvent)
	err = c.Watch(&source.Channel{
		Source: ch,
	}, &handler.EnqueueRequestForObject{})

	// go routine runs in background so that it does not block
	go func() {
		log.Info("Queueing first ImageCollector reconcile...")
		// goes to ImageCollector Reconcile
		ch <- event.GenericEvent{
			Object: &eraserv1alpha1.ImageCollector{},
		}
		log.Info("Queued first ImageCollector reconcile")
	}()

	return nil
}

//+kubebuilder:rbac:groups=eraser.sh.azure-eraser.io,resources=imagecollectors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=eraser.sh.azure-eraser.io,resources=imagecollectors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=eraser.sh.azure-eraser.io,resources=imagecollectors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ImageCollector object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Info("ImageCollector Reconcile")

	imageJob := &eraserv1alpha1.ImageJob{}
	if err := r.Get(ctx, req.NamespacedName, imageJob); err != nil {
		if errors.IsNotFound(err) {
			if res, err := r.createImageJob(ctx); err != nil {
				return res, err
			}
		} else {
			return reconcile.Result{}, err
		}
	} else {
		if imageJob.Status.Phase == eraserv1alpha1.PhaseCompleted {
			log.Info("completed phase")
			if res, err := r.createSharedCRD(ctx); err != nil {
				return res, err
			}
		}
	}

	// else there is a job in running phase, requeue
	return ctrl.Result{RequeueAfter: time.Second * 30}, nil
}

func (r *Reconciler) createImageJob(ctx context.Context) (ctrl.Result, error) {
	job := &eraserv1alpha1.ImageJob{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "imagejob-",
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(&eraserv1alpha1.ImageCollector{
					TypeMeta:   metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{Name: "test", UID: types.UID("test")},
				}, schema.GroupVersionKind{Group: "eraser.sh", Version: "v1alpha1", Kind: "ImageCollector"}),
			},
		},
		Spec: eraserv1alpha1.ImageJobSpec{
			JobTemplate: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:            "collector",
							Image:           *collectorImage,
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
					ServiceAccountName: "eraser-controller-manager",
				},
			},
		},
	}

	err := r.Create(ctx, job)
	if err != nil {
		log.Info("Could not create collector ImageJob")
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	log.Info("Successfully created collector ImageJob", "job", job.Name)
	return reconcile.Result{}, nil
}

func (r *Reconciler) createSharedCRD(ctx context.Context) (ctrl.Result, error) {
	// implement once working
	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eraserv1alpha1.ImageCollector{}).
		Complete(r)
}
