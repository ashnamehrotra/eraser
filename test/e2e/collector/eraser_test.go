//go:build collector
// +build collector

package collector

import (
	"context"
	"testing"
	"time"

	eraserv1alpha1 "github.com/Azure/eraser/api/v1alpha1"
	"github.com/Azure/eraser/test/e2e/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"os"
	"path/filepath"
)

func TestRemoveImagesFromAllNodes(t *testing.T) {
	const (
		alpine = "alpine"
		nginx  = "nginx"
	)

	collectScanErasePipelineFeat := features.New("Test Remove Image From All Nodes").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			wd, err := os.Getwd()
			if err != nil {
				t.Error("Could not get wd")
			}

			providerResourceAbsolutePath, err := filepath.Abs(filepath.Join(wd, "/../../../", providerResourceDirectory, "eraser"))
			if err != nil {
				t.Error("Could not get provider resource absolute pathy")
			}
			// start deployment
			if err := util.HelmInstall(cfg.KubeconfigFile(), "eraser-system", []string{providerResourceAbsolutePath}); err != nil {
				t.Error("Unable to helm install deployment", err)
			}

			return ctx
		}).
		Assess("ImageCollector CR is generated", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			c, err := cfg.NewClient()
			if err != nil {
				t.Fatal("Failed to create new client", err)
			}

			resource := eraserv1alpha1.ImageCollector{}
			wait.For(func() (bool, error) {
				err := c.Resources().Get(ctx, "imagecollector-shared", "default", &resource)
				if err != nil {
					return false, err
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
				t.Fatal("Failed to create new client", err)
			}

			resource := eraserv1alpha1.ImageList{}
			wait.For(func() (bool, error) {
				err := c.Resources().Get(ctx, "imagelist", "default", &resource)
				if util.IsNotFound(err) {
					return false, nil
				}

				if err != nil {
					return false, err
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
			util.CheckImageRemoved(ctxT, t, util.GetClusterNodes(t), alpine)

			return ctx
		}).
		Assess("Pods from imagejobs are cleaned up", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			c, err := cfg.NewClient()
			if err != nil {
				t.Fatal("Failed to create new client", err)
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
			c, err := cfg.NewClient()
			if err != nil {
				t.Fatal("Failed to create new client", err)
			}

			err = util.HelmUninstall(cfg.KubeconfigFile(), "eraser-system", []string{})
			if err != nil {
				t.Error("Unable to uninstall deployment for teardown", err)
			}

			var ls corev1.PodList
			err = c.Resources().List(ctx, &ls, func(o *metav1.ListOptions) {
				o.LabelSelector = labels.SelectorFromSet(map[string]string{"name": "eraser-manager"}).String()
			})
			if err != nil {
				t.Error("could not list eraser manager pod")
			}

			err = wait.For(conditions.New(c.Resources()).ResourcesDeleted(&ls), wait.WithTimeout(time.Minute))
			if err != nil {
				t.Errorf("error waiting for eraser-manager to be deleted: %v", err)
			}

			return ctx
		}).
		Feature()

	disableScanFeat := features.New("Test Scanner Disabled Prune").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			wd, err := os.Getwd()
			if err != nil {
				t.Error("Could not get working directory", err)
			}

			providerResourceAbsolutePath, err := filepath.Abs(filepath.Join(wd, "/../../../", providerResourceDirectory, "eraser"))
			if err != nil {
				t.Error("Unable to get provider resource absolute path", err)
			}

			err = util.HelmInstall(cfg.KubeconfigFile(), "eraser-system", []string{providerResourceAbsolutePath, "--set", "scanner.image.repository="})
			if err != nil {
				t.Error("Unable to install deployment with scanner disabled", err)
			}

			return ctx
		}).
		Assess("ImageCollector CR is generated", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			c, err := cfg.NewClient()
			if err != nil {
				t.Error("Failed to create new client", err)
			}

			imagecollector := eraserv1alpha1.ImageCollector{}
			wait.For(func() (bool, error) {
				err := c.Resources().Get(ctx, "imagecollector-shared", "default", &imagecollector)
				if err != nil {
					t.Error("Could not get imagecollector-shared")
				}

				if imagecollector.ObjectMeta.Name == "imagecollector-shared" {
					return true, nil
				}

				return false, nil
			}, wait.WithTimeout(time.Minute*3))

			return ctx
		}).
		Assess("ImageList Spec Contains Same Images As ImageCollecotor Shared", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			c, err := cfg.NewClient()
			if err != nil {
				t.Error("Failed to create new client", err)
			}

			// verify imagelist created
			imagelist := eraserv1alpha1.ImageList{}
			wait.For(func() (bool, error) {
				err := c.Resources().Get(ctx, "imagelist", "default", &imagelist)
				if util.IsNotFound(err) {
					return false, nil
				}

				if err != nil {
					return false, err
				}

				if imagelist.ObjectMeta.Name == "imagelist" {
					return true, nil
				}

				return false, nil
			}, wait.WithTimeout(time.Minute*3))

			imagecollectorShared := eraserv1alpha1.ImageCollector{}
			err = c.Resources().Get(ctx, "imagecollector-shared", "default", &imagecollectorShared)
			if err != nil {
				t.Error("Could not get imagecollector-shared")
			}

			// verify imagecollector-shared status fields are empty
			if imagecollectorShared.Status.Vulnerable != nil || imagecollectorShared.Status.Failed != nil {
				t.Error("Scan job has run, should be disabled")
			}

			imagelistSpec := make(map[string]struct{}, len(imagelist.Spec.Images))
			for _, img := range imagelist.Spec.Images {
				imagelistSpec[img] = struct{}{}
			}

			// verify the images in both lists match
			for _, img := range imagecollectorShared.Spec.Images {
				// check by digest as we add to imagelist by digest when pruning without scanner
				if _, contains := imagelistSpec[img.Digest]; !contains {
					t.Error("imagelist spec does not match imagecollector-shared: ", img.Digest)
				}
			}

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			c, err := cfg.NewClient()
			if err != nil {
				t.Fatal("Failed to create new client", err)
			}

			err = util.HelmUninstall(cfg.KubeconfigFile(), "eraser-system", []string{})
			if err != nil {
				t.Error("Unable to uninstall deployment for teardown", err)
			}

			var ls corev1.PodList
			err = c.Resources().List(ctx, &ls, func(o *metav1.ListOptions) {
				o.LabelSelector = labels.SelectorFromSet(map[string]string{"name": "eraser-manager"}).String()
			})
			if err != nil {
				t.Error("could not list eraser manager pod")
			}

			err = wait.For(conditions.New(c.Resources()).ResourcesDeleted(&ls), wait.WithTimeout(time.Minute))
			if err != nil {
				t.Errorf("error waiting for eraser-manager to be deleted: %v", err)
			}

			return ctx
		}).
		Feature()

	excludedFeat := features.New("Test Excluded").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			wd, err := os.Getwd()
			if err != nil {
				t.Error("Could not get working directory", err)
			}

			providerResourceAbsolutePath, err := filepath.Abs(filepath.Join(wd, "/../../../", providerResourceDirectory, "eraser"))
			if err != nil {
				t.Error("Unable to get provider resource absolute path", err)
			}

			err = util.HelmInstall(cfg.KubeconfigFile(), "eraser-system", []string{providerResourceAbsolutePath, "--set", "collector.image.repository="})
			if err != nil {
				t.Error("Unable to install deployment with no collector disabled", err)
			}

			podSelectorLabels := map[string]string{"app": nginx}
			nginxDep := util.NewDeployment(cfg.Namespace(), nginx, 2, podSelectorLabels, corev1.Container{Image: nginx, Name: nginx})
			if err := cfg.Client().Resources().Create(ctx, nginxDep); err != nil {
				t.Error("Failed to create the deployment", err)
			}

			// create excluded configmap and add docker.io/library/*
			excluded := corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "excluded",
					Namespace: "eraser-system",
				},
				Data: map[string]string{"excluded": "docker.io/library/*"},
			}
			if err := cfg.Client().Resources().Create(ctx, &excluded); err != nil {
				t.Error("failed to create excluded configmap", err)
			}

			return ctx
		}).
		Assess("deployment successfully deployed", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client, err := cfg.NewClient()
			if err != nil {
				t.Error("Failed to create new client", err)
			}

			resultDeployment := appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{Name: nginx, Namespace: cfg.Namespace()},
			}

			if err = wait.For(conditions.New(client.Resources()).DeploymentConditionMatch(&resultDeployment, appsv1.DeploymentAvailable, corev1.ConditionTrue),
				wait.WithTimeout(time.Minute*3)); err != nil {
				t.Error("deployment not found", err)
			}

			return context.WithValue(ctx, nginx, &resultDeployment)
		}).
		Assess("Check image remains in all nodes", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			// delete deployment
			client, err := cfg.NewClient()
			if err != nil {
				t.Error("Failed to create new client", err)
			}

			var pods corev1.PodList
			err = client.Resources().List(ctx, &pods, func(o *metav1.ListOptions) {
				o.LabelSelector = labels.SelectorFromSet(labels.Set{"app": nginx}).String()
			})
			if err != nil {
				t.Fatal(err)
			}

			dep := ctx.Value(nginx).(*appsv1.Deployment)
			if err := client.Resources().Delete(ctx, dep); err != nil {
				t.Error("Failed to delete the dep", err)
			}

			for _, nodeName := range util.GetClusterNodes(t) {
				err := wait.For(util.ContainerNotPresentOnNode(nodeName, nginx), wait.WithTimeout(time.Minute*2))
				if err != nil {
					t.Logf("error while waiting for deployment deletion: %v", err)
				}
			}

			// create imagelist to trigger deletion
			if err := util.DeployEraserConfig(cfg.KubeconfigFile(), "eraser-system", "../test-data", "eraser_v1alpha1_imagelist.yaml"); err != nil {
				t.Error("Failed to deploy image list config", err)
			}

			ctxT, cancel := context.WithTimeout(ctx, time.Minute)
			defer cancel()
			// since docker.io/library/* was excluded, nginx should still exist following deletion
			util.CheckImagesExist(ctxT, t, util.GetClusterNodes(t), nginx)

			return ctx
		}).
		Assess("Pods from imagejobs are cleaned up", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			c, err := cfg.NewClient()
			if err != nil {
				t.Error("Failed to create new client", err)
			}

			var ls corev1.PodList
			err = c.Resources().List(ctx, &ls, func(o *metav1.ListOptions) {
				o.LabelSelector = labels.SelectorFromSet(map[string]string{"name": "eraser"}).String()
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
			c, err := cfg.NewClient()
			if err != nil {
				t.Fatal("Failed to create new client", err)
			}

			err = util.HelmUninstall(cfg.KubeconfigFile(), "eraser-system", []string{})
			if err != nil {
				t.Error("Unable to uninstall deployment for teardown", err)
			}

			var ls corev1.PodList
			err = c.Resources().List(ctx, &ls, func(o *metav1.ListOptions) {
				o.LabelSelector = labels.SelectorFromSet(map[string]string{"name": "eraser-controller"}).String()
			})
			if err != nil {
				t.Error("could not list eraser manager pod")
			}

			err = wait.For(conditions.New(c.Resources()).ResourcesDeleted(&ls), wait.WithTimeout(time.Minute))
			if err != nil {
				t.Errorf("error waiting for eraser-manager to be deleted: %v", err)
			}
			return ctx
		}).
		Feature()

	testenv.Test(t, disableScanFeat)
	testenv.Test(t, collectScanErasePipelineFeat)
	testenv.Test(t, excludedFeat)
}
