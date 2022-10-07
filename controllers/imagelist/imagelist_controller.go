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

package imagelist

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	eraserv1alpha1 "github.com/Azure/eraser/api/v1alpha1"
	"github.com/Azure/eraser/controllers/util"
	"github.com/Azure/eraser/pkg/logger"
	"github.com/Azure/eraser/pkg/utils"
	"go.opentelemetry.io/otel/metric/unit"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

const (
	imgListPath = "/run/eraser.sh/imagelist"
)

var (
	log       = logf.Log.WithName("controller").WithValues("process", "imagelist-controller")
	imageList = types.NamespacedName{Name: "imagelist"}
	startTime time.Time
)

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler.
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &Reconciler{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
	}
}

// ImageJobReconciler reconciles a ImageJob object.
type ImageJobReconciler struct {
	client.Client
}

// ImageListReconciler reconciles a ImageList object.
type Reconciler struct {
	client.Client
	scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=eraser.sh,resources=imagelists,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=eraser.sh,resources=imagelists/status,verbs=get;update;patch
//+kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;update;create;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ImageList object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Ignore unsupported lists
	if req.NamespacedName != imageList {
		log.Info("Ignoring unsupported imagelist name", "name", req.Name)
		return reconcile.Result{}, nil
	}

	imageList := eraserv1alpha1.ImageList{}
	err := r.Get(ctx, req.NamespacedName, &imageList)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	jobList := eraserv1alpha1.ImageJobList{}
	err = r.List(ctx, &jobList)
	if client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	}

	items := util.FilterJobListByOwner(jobList.Items, metav1.NewControllerRef(&imageList, imageList.GroupVersionKind()))

	switch len(items) {
	case 0:
		return r.handleImageListEvent(ctx, &req, &imageList)
	case 1:
		job := items[0]

		// If we got here because of a completed ImageJob:
		if util.IsCompletedOrFailed(job.Status.Phase) {
			return r.handleJobListEvent(ctx, &imageList, &job)
		}

		// If we got here due to an update to the ImageList, and there is an ImageJob already running,
		// keep requeueing it until that job is completed.
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	default:
		return ctrl.Result{}, fmt.Errorf("there are multiple child imagejobs running")
	}
}

func (r *Reconciler) handleJobListEvent(ctx context.Context, imageList *eraserv1alpha1.ImageList, job *eraserv1alpha1.ImageJob) (ctrl.Result, error) {
	phase := job.Status.Phase
	if phase == eraserv1alpha1.PhaseCompleted || phase == eraserv1alpha1.PhaseFailed {
		err := r.handleJobCompletion(ctx, imageList, job)
		if err != nil {
			return ctrl.Result{}, err
		}

		if job.Status.DeleteAfter == nil {
			if job.Status.Phase == eraserv1alpha1.PhaseCompleted {
				job.Status.DeleteAfter = util.After(time.Now(), int64(util.SuccessDel.Seconds()))
			} else if job.Status.Phase == eraserv1alpha1.PhaseFailed {
				job.Status.DeleteAfter = util.After(time.Now(), int64(util.ErrDel.Seconds()))
			}

			if err := r.Status().Update(ctx, job); err != nil {
				log.Info("Could not update Delete After for job " + job.Name)
			}
			return ctrl.Result{}, nil
		}

		// record metrics
		ctxB := context.Background()
		ctx, cancel := signal.NotifyContext(ctxB, os.Interrupt, syscall.SIGTERM)
		defer cancel()

		exporter, err := otlpmetrichttp.New(ctx, otlpmetrichttp.WithInsecure(), os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
		if err != nil {
			panic(err)
		}

		reader := sdkmetric.NewPeriodicReader(exporter)
		provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
		defer func() {
			fmt.Fprintln(os.Stderr, "collecting final metrics...")
			m, err := reader.Collect(ctxB)
			if err != nil {
				fmt.Fprintln(os.Stderr, "failed to collect metrics:", err)
				return
			}
			if err := exporter.Export(ctxB, m); err != nil {
				fmt.Fprintln(os.Stderr, "failed to export metrics:", err)
			}
			if err := provider.Shutdown(ctxB); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}()
		global.SetMeterProvider(provider)
		recordMetrics(ctx, float64(time.Since(startTime).Milliseconds()), int64(job.Status.Succeeded), int64(job.Status.Failed))

		return r.handleJobDeletion(ctx, job)
	}

	return ctrl.Result{}, fmt.Errorf("unexpected job phase: '%s'", job.Status.Phase)
}

func (r *Reconciler) handleJobDeletion(ctx context.Context, job *eraserv1alpha1.ImageJob) (ctrl.Result, error) {
	until := time.Until(job.Status.DeleteAfter.Time)
	if until > 0 {
		log.Info("Delaying imagejob delete", "job", job.Name, "deleteAter", job.Status.DeleteAfter)
		return ctrl.Result{RequeueAfter: until}, nil
	}

	log.Info("Deleting imagejob", "job", job.Name)
	err := r.Delete(ctx, job)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler) handleImageListEvent(ctx context.Context, req *ctrl.Request, imageList *eraserv1alpha1.ImageList) (ctrl.Result, error) {
	imgListJSON, err := json.Marshal(imageList.Spec.Images)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("marshal image list: %w", err)
	}

	configMap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "imagelist-",
			Namespace:    utils.GetNamespace(),
		},
		Immutable: utils.BoolPtr(true),
		Data:      map[string]string{"images": string(imgListJSON)},
	}
	if err := r.Create(ctx, &configMap); err != nil {
		return ctrl.Result{}, fmt.Errorf("create configmap: %w", err)
	}

	configName := configMap.Name
	args := []string{
		"--imagelist=" + filepath.Join(imgListPath, "images"),
		"--log-level=" + logger.GetLevel(),
	}
	args = append(args, util.EraserArgs...)

	job := &eraserv1alpha1.ImageJob{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "imagejob-",
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(imageList, imageList.GroupVersionKind()),
			},
		},
		Spec: eraserv1alpha1.ImageJobSpec{
			JobTemplate: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: configName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: configName}},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:            "eraser",
							Image:           *util.EraserImage,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Args:            args,
							VolumeMounts: []corev1.VolumeMount{
								{MountPath: imgListPath, Name: configName},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse("7m"),
									"memory": resource.MustParse("25Mi"),
								},
								Limits: corev1.ResourceList{
									"cpu":    resource.MustParse("8m"),
									"memory": resource.MustParse("30Mi"),
								},
							},
							SecurityContext: utils.SharedSecurityContext,
							Env:             []corev1.EnvVar{{Name: "OTEL_EXPORTER_OTLP_ENDPOINT", Value: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")}, {Name: "NODE_NAME", ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "spec.nodeName"}}}},
						},
					},
					ServiceAccountName: "eraser-imagejob-pods",
				},
			},
		},
	}

	configmapList := &corev1.ConfigMapList{}
	if err := r.List(ctx, configmapList); err != nil {
		log.Info("Could not get list of configmaps")
		return reconcile.Result{}, err
	}

	exclusionMount, exclusionVolume, err := util.GetExclusionVolume(configmapList)
	if err != nil {
		log.Info("Could not get exclusion mounts and volumes")
		return reconcile.Result{}, err
	}

	for i := range job.Spec.JobTemplate.Spec.Containers {
		job.Spec.JobTemplate.Spec.Containers[i].VolumeMounts = append(job.Spec.JobTemplate.Spec.Containers[i].VolumeMounts, exclusionMount...)
	}

	job.Spec.JobTemplate.Spec.Volumes = append(job.Spec.JobTemplate.Spec.Volumes, exclusionVolume...)

	err = r.Create(ctx, job)
	startTime = time.Now()
	log.Info("creating imagejob", "job", job.Name)

	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	configMap.ObjectMeta.OwnerReferences = []metav1.OwnerReference{*metav1.NewControllerRef(job, schema.GroupVersionKind{
		Group:   "eraser.sh",
		Version: "v1alpha1",
		Kind:    "ImageJob",
	})}
	err = r.Update(ctx, &configMap)
	if err != nil {
		return reconcile.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler) handleJobCompletion(ctx context.Context, imageList *eraserv1alpha1.ImageList, job *eraserv1alpha1.ImageJob) error {
	now := metav1.Now()

	imageList.Status.Success = int64(job.Status.Succeeded)
	imageList.Status.Failed = int64(job.Status.Failed)
	imageList.Status.Skipped = int64(job.Status.Skipped)
	imageList.Status.Timestamp = &now

	err := r.Status().Update(ctx, imageList)
	if err != nil {
		return err
	}

	return nil
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("imagelist-controller", mgr, controller.Options{
		Reconciler: r,
	})
	if err != nil {
		return err
	}

	err = c.Watch(
		&source.Kind{Type: &eraserv1alpha1.ImageList{}},
		&handler.EnqueueRequestForObject{}, predicate.GenerationChangedPredicate{})
	if err != nil {
		return err
	}
	err = c.Watch(
		&source.Kind{Type: &eraserv1alpha1.ImageJob{}},
		&handler.EnqueueRequestForOwner{OwnerType: &eraserv1alpha1.ImageList{}, IsController: true},
		predicate.Funcs{
			// Do nothing on Create, Delete, or Generic events
			CreateFunc:  util.NeverOnCreate,
			DeleteFunc:  util.NeverOnDelete,
			GenericFunc: util.NeverOnGeneric,
			UpdateFunc: func(e event.UpdateEvent) bool {
				if job, ok := e.ObjectNew.(*eraserv1alpha1.ImageJob); ok && util.IsCompletedOrFailed(job.Status.Phase) {
					return true
				}

				return false
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func recordMetrics(ctx context.Context, jobDuration float64, podsCompleted int64, podsFailed int64) {
	p := global.MeterProvider()

	duration, err := p.Meter("eraser").SyncFloat64().Histogram("ImageJobEraserDuration", instrument.WithDescription("duration of eraser imagejob"), instrument.WithUnit(unit.Milliseconds))
	if err != nil {
		panic(err)
	}
	duration.Record(ctx, jobDuration)

	completed, err := p.Meter("eraser").SyncInt64().Counter("PodsCompleted", instrument.WithDescription("total pods completed"), instrument.WithUnit("1"))
	if err != nil {
		panic(err)
	}
	completed.Add(ctx, podsCompleted)

	failed, err := p.Meter("eraser").SyncInt64().Counter("PodsFailed", instrument.WithDescription("total pods failed"), instrument.WithUnit("1"))
	if err != nil {
		panic(err)
	}
	failed.Add(ctx, podsFailed)

	jobTotal, err := p.Meter("eraser").SyncInt64().Counter("ImageJobEraserTotal", instrument.WithDescription("total number of eraser imagejobs completed"), instrument.WithUnit("1"))
	if err != nil {
		panic(err)
	}
	jobTotal.Add(ctx, 1)
}
