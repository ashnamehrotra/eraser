//go:build collector
// +build collector

package collector

import (
	"context"
	"testing"
	"time"

	eraserv1alpha1 "github.com/Azure/eraser/api/v1alpha1"
	"github.com/Azure/eraser/test/e2e/util"
	// appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	// clientgo "k8s.io/client-go/kubernetes"

	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestRemoveImagesFromAllNodes(t *testing.T) {
	const (
		nginx         = "nginx"
		nginxLatest   = "docker.io/library/nginx:latest"
		nginxAliasOne = "docker.io/library/nginx:one"
		nginxAliasTwo = "docker.io/library/nginx:two"
		redis         = "redis"
		caddy         = "caddy"

		prune               = "imagelist"
		skippedNodeName     = "eraser-e2e-test-worker"
		skippedNodeSelector = "kubernetes.io/hostname=eraser-e2e-test-worker"
		skipLabelKey        = "eraser.sh/cleanup.skip"
		skipLabelValue      = "true"
	)

	collectScanErasePipelineFeat := features.New("Test Remove Image From All Nodes").
		Assess("ImageCollector CR is generated", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			c, err := cfg.NewClient()
			if err != nil {
				t.Error("Failed to create new client", err)
			}

			resource := eraserv1alpha1.ImageCollector{}
			wait.For(func() (bool, error) {
				err := c.Resources().Get(ctx, "imagecollector-shared", "default", &resource)
				if err != nil {
					t.Logf("WE ARE HERE")
				}

				if resource.ObjectMeta.Name == "imagecollector-shared" {
					return true, nil
				}

				return false, nil
			}, wait.WithTimeout(time.Minute*3))

			return ctx
		}).
		Assess("ImageList CR is generated", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			c, err := cfg.NewClient()
			if err != nil {
				t.Error("Failed to create new client", err)
			}

			resource := eraserv1alpha1.ImageList{}
			wait.For(func() (bool, error) {
				err := c.Resources().Get(ctx, "imagelist", "default", &resource)
				if err != nil {
					t.Logf("WE ARE HERE")
				}

				if resource.ObjectMeta.Name == "imagelist" {
					return true, nil
				}

				return false, nil
			}, wait.WithTimeout(time.Minute*3))

			return ctx
		}).
		Assess("Images successfully deleted from all nodes", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			ctxT, cancel := context.WithTimeout(ctx, 3*time.Minute)
			defer cancel()
			util.CheckImageRemoved(ctxT, t, util.GetClusterNodes(t), nginx)

			return ctx
		}).
		Assess("Pods from imagejobs are cleaned up", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			c, err := cfg.NewClient()
			if err != nil {
				t.Error("Failed to create new client", err)
			}

			var ls corev1.PodList
			err = c.Resources().List(ctx, &ls, func(o *metav1.ListOptions) {
				o.LabelSelector = labels.SelectorFromSet(map[string]string{"name": "collector"}).String()
			})
			if err != nil {
				t.Errorf("could not list pods: %v", err)
			}

			err = wait.For(conditions.New(c.Resources()).ResourcesDeleted(&ls), wait.WithTimeout(time.Minute))
			if err != nil {
				t.Errorf("error waiting for pods to be deleted: %v", err)
			}

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			if err := util.DeleteImageListsAndJobs(cfg.KubeconfigFile()); err != nil {
				t.Error("Failed to clean eraser obejcts ", err)
			}
			return ctx
		}).
		Feature()

	testenv.Test(t, collectScanErasePipelineFeat)
}
