//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/Azure/eraser/test/e2e/util"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"regexp"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"strconv"
)

const (
	ExpectedImagesRemoved = 3
)

func TestMetrics(t *testing.T) {
	metrics := features.New("Images_removed_run_total metric should report 1").
		Assess("Alpine image is removed", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			// deploy imagelist config
			if err := util.DeployEraserConfig(cfg.KubeconfigFile(), util.TestNamespace, "../../test-data", "imagelist_alpine.yaml"); err != nil {
				t.Error("Failed to deploy image list config", err)
			}

			ctxT, cancel := context.WithTimeout(ctx, util.Timeout)
			defer cancel()
			util.CheckImageRemoved(ctxT, t, util.GetClusterNodes(t), util.VulnerableImage)

			return ctx
		}).
		Assess("Check images_removed_run_total metric", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			c, err := cfg.NewClient()
			if err != nil {
				t.Fatal("Failed to create new client", err)
			}

			var ls corev1.PodList
			err = c.Resources().List(ctx, &ls, func(o *metav1.ListOptions) {
				o.LabelSelector = labels.SelectorFromSet(map[string]string{"component": "otel-collector"}).String()
			})
			if err != nil {
				t.Errorf("could not list pods: %v", err)
			}

			otelcollector := ls.Items[0]

			if err := util.KubectlPortForward(cfg.KubeconfigFile(), otelcollector.Name, util.TestNamespace); err != nil {
				t.Error(err, "error in kubectl port-forward otel-collector")
			}

			if _, err := util.KubectlCurlPod(cfg.KubeconfigFile()); err != nil {
				t.Error(err, "error running curl pod")
			}

			if _, err := util.KubectlWait(cfg.KubeconfigFile(), "temp"); err != nil {
				t.Error(err, "error waiting for temp curl pod")
			}

			service, err := util.KubectlDescribeService(cfg.KubeconfigFile(), "otel-collector", util.TestNamespace)
			if err != nil {
				t.Error(err, "could not get otel collector service")
			}

			regex := regexp.MustCompile(`IP:\s+(\d+.\d+.\d+.\d+)`)
			match := regex.FindStringSubmatch(service)

			otelEndpoint := "http://" + match[1] + ":8889/metrics"

			output, err := util.KubectlExecCurl(cfg.KubeconfigFile(), "temp", otelEndpoint)
			if err != nil {
				t.Error(err, "error with otlp curl request")
			}

			t.Log("OUTPUT ", output)

			r := regexp.MustCompile(`images_removed_run_total{job="eraser",node_name=".+"} (\d+)`)
			results := r.FindAllStringSubmatch(output, -1)

			totalRemoved := 0
			for i, _ := range results {
				val, _ := strconv.Atoi(results[i][1])
				totalRemoved += val
			}

			if totalRemoved != 3 {
				t.Error("images_removed_run_total incorrect, expected 3, got", totalRemoved)
			}

			return ctx
		}).
		Feature()

	util.Testenv.Test(t, metrics)
}
