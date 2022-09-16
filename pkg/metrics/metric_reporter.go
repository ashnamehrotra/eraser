package metrics

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
	"k8s.io/klog/v2"
)

var (
	imageJobCollectorDuration syncfloat64.Histogram
	imageJobEraserDuration    syncfloat64.Histogram
	imagesRemoved             syncint64.Counter
	nonCompliantImages        syncint64.Counter
	imageJobCollectorTotal    syncint64.Counter
	imageJobEraserTotal       syncint64.Counter
	podsCompleted             syncint64.Counter
	podsFailed                syncint64.Counter
	metricsAddr               = ":8088"
	meter                     metric.Meter
)

func InitMetricInstruments() error {
	err, exporter := InitPrometheusExporter(metricsAddr)
	if err != nil {
		fmt.Println("unable to initialize prometheus exporter")
		return err
	}

	http.HandleFunc("/metrics", exporter.ServeHTTP)
	go func() {
		if err := http.ListenAndServe(metricsAddr, nil); err != nil {
			klog.ErrorS(err, "failed to register prometheus endpoint", "metricsAddress", metricsAddr)
			os.Exit(1)
		}
	}()

	klog.InfoS("Prometheus metrics server running", "address", metricsAddr)

	meter = global.MeterProvider().Meter("eraser")

	imageJobCollectorDuration, err = meter.SyncFloat64().Histogram("imagejob_collector_duration", instrument.WithDescription("Distribution of how long it took for collector imagejobs"), instrument.WithUnit(unit.Milliseconds))
	if err != nil {
		klog.InfoS("Failed to register instrument: ImageJobCollectorDuration")
		return err
	}

	if imageJobEraserDuration, err = meter.SyncFloat64().Histogram("imagejob_eraser_duration", instrument.WithDescription("Distribution of how long it took for eraser imagejobs"), instrument.WithUnit(unit.Milliseconds)); err != nil {
		klog.InfoS("Failed to register instrument: ImageJobEraserDuration")
		return err
	}

	if imageJobCollectorTotal, err = meter.SyncInt64().Counter("imagejob_collector_total", instrument.WithDescription("Count of total number of collector imagejobs scheduled")); err != nil {
		klog.InfoS("Failed to register instrument: ImageJobCollectorTotal")
		return err
	}

	if imageJobEraserTotal, err = meter.SyncInt64().Counter("imagejob_eraser_total", instrument.WithDescription("Count of total number of eraser imagejobs scheduled")); err != nil {
		klog.InfoS("Failed to register instrument: ImageJobEraserTotal")
		return err
	}

	if podsCompleted, err = meter.SyncInt64().Counter("pods_completed", instrument.WithDescription("Count of total number of imagejob pods succeeded")); err != nil {
		klog.InfoS("Failed to register instrument: PodsCompleted")
		return err
	}

	if podsFailed, err = meter.SyncInt64().Counter("pods_failed", instrument.WithDescription("Count of total number of imagejob pods failed")); err != nil {
		klog.InfoS("Failed to register instrument: PodsFailed")
		return err
	}

	return nil
}

func RecordImagesRemoved() {
	var err error

	if imagesRemoved != nil {
		imagesRemoved.Add(context.Background(), 1)
	} else {
		if imagesRemoved, err = meter.SyncInt64().Counter("images_removed", instrument.WithDescription("Count of total number of images removed")); err != nil {
			klog.Errorf("Failed to register instrument: ImagesRemoved %v", err)
		}

		imagesRemoved.Add(context.Background(), 1)
	}
}

func RecordImageJobCollectorDuration(duration float64) {
	imageJobCollectorDuration.Record(context.Background(), duration)
}

func RecordImageJobEraserDuration(duration float64) {
	imageJobEraserDuration.Record(context.Background(), duration)
}

func RecordNonCompliantImages(count int64) {
	var err error

	if nonCompliantImages != nil {
		nonCompliantImages.Add(context.Background(), count)
	} else {
		if nonCompliantImages, err = meter.SyncInt64().Counter("non_compliant_images", instrument.WithDescription("Count of total number of vulnerable images found")); err != nil {
			klog.Errorf("Failed to register instrument: NonCompliantImages %v", err)
		}
		nonCompliantImages.Add(context.Background(), count)
	}
}

func RecordImageJobCollectorTotal() {
	imageJobCollectorTotal.Add(context.Background(), 1)
}

func RecordImageJobEraserTotal() {
	imageJobEraserTotal.Add(context.Background(), 1)
}

func RecordPodsCompleted(count int64) {
	podsCompleted.Add(context.Background(), count)
}

func RecordPodsFailed(count int64) {
	podsFailed.Add(context.Background(), count)
}
