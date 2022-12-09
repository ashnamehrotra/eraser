package template

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	eraserv1alpha1 "github.com/Azure/eraser/api/v1alpha1"
	"github.com/go-logr/logr"
	"golang.org/x/sys/unix"

	"github.com/Azure/eraser/pkg/metrics"
	util "github.com/Azure/eraser/pkg/utils"
	"go.opentelemetry.io/otel/metric/global"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type ImageProvider interface {
	ReceiveImages() []eraserv1alpha1.Image
	SendImages(vulnerableImages, failedImages []eraserv1alpha1.Image)
	Finish()
}

type config struct {
	ctx                    context.Context
	log                    logr.Logger
	deleteScanFailedImages bool
	reportMetrics          bool
}

type ConfigFunc func(*config)

func NewImageProvider(funcs ...ConfigFunc) ImageProvider {
	// default config
	cfg := &config{
		ctx:                    context.Background(),
		log:                    logf.Log.WithName("scanner"),
		deleteScanFailedImages: true,
		reportMetrics:          false,
	}

	// apply user config
	for _, f := range funcs {
		f(cfg)
	}

	return cfg
}

func (cfg *config) ReceiveImages() []eraserv1alpha1.Image {
	var err error

	if err := unix.Mkfifo(util.EraseCompleteScanPath, util.PipeMode); err != nil {
		cfg.log.Error(err, "failed to create pipe", "pipeName", util.EraseCompleteScanPath)
		os.Exit(1)
	}

	err = os.Chmod(util.EraseCompleteScanPath, 0o666)
	if err != nil {
		cfg.log.Error(err, "unable to enable pipe for writing", "pipeName", util.EraseCompleteScanPath)
		os.Exit(1)
	}

	allImages, err := util.ReadCollectScanPipe(cfg.ctx)
	if err != nil {
		cfg.log.Error(err, "unable to read images from collect scan pipe")
		os.Exit(1)
	}

	return allImages
}

func (cfg *config) SendImages(vulnerableImages, failedImages []eraserv1alpha1.Image) {
	if cfg.deleteScanFailedImages {
		vulnerableImages = append(vulnerableImages, failedImages...)
	}

	if err := util.WriteScanErasePipe(vulnerableImages); err != nil {
		cfg.log.Error(err, "unable to write non-compliant images to scan erase pipe")
		os.Exit(1)
	}

	if cfg.reportMetrics {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		exporter, reader, provider := metrics.ConfigureMetrics(ctx, cfg.log, os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
		global.SetMeterProvider(provider)

		defer metrics.ExportMetrics(cfg.log, exporter, reader, provider)

		if err := metrics.RecordMetricsScanner(ctx, global.MeterProvider(), len(vulnerableImages)); err != nil {
			cfg.log.Error(err, "error recording metrics")
		}
	}
}

func (cfg *config) Finish() {
	file, err := os.OpenFile(util.EraseCompleteScanPath, os.O_RDONLY, 0)
	if err != nil {
		cfg.log.Error(err, "failed to open pipe", "pipeName", util.EraseCompleteScanPath)
		os.Exit(1)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		cfg.log.Error(err, "failed to read pipe", "pipeName", util.EraseCompleteScanPath)
		os.Exit(1)
	}

	file.Close()

	if string(data) != util.EraseCompleteMessage {
		cfg.log.Info("garbage in pipe", "pipeName", util.EraseCompleteScanPath, "in_pipe", string(data))
		os.Exit(1)
	}

	cfg.log.Info("scanning complete, exiting")
}

func WithContext(ctx context.Context) ConfigFunc {
	return func(cfg *config) {
		cfg.ctx = ctx
	}
}

func WithDeleteScanFailedImages(deleteScanFailedImages bool) ConfigFunc {
	return func(cfg *config) {
		cfg.deleteScanFailedImages = deleteScanFailedImages
	}
}

func WithLogger(log logr.Logger) ConfigFunc {
	return func(cfg *config) {
		cfg.log = log
	}
}

func WithMetrics(reportMetrics bool) ConfigFunc {
	return func(cfg *config) {
		cfg.reportMetrics = reportMetrics
	}
}
