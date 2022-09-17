package metrics

import (
	"fmt"
	"net/http"
	"os"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"k8s.io/klog/v2"
)

func InitPrometheusExporter(metricsAddr string) error {
	config := prometheus.Config{}

	ctrl := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				// change distribution
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)

	exporter, err := prometheus.New(config, ctrl)
	if err != nil {
		return fmt.Errorf("failed to register prometheus exporter: %v", err)
	}

	http.HandleFunc("/metrics", exporter.ServeHTTP)
	go func() {
		if err := http.ListenAndServe(metricsAddr, nil); err != nil {
			klog.ErrorS(err, "failed to register prometheus endpoint", "metricsAddress", metricsAddr)
			os.Exit(1)
		}
	}()

	global.SetMeterProvider(exporter.MeterProvider())

	return nil
}
