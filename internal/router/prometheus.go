package router

// Copied code from https://github.com/0neSe7en/echo-prometheus
import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type (
	// PrometheusConfig represents basic prometheus exporter
	PrometheusConfig struct {
		// Skipper echo skipper
		Skipper   middleware.Skipper
		Namespace string
	}

	// EchoMetrics represents prometheus metrics for echo
	EchoMetrics struct {
		echoReqQPS      *prometheus.CounterVec
		echoReqDuration *prometheus.SummaryVec
		echoOutBytes    prometheus.Summary
	}
)

// NewEchoMetrics creates and registers echo metrics. This function will panic on multiple call.
func NewEchoMetrics(namespace string) EchoMetrics {
	var em EchoMetrics

	em.echoReqQPS = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_request_total",
			Help:      "HTTP requests processed.",
		},
		[]string{"code", "method", "host", "url"},
	)
	em.echoReqDuration = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request latencies in seconds.",
		},
		[]string{"method", "host", "url"},
	)
	em.echoOutBytes = promauto.NewSummary(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "http_response_size_bytes",
			Help:      "HTTP response bytes.",
		},
	)

	return em
}

// NewPrometheusMiddleware returns new prometheus exporter with default config
func NewPrometheusMiddleware(namespace string) echo.MiddlewareFunc {
	return NewMetricWithConfig(PrometheusConfig{
		Skipper:   middleware.DefaultSkipper,
		Namespace: namespace,
	})
}

// NewMetricWithConfig creates a new prometheus with config
func NewMetricWithConfig(config PrometheusConfig) echo.MiddlewareFunc {
	em := NewEchoMetrics(config.Namespace)

	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()

			if err := next(c); err != nil {
				c.Error(err)
			}

			uri := req.URL.Path
			status := strconv.Itoa(res.Status)
			elapsed := time.Since(start).Seconds()
			bytesOut := float64(res.Size)

			em.echoReqQPS.WithLabelValues(status, req.Method, req.Host, uri).Inc()
			em.echoReqDuration.WithLabelValues(req.Method, req.Host, uri).Observe(elapsed)
			em.echoOutBytes.Observe(bytesOut)

			return nil
		}
	}
}
