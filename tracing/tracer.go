package tracing

import (
	"errors"
	"fmt"
	"io"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"
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

func (t *Tracer) Setup() error {
	switch strings.ToLower(t.Provider) {
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
			return err
		}

		t.closer = closer
		t.tracer = opentracing.GlobalTracer()
		t.Logger.Infof("Jaeger tracer configured!")
	case "":
		t.Logger.Infof("No tracer configured - skipping tracing setup")
	default:
		return errors.New(fmt.Sprintf("unknown tracer: %s", t.Provider))
	}
	return nil
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

func HelpMessage() string {
	return `- TRACING_PROVIDER: Set this to the tracing backend you wish to use.

	Supported tracing backends: [jaeger]

	Example: TRACING_PROVIDER=jaeger

- TRACING_PROVIDER_JAEGER_SAMPLING_SERVER_URL: The address of jaeger-agent's HTTP sampling server

	Example: TRACING_PROVIDER_JAEGER_SAMPLING_SERVER_URL=http://localhost:5778/sampling

- TRACING_PROVIDER_JAEGER_SAMPLING_TYPE: The type of the sampler you want to use

	Supported values: [const, probabilistic, ratelimiting]

	Default: const

	Example: TRACING_PROVIDER_JAEGER_SAMPLING_TYPE=const

- TRACING_PROVIDER_JAEGER_SAMPLING_VALUE: The value passed to the sampler type that has been configured.

	Supported values: This is dependant on the sampling strategy used:
		- const: 0 or 1 (all or nothing)
		- rateLimiting: a constant rate (e.g. setting this to 3 will sample requests with the rate of 3 traces per second)
		- probabilistic: a value between 0..1

	Example: TRACING_PROVIDER_JAEGER_SAMPLING_VALUE=1

- TRACING_PROVIDER_JAEGER_LOCAL_AGENT_ADDRESS: The address of the jaeger-agent where spans should be sent to

	Example: TRACING_PROVIDER_JAEGER_LOCAL_AGENT_ADDRESS=127.0.0.1:6831

- TRACING_SERVICE_NAME: Specifies the service name to use on the tracer.

	Default: ORY Hydra

	Example: TRACING_SERVICE_NAME="ORY Hydra"
`
}
