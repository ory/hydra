package tracing

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	jeagerConf "github.com/uber/jaeger-client-go/config"
)

type Tracer struct {
	ServiceName  string
	Provider     string
	Logger       *logrus.Logger
	JaegerConfig *JaegerConfig

	tracer opentracing.Tracer
	closer io.Closer
}

type JaegerConfig struct {
	LocalAgentHostPort string
	SamplerType        string
	SamplerValue       float64
	SamplerServerUrl   string
}

func (t *Tracer) Setup() {
	switch t.Provider {
	case "jaeger":
		jc := jeagerConf.Configuration{
			Sampler: &jeagerConf.SamplerConfig{
				SamplingServerURL: t.JaegerConfig.SamplerServerUrl,
				Type:              t.JaegerConfig.SamplerType,
				Param:             t.JaegerConfig.SamplerValue,
			},
			Reporter: &jeagerConf.ReporterConfig{
				LocalAgentHostPort: t.JaegerConfig.LocalAgentHostPort,
			},
		}

		closer, err := jc.InitGlobalTracer(
			t.ServiceName,
		)

		if err != nil {
			t.Logger.Warnf("Could not initialize jaeger tracer: %s", err.Error())
			return
		}

		t.closer = closer
		t.tracer = opentracing.GlobalTracer()
		t.Logger.Infof("Jaeger tracer configured!")
	default:
		if len(t.Provider) > 0 {
			t.Logger.Warnf("Unknown tracer %q - tracer not initialized", t.Provider)
		}
	}
}

func (t *Tracer) IsLoaded() bool {
	if t == nil || t.tracer == nil {
		return false
	}
	return true
}

func (t *Tracer) Close() {
	if t.closer != nil {
		err := t.closer.Close()
		if err != nil {
			t.Logger.Warn(err)
		}
	}
}
