package tracing

import (
	"io"
	"strings"

	"github.com/opentracing/opentracing-go"
	zipkinOT "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinHttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/pkg/errors"
	jaegerConf "github.com/uber/jaeger-client-go/config"
	jaegerZipkin "github.com/uber/jaeger-client-go/zipkin"

	"github.com/ory/x/logrusx"
)

type Tracer struct {
	ServiceName  string
	Provider     string
	Logger       *logrusx.Logger
	JaegerConfig *JaegerConfig
	ZipkinConfig *ZipkinConfig

	tracer opentracing.Tracer
	closer io.Closer
}

type JaegerConfig struct {
	LocalAgentHostPort string
	SamplerType        string
	SamplerValue       float64
	SamplerServerUrl   string
	Propagation        string
}

type ZipkinConfig struct {
	ServerURL string
}

func (t *Tracer) Setup() error {
	switch strings.ToLower(t.Provider) {
	case "jaeger":
		jc, err := jaegerConf.FromEnv()

		if err != nil {
			return err
		}

		if t.JaegerConfig.SamplerServerUrl != "" {
			jc.Sampler.SamplingServerURL = t.JaegerConfig.SamplerServerUrl
		}

		if t.JaegerConfig.SamplerType != "" {
			jc.Sampler.Type = t.JaegerConfig.SamplerType
		}

		if t.JaegerConfig.SamplerValue != 0 {
			jc.Sampler.Param = t.JaegerConfig.SamplerValue
		}

		if t.JaegerConfig.LocalAgentHostPort != "" {
			jc.Reporter.LocalAgentHostPort = t.JaegerConfig.LocalAgentHostPort
		}

		var configs []jaegerConf.Option

		// This works in other jaeger clients, but is not part of jaeger-client-go
		if t.JaegerConfig.Propagation == "b3" {
			zipkinPropagator := jaegerZipkin.NewZipkinB3HTTPHeaderPropagator()
			configs = append(
				configs,
				jaegerConf.Injector(opentracing.HTTPHeaders, zipkinPropagator),
				jaegerConf.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
			)
		}

		closer, err := jc.InitGlobalTracer(
			t.ServiceName,
			configs...,
		)

		if err != nil {
			return err
		}

		t.closer = closer
		t.tracer = opentracing.GlobalTracer()
		t.Logger.Infof("Jaeger tracer configured!")
	case "zipkin":
		if t.ZipkinConfig.ServerURL == "" {
			return errors.Errorf("Zipkin's server url is required")
		}

		reporter := zipkinHttp.NewReporter(t.ZipkinConfig.ServerURL)

		endpoint, err := zipkin.NewEndpoint(t.ServiceName, "")

		if err != nil {
			return err
		}

		nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))

		if err != nil {
			return err
		}

		opentracing.SetGlobalTracer(zipkinOT.Wrap(nativeTracer))

		t.closer = reporter
		t.tracer = opentracing.GlobalTracer()
		t.Logger.Infof("Zipkin tracer configured!")
	case "":
		t.Logger.Infof("No tracer configured - skipping tracing setup")
	default:
		return errors.Errorf("unknown tracer: %s", t.Provider)
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

	Supported tracing backends: [jaeger, zipkin]

	Example: TRACING_PROVIDER=jaeger

- TRACING_PROVIDERS_JAEGER_SAMPLING_SERVER_URL: The address of jaeger-agent's HTTP sampling server

	Example: TRACING_PROVIDER_JAEGER_SAMPLING_SERVER_URL=http://localhost:5778/sampling

- TRACING_PROVIDERS_JAEGER_SAMPLING_TYPE: The type of the sampler you want to use

	Supported values: [const, probabilistic, ratelimiting]

	Default: const

	Example: TRACING_PROVIDER_JAEGER_SAMPLING_TYPE=const

- TRACING_PROVIDERS_JAEGER_SAMPLING_VALUE: The value passed to the sampler type that has been configured.

	Supported values: This is dependant on the sampling strategy used:
		- const: 0 or 1 (all or nothing)
		- rateLimiting: a constant rate (e.g. setting this to 3 will sample requests with the rate of 3 traces per second)
		- probabilistic: a value between 0..1

	Example: TRACING_PROVIDER_JAEGER_SAMPLING_VALUE=1

- TRACING_PROVIDERS_JAEGER_LOCAL_AGENT_ADDRESS: The address of the jaeger-agent where spans should be sent to

	Example: TRACING_PROVIDER_JAEGER_LOCAL_AGENT_ADDRESS=127.0.0.1:6831

- TRACING_PROVIDERS_JAEGER_PROPAGATION: The tracing header propagation format. Defaults to jaeger.

	Example: TRACING_PROVIDER_JAEGER_PROPAGATION=b3

- TRACING_PROVIDERS_ZIPKIN_SERVER_URL: The address of Zipkin server

	Example: TRACING_PROVIDERS_ZIPKIN_SERVER_URL=http://localhost:9411/api/v2/spans

- TRACING_SERVICE_NAME: Specifies the service name to use on the tracer.

	Default: ORY Hydra

	Example: TRACING_SERVICE_NAME="ORY Hydra"
`
}
