package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"log"
	"time"
)

func Jaeger(logger *log.Logger) gin.HandlerFunc {
	logger.SetPrefix("jaeger:")
	return func(c *gin.Context) {
		var newCtx = context.Background()
		var parentSpan opentracing.Span
		tracer, closer, err := NewJaegerTracer("micro-shop", viper.GetString("jaeger.addr"))
		if err != nil {
			logger.Printf("new tracer error:%s", err.Error())
			c.Next()
			return
		}
		defer closer.Close()

		opentracing.SetGlobalTracer(tracer)
		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			parentSpan = tracer.StartSpan(c.Request.URL.Path)
			defer parentSpan.Finish()
		} else {
			parentSpan, newCtx = opentracing.StartSpanFromContextWithTracer(
				c.Request.Context(),
				tracer,
				c.Request.URL.Path,
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			)
			defer parentSpan.Finish()
		}
		newCtx, cancel := context.WithTimeout(newCtx, time.Second*5)
		defer cancel()
		c.Request.WithContext(newCtx)
		c.Next()
	}
}

func NewJaegerTracer(srvName string, endpoint string) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: endpoint,
		},
		ServiceName: srvName,
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}
