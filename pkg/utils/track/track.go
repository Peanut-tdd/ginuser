package track

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func InitJaeger(serviceName, agentHost string, agentPort int) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: fmt.Sprintf("%s:%d", agentHost, agentPort),
		},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}

func StartSpan(tracer opentracing.Tracer, operationName string) opentracing.Span {
	return tracer.StartSpan(operationName)
}

func GetParentSpan(operationName, traceId string, headers map[string][]string) (opentracing.Span, error) {
	tracer := opentracing.GlobalTracer()

	textMapCarrier := make(opentracing.TextMapCarrier)
	for k, v := range headers {
		if len(v) > 0 {
			textMapCarrier[k] = v[0]
		}
	}

	spanCtx, err := tracer.Extract(opentracing.TextMap, textMapCarrier)
	if err != nil {
		return nil, err
	}

	return tracer.StartSpan(operationName, opentracing.ChildOf(spanCtx)), nil
}

func ExtractTraceId(span opentracing.Span) string {
	if span == nil {
		return ""
	}

	spanCtx, ok := span.Context().(jaeger.SpanContext)
	if !ok {
		return ""
	}

	return spanCtx.TraceID().String()
}
