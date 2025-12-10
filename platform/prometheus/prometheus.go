package prometheus

import (
	"log/slog"

	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Métricas HTTP
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total de solicitudes HTTP recibidas",
		},
		[]string{"method", "endpoint", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duración de las solicitudes HTTP en segundos",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Métricas de base de datos
	DBConnectionsActive = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Número de conexiones activas a la base de datos",
		},
	)

	DBQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Total de queries ejecutadas en la base de datos",
		},
		[]string{"operation", "status"},
	)

	// Métricas de aplicación
	EmployeesRegistered = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "employees_registered_total",
			Help: "Total de empleados registrados",
		},
	)

	MessagesCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "messages_created_total",
			Help: "Total de mensajes creados",
		},
	)
)

// InitMetrics inicializa y registra todas las métricas de Prometheus
func InitMetrics() error {
	slog.Info(logger.LogPrometheusInit)

	// Registrar métricas HTTP
	if err := prometheus.Register(HTTPRequestsTotal); err != nil {
		slog.Error(logger.LogPrometheusInitError, slog.String("metric", "http_requests_total"), slog.String("error", err.Error()))
		return err
	}

	if err := prometheus.Register(HTTPRequestDuration); err != nil {
		slog.Error(logger.LogPrometheusInitError, slog.String("metric", "http_request_duration_seconds"), slog.String("error", err.Error()))
		return err
	}

	// Registrar métricas de base de datos
	if err := prometheus.Register(DBConnectionsActive); err != nil {
		slog.Error(logger.LogPrometheusInitError, slog.String("metric", "db_connections_active"), slog.String("error", err.Error()))
		return err
	}

	if err := prometheus.Register(DBQueriesTotal); err != nil {
		slog.Error(logger.LogPrometheusInitError, slog.String("metric", "db_queries_total"), slog.String("error", err.Error()))
		return err
	}

	// Registrar métricas de aplicación
	if err := prometheus.Register(EmployeesRegistered); err != nil {
		slog.Error(logger.LogPrometheusInitError, slog.String("metric", "employees_registered_total"), slog.String("error", err.Error()))
		return err
	}

	if err := prometheus.Register(MessagesCreated); err != nil {
		slog.Error(logger.LogPrometheusInitError, slog.String("metric", "messages_created_total"), slog.String("error", err.Error()))
		return err
	}

	slog.Info(logger.LogPrometheusInitOK)
	return nil
}

// Handler devuelve el handler HTTP de Prometheus para el endpoint /metrics
func Handler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
