package prometheusx

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/huanggze/x/httpx"
)

type Metrics struct {
	responseTime    *prometheus.HistogramVec
	totalRequests   *prometheus.CounterVec
	duration        *prometheus.HistogramVec
	responseSize    *prometheus.HistogramVec
	requestSize     *prometheus.HistogramVec
	handlerStatuses *prometheus.CounterVec
}

func (h Metrics) instrumentHandlerStatusBucket(next http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(rw, r)

		status, _ := httpx.GetResponseMeta(rw)

		statusBucket := "unknown"
		switch {
		case status >= 200 && status <= 299:
			statusBucket = "2xx"
		case status >= 300 && status <= 399:
			statusBucket = "3xx"
		case status >= 400 && status <= 499:
			statusBucket = "4xx"
		case status >= 500 && status <= 599:
			statusBucket = "5xx"
		}

		h.handlerStatuses.With(prometheus.Labels{"method": r.Method, "status_bucket": statusBucket}).
			Inc()
	}
}

// Instrument will instrument any http.HandlerFunc with custom metrics
func (h Metrics) Instrument(rw http.ResponseWriter, next http.HandlerFunc, endpoint string) http.HandlerFunc {
	labels := prometheus.Labels{}
	labelsWithEndpoint := prometheus.Labels{"endpoint": endpoint}
	if status, _ := httpx.GetResponseMeta(rw); status != 0 {
		labels = prometheus.Labels{"code": strconv.Itoa(status)}
		labelsWithEndpoint["code"] = labels["code"]
	}
	wrapped := promhttp.InstrumentHandlerResponseSize(h.responseSize.MustCurryWith(labels), next)
	wrapped = promhttp.InstrumentHandlerCounter(h.totalRequests.MustCurryWith(labelsWithEndpoint), wrapped)
	wrapped = promhttp.InstrumentHandlerDuration(h.duration.MustCurryWith(labelsWithEndpoint), wrapped)
	wrapped = promhttp.InstrumentHandlerDuration(h.responseTime.MustCurryWith(prometheus.Labels{"endpoint": endpoint}), wrapped)
	wrapped = promhttp.InstrumentHandlerRequestSize(h.requestSize.MustCurryWith(labels), wrapped)
	wrapped = h.instrumentHandlerStatusBucket(wrapped)

	return wrapped.ServeHTTP
}
