package metrics

import (
	"context"

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
)

type reporter struct {
	meter metric.Meter
}

type StatsReporter interface {
	RecordImagesRemoved()
	RecordImageJobCollectorDuration(duration float64)
	RecordImageJobEraserDuration(duration float64)
	RecordNonCompliantImages(count int64)
	RecordImageJobCollectorTotal()
	RecordImageJobEraserTotal()
	RecordPodsCompleted(count int64)
	RecordPodsFailed(count int64)
}

func NewStatsReporter() (StatsReporter, error) {
	meter := global.Meter("eraser")

	var err error

	imageJobCollectorDuration, err = meter.SyncFloat64().Histogram("imagejob_collector_duration", instrument.WithDescription("Distribution of how long it took for collector imagejobs"), instrument.WithUnit(unit.Milliseconds))
	if err != nil {
		klog.InfoS("Failed to register instrument: ImageJobCollectorDuration")
		return nil, err
	}

	if imageJobEraserDuration, err = meter.SyncFloat64().Histogram("imagejob_eraser_duration", instrument.WithDescription("Distribution of how long it took for eraser imagejobs"), instrument.WithUnit(unit.Milliseconds)); err != nil {
		klog.InfoS("Failed to register instrument: ImageJobEraserDuration")
		return nil, err
	}

	if imageJobCollectorTotal, err = meter.SyncInt64().Counter("imagejob_collector_total", instrument.WithDescription("Count of total number of collector imagejobs scheduled")); err != nil {
		klog.InfoS("Failed to register instrument: ImageJobCollectorTotal")
		return nil, err
	}

	if imageJobEraserTotal, err = meter.SyncInt64().Counter("imagejob_eraser_total", instrument.WithDescription("Count of total number of eraser imagejobs scheduled")); err != nil {
		klog.InfoS("Failed to register instrument: ImageJobEraserTotal")
		return nil, err
	}

	if podsCompleted, err = meter.SyncInt64().Counter("pods_completed", instrument.WithDescription("Count of total number of imagejob pods succeeded")); err != nil {
		klog.InfoS("Failed to register instrument: PodsCompleted")
		return nil, err
	}

	if podsFailed, err = meter.SyncInt64().Counter("pods_failed", instrument.WithDescription("Count of total number of imagejob pods failed")); err != nil {
		klog.InfoS("Failed to register instrument: PodsFailed")
		return nil, err
	}

	if imagesRemoved, err = meter.SyncInt64().Counter("images_removed", instrument.WithDescription("Count of total number of images removed")); err != nil {
		klog.InfoS("Failed to register instrument: ImagesRemoved")
		return nil, err
	}

	if nonCompliantImages, err = meter.SyncInt64().Counter("non_compliant_images", instrument.WithDescription("Count of total number of non-compliant images found")); err != nil {
		klog.InfoS("Failed to register instrument: NonCompliantImages")
		return nil, err
	}

	return &reporter{meter: meter}, nil
}

func (r *reporter) RecordImagesRemoved() {
	imagesRemoved.Add(context.Background(), 1)
}

func (r *reporter) RecordImageJobCollectorDuration(duration float64) {
	imageJobCollectorDuration.Record(context.Background(), duration)
}

func (r *reporter) RecordImageJobEraserDuration(duration float64) {
	imageJobEraserDuration.Record(context.Background(), duration)
}

func (r *reporter) RecordNonCompliantImages(count int64) {
	nonCompliantImages.Add(context.Background(), count)
}

func (r *reporter) RecordImageJobCollectorTotal() {
	imageJobCollectorTotal.Add(context.Background(), 1)
}

func (r *reporter) RecordImageJobEraserTotal() {
	imageJobEraserTotal.Add(context.Background(), 1)
}

func (r *reporter) RecordPodsCompleted(count int64) {
	podsCompleted.Add(context.Background(), count)
}

func (r *reporter) RecordPodsFailed(count int64) {
	podsFailed.Add(context.Background(), count)
}
