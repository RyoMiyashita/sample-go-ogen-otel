// Code generated by ogen, DO NOT EDIT.

package logo

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/otelogen"
)

var (
	// Allocate option closure once.
	clientSpanKind = trace.WithSpanKind(trace.SpanKindClient)
	// Allocate option closure once.
	serverSpanKind = trace.WithSpanKind(trace.SpanKindServer)
)

type (
	optionFunc[C any] func(*C)
	otelOptionFunc    func(*otelConfig)
)

type otelConfig struct {
	TracerProvider trace.TracerProvider
	Tracer         trace.Tracer
	MeterProvider  metric.MeterProvider
	Meter          metric.Meter
}

func (cfg *otelConfig) initOTEL() {
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}
	if cfg.MeterProvider == nil {
		cfg.MeterProvider = otel.GetMeterProvider()
	}
	cfg.Tracer = cfg.TracerProvider.Tracer(otelogen.Name,
		trace.WithInstrumentationVersion(otelogen.SemVersion()),
	)
	cfg.Meter = cfg.MeterProvider.Meter(otelogen.Name,
		metric.WithInstrumentationVersion(otelogen.SemVersion()),
	)
}

// ErrorHandler is error handler.
type ErrorHandler = ogenerrors.ErrorHandler

type serverConfig struct {
	otelConfig
	NotFound           http.HandlerFunc
	MethodNotAllowed   func(w http.ResponseWriter, r *http.Request, allowed string)
	ErrorHandler       ErrorHandler
	Prefix             string
	Middleware         Middleware
	MaxMultipartMemory int64
}

// ServerOption is server config option.
type ServerOption interface {
	applyServer(*serverConfig)
}

var _ ServerOption = (optionFunc[serverConfig])(nil)

func (o optionFunc[C]) applyServer(c *C) {
	o(c)
}

var _ ServerOption = (otelOptionFunc)(nil)

func (o otelOptionFunc) applyServer(c *serverConfig) {
	o(&c.otelConfig)
}

func newServerConfig(opts ...ServerOption) serverConfig {
	cfg := serverConfig{
		NotFound: http.NotFound,
		MethodNotAllowed: func(w http.ResponseWriter, r *http.Request, allowed string) {
			status := http.StatusMethodNotAllowed
			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Methods", allowed)
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				status = http.StatusNoContent
			} else {
				w.Header().Set("Allow", allowed)
			}
			w.WriteHeader(status)
		},
		ErrorHandler:       ogenerrors.DefaultErrorHandler,
		Middleware:         nil,
		MaxMultipartMemory: 32 << 20, // 32 MB
	}
	for _, opt := range opts {
		opt.applyServer(&cfg)
	}
	cfg.initOTEL()
	return cfg
}

type baseServer struct {
	cfg      serverConfig
	requests metric.Int64Counter
	errors   metric.Int64Counter
	duration metric.Float64Histogram
}

func (s baseServer) notFound(w http.ResponseWriter, r *http.Request) {
	s.cfg.NotFound(w, r)
}

func (s baseServer) notAllowed(w http.ResponseWriter, r *http.Request, allowed string) {
	s.cfg.MethodNotAllowed(w, r, allowed)
}

func (cfg serverConfig) baseServer() (s baseServer, err error) {
	s = baseServer{cfg: cfg}
	if s.requests, err = otelogen.ServerRequestCountCounter(s.cfg.Meter); err != nil {
		return s, err
	}
	if s.errors, err = otelogen.ServerErrorsCountCounter(s.cfg.Meter); err != nil {
		return s, err
	}
	if s.duration, err = otelogen.ServerDurationHistogram(s.cfg.Meter); err != nil {
		return s, err
	}
	return s, nil
}

type clientConfig struct {
	otelConfig
	Client ht.Client
}

// ClientOption is client config option.
type ClientOption interface {
	applyClient(*clientConfig)
}

var _ ClientOption = (optionFunc[clientConfig])(nil)

func (o optionFunc[C]) applyClient(c *C) {
	o(c)
}

var _ ClientOption = (otelOptionFunc)(nil)

func (o otelOptionFunc) applyClient(c *clientConfig) {
	o(&c.otelConfig)
}

func newClientConfig(opts ...ClientOption) clientConfig {
	cfg := clientConfig{
		Client: http.DefaultClient,
	}
	for _, opt := range opts {
		opt.applyClient(&cfg)
	}
	cfg.initOTEL()
	return cfg
}

type baseClient struct {
	cfg      clientConfig
	requests metric.Int64Counter
	errors   metric.Int64Counter
	duration metric.Float64Histogram
}

func (cfg clientConfig) baseClient() (c baseClient, err error) {
	c = baseClient{cfg: cfg}
	if c.requests, err = otelogen.ClientRequestCountCounter(c.cfg.Meter); err != nil {
		return c, err
	}
	if c.errors, err = otelogen.ClientErrorsCountCounter(c.cfg.Meter); err != nil {
		return c, err
	}
	if c.duration, err = otelogen.ClientDurationHistogram(c.cfg.Meter); err != nil {
		return c, err
	}
	return c, nil
}

// Option is config option.
type Option interface {
	ServerOption
	ClientOption
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
//
// If none is specified, the global provider is used.
func WithTracerProvider(provider trace.TracerProvider) Option {
	return otelOptionFunc(func(cfg *otelConfig) {
		if provider != nil {
			cfg.TracerProvider = provider
		}
	})
}

// WithMeterProvider specifies a meter provider to use for creating a meter.
//
// If none is specified, the otel.GetMeterProvider() is used.
func WithMeterProvider(provider metric.MeterProvider) Option {
	return otelOptionFunc(func(cfg *otelConfig) {
		if provider != nil {
			cfg.MeterProvider = provider
		}
	})
}

// WithClient specifies http client to use.
func WithClient(client ht.Client) ClientOption {
	return optionFunc[clientConfig](func(cfg *clientConfig) {
		if client != nil {
			cfg.Client = client
		}
	})
}

// WithNotFound specifies Not Found handler to use.
func WithNotFound(notFound http.HandlerFunc) ServerOption {
	return optionFunc[serverConfig](func(cfg *serverConfig) {
		if notFound != nil {
			cfg.NotFound = notFound
		}
	})
}

// WithMethodNotAllowed specifies Method Not Allowed handler to use.
func WithMethodNotAllowed(methodNotAllowed func(w http.ResponseWriter, r *http.Request, allowed string)) ServerOption {
	return optionFunc[serverConfig](func(cfg *serverConfig) {
		if methodNotAllowed != nil {
			cfg.MethodNotAllowed = methodNotAllowed
		}
	})
}

// WithErrorHandler specifies error handler to use.
func WithErrorHandler(h ErrorHandler) ServerOption {
	return optionFunc[serverConfig](func(cfg *serverConfig) {
		if h != nil {
			cfg.ErrorHandler = h
		}
	})
}

// WithPathPrefix specifies server path prefix.
func WithPathPrefix(prefix string) ServerOption {
	return optionFunc[serverConfig](func(cfg *serverConfig) {
		cfg.Prefix = prefix
	})
}

// WithMiddleware specifies middlewares to use.
func WithMiddleware(m ...Middleware) ServerOption {
	return optionFunc[serverConfig](func(cfg *serverConfig) {
		switch len(m) {
		case 0:
			cfg.Middleware = nil
		case 1:
			cfg.Middleware = m[0]
		default:
			cfg.Middleware = middleware.ChainMiddlewares(m...)
		}
	})
}

// WithMaxMultipartMemory specifies limit of memory for storing file parts.
// File parts which can't be stored in memory will be stored on disk in temporary files.
func WithMaxMultipartMemory(max int64) ServerOption {
	return optionFunc[serverConfig](func(cfg *serverConfig) {
		if max > 0 {
			cfg.MaxMultipartMemory = max
		}
	})
}
