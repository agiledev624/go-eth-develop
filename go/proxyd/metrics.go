package proxyd

import (
	"context"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	MetricsNamespace = "proxyd"

	RPCRequestSourceHTTP = "http"
	RPCRequestSourceWS   = "ws"

	BackendProxyd = "proxyd"
	SourceClient  = "client"
	SourceBackend = "backend"
	MethodUnknown = "unknown"
)

var PayloadSizeBuckets = []float64{10, 50, 100, 500, 1000, 5000, 10000, 100000, 1000000}

var (
	rpcRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Name:      "rpc_requests_total",
		Help:      "Count of total client RPC requests.",
	})

	rpcForwardsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Name:      "rpc_forwards_total",
		Help:      "Count of total RPC requests forwarded to each backend.",
	}, []string{
		"auth",
		"backend_name",
		"method_name",
		"source",
	})

	rpcBackendHTTPResponseCodesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Name:      "rpc_backend_http_response_codes_total",
		Help:      "Count of total backend responses by HTTP status code.",
	}, []string{
		"auth",
		"backend_name",
		"method_name",
		"status_code",
	})

	rpcErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Name:      "rpc_errors_total",
		Help:      "Count of total RPC errors.",
	}, []string{
		"auth",
		"backend_name",
		"method_name",
		"error_code",
	})

	rpcSpecialErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Name:      "rpc_special_errors_total",
		Help:      "Count of total special RPC errors.",
	}, []string{
		"auth",
		"backend_name",
		"method_name",
		"error_type",
	})

	rpcBackendRequestDurationSumm = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  MetricsNamespace,
		Name:       "rpc_backend_request_duration_seconds",
		Help:       "Summary of backend response times broken down by backend and method name.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
	}, []string{
		"backend_name",
		"method_name",
	})

	activeClientWsConnsGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: MetricsNamespace,
		Name:      "active_client_ws_conns",
		Help:      "Gauge of active client WS connections.",
	}, []string{
		"auth",
	})

	activeBackendWsConnsGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: MetricsNamespace,
		Name:      "active_backend_ws_conns",
		Help:      "Gauge of active backend WS connections.",
	}, []string{
		"backend_name",
	})

	unserviceableRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Name:      "unserviceable_requests_total",
		Help:      "Count of total requests that were rejected due to no backends being available.",
	}, []string{
		"auth",
		"request_source",
	})

	httpResponseCodesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Name:      "http_response_codes_total",
		Help:      "Count of total HTTP response codes.",
	}, []string{
		"status_code",
	})

	httpRequestDurationSumm = promauto.NewSummary(prometheus.SummaryOpts{
		Namespace:  MetricsNamespace,
		Name:       "http_request_duration_seconds",
		Help:       "Summary of HTTP request durations, in seconds.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
	})

	wsMessagesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Name:      "ws_messages_total",
		Help:      "Count of total websocket messages including protocol control.",
	}, []string{
		"auth",
		"backend_name",
		"source",
	})

	redisErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Name:      "redis_errors_total",
		Help:      "Count of total Redis errors.",
	}, []string{
		"source",
	})

	requestPayloadSizesGauge = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: MetricsNamespace,
		Name:      "request_payload_sizes",
		Help:      "Histogram of client request payload sizes.",
		Buckets:   PayloadSizeBuckets,
	}, []string{
		"auth",
	})

	responsePayloadSizesGauge = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: MetricsNamespace,
		Name:      "response_payload_sizes",
		Help:      "Histogram of client response payload sizes.",
		Buckets:   PayloadSizeBuckets,
	}, []string{
		"auth",
	})

	cacheHitsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Name:      "cache_hits_total",
		Help:      "Number of cache hits.",
	}, []string{
		"method",
	})

	cacheMissesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Name:      "cache_misses_total",
		Help:      "Number of cache misses.",
	}, []string{
		"method",
	})

	lvcErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: MetricsNamespace,
		Name:      "lvc_errors_total",
		Help:      "Count of lvc errors.",
	}, []string{
		"key",
	})

	lvcPollTimeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: MetricsNamespace,
		Name:      "lvc_poll_time_gauge",
		Help:      "Gauge of lvc poll time.",
	}, []string{
		"key",
	})

	rpcSpecialErrors = []string{
		"nonce too low",
		"gas price too high",
		"gas price too low",
		"invalid parameters",
	}
)

func RecordRedisError(source string) {
	redisErrorsTotal.WithLabelValues(source).Inc()
}

func RecordRPCError(ctx context.Context, backendName, method string, err error) {
	rpcErr, ok := err.(*RPCErr)
	var code int
	if ok {
		MaybeRecordSpecialRPCError(ctx, backendName, method, rpcErr)
		code = rpcErr.Code
	} else {
		code = -1
	}

	rpcErrorsTotal.WithLabelValues(GetAuthCtx(ctx), backendName, method, strconv.Itoa(code)).Inc()
}

func RecordWSMessage(ctx context.Context, backendName, source string) {
	wsMessagesTotal.WithLabelValues(GetAuthCtx(ctx), backendName, source).Inc()
}

func RecordUnserviceableRequest(ctx context.Context, source string) {
	unserviceableRequestsTotal.WithLabelValues(GetAuthCtx(ctx), source).Inc()
}

func RecordRPCForward(ctx context.Context, backendName, method, source string) {
	rpcForwardsTotal.WithLabelValues(GetAuthCtx(ctx), backendName, method, source).Inc()
}

func MaybeRecordSpecialRPCError(ctx context.Context, backendName, method string, rpcErr *RPCErr) {
	errMsg := strings.ToLower(rpcErr.Message)
	for _, errStr := range rpcSpecialErrors {
		if strings.Contains(errMsg, errStr) {
			rpcSpecialErrorsTotal.WithLabelValues(GetAuthCtx(ctx), backendName, method, errStr).Inc()
			return
		}
	}
}

func RecordRequestPayloadSize(ctx context.Context, payloadSize int) {
	requestPayloadSizesGauge.WithLabelValues(GetAuthCtx(ctx)).Observe(float64(payloadSize))
}

func RecordResponsePayloadSize(ctx context.Context, payloadSize int) {
	responsePayloadSizesGauge.WithLabelValues(GetAuthCtx(ctx)).Observe(float64(payloadSize))
}

func RecordCacheHit(method string) {
	cacheHitsTotal.WithLabelValues(method).Inc()
}

func RecordCacheMiss(method string) {
	cacheMissesTotal.WithLabelValues(method).Inc()
}
