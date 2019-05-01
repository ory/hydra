package metrics

import (
	"context"
	"fmt"
	"log"
	"net"
	"runtime"
	"time"

	"github.com/pborman/uuid"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/segmentio/analytics-go"
	"github.com/sirupsen/logrus"

	"github.com/ory/x/resilience"
)

type void struct {
}

func (v *void) Logf(format string, args ...interface{}) {
}

func (v *void) Errorf(format string, args ...interface{}) {
}

// SegmentOptions provides configuration settings for the SegmentBridge
type SegmentOptions struct {
	// Service represents the service name, for example "ory-hydra".
	Service string

	// ClusterID represents the cluster id, typically a hash of some unique configuration properties.
	ClusterID string

	// IsDevelopment should be true if we assume that we're in a development environment.
	IsDevelopment bool

	// WriteKey is the segment API key.
	WriteKey string
	// BuildVersion represents the build version.
	BuildVersion string

	// BuildHash represents the build git hash.
	BuildHash string

	// BuildTime represents the build time.
	BuildTime string
}

type SegmentBridge struct {
	o       *SegmentOptions
	client  analytics.Client
	context *analytics.Context
	g       prometheus.Gatherer
	l       logrus.FieldLogger
}

func (s *SegmentBridge) Push(ctx context.Context) error {
	mfs, err := s.g.Gather()
	if err != nil {
		return err
	}

	for _, v := range mfs {
		if v.GetType() == dto.MetricType_HISTOGRAM || v.GetType() == dto.MetricType_SUMMARY {
			continue
		}
		for _, m := range v.Metric {
			s.enqueueMetric(v.GetName(), v.GetType(), m)
		}
	}

	return nil
}

func (s *SegmentBridge) getValueFromMetric(t dto.MetricType, m *dto.Metric) float64 {
	switch t {
	case dto.MetricType_GAUGE:
		return m.GetGauge().GetValue()
	case dto.MetricType_COUNTER:
		return m.GetCounter().GetValue()
	}
	return 0
}

func (s *SegmentBridge) enqueueMetric(name string, t dto.MetricType, m *dto.Metric) error {
	p := analytics.Properties{}
	for _, label := range m.GetLabel() {
		if label.GetName() != "" && label.GetValue() != "" {
			log.Println("Tracking label", label.GetName(), label.GetValue())
			p.Set(*label.Name, *label.Value)
		}
	}
	val := s.getValueFromMetric(t, m)
	p.Set("_value", val)
	log.Println("_value=", val)
	err := s.client.Enqueue(analytics.Track{
		UserId:     s.o.ClusterID,
		Properties: p,
		Context:    s.context,
		Event:      name,
	})

	if err != nil {
		log.Println(err)
	}

	return nil
}

func NewSegmentBridge(ctx context.Context, o *SegmentOptions, logger logrus.FieldLogger, gatherer prometheus.Gatherer) (*SegmentBridge, error) {
	client, err := analytics.NewWithConfig(o.WriteKey, analytics.Config{
		Interval:  time.Hour * 24,
		BatchSize: 100,
	})

	if err != nil {
		return nil, err
	}

	oi := analytics.OSInfo{
		Version: fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH),
	}
	actx := &analytics.Context{
		IP: net.IPv4(0, 0, 0, 0),
		App: analytics.AppInfo{
			Name:    o.Service,
			Version: o.BuildVersion,
			Build:   fmt.Sprintf("%s/%s/%s", o.BuildVersion, o.BuildHash, o.BuildTime),
		},
		OS: oi,
		Traits: analytics.NewTraits().
			Set("optedOut", false).
			Set("instanceId", uuid.New()).
			Set("isDevelopment", o.IsDevelopment),
		UserAgent: "github.com/ory/x/metricsx.Service/v0.0.1",
	}

	if err := resilience.Retry(logger, time.Minute*5, time.Hour*24*30, func() error {
		return client.Enqueue(analytics.Identify{
			UserId:  o.ClusterID,
			Traits:  actx.Traits,
			Context: actx,
		})
	}); err != nil {
		logger.WithError(err).Debug("Could not commit anonymized environment information")
	}
	return &SegmentBridge{
		client: client,
		g:      gatherer,
		o:      o,
		l:      logger,
	}, nil
}
