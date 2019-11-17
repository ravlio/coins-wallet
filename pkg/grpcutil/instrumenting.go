package grpcutil

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/ravlio/wallet/pkg/errutil"
	"google.golang.org/grpc"
)

var (
	requestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "requests total",
	}, []string{"service", "method", "code"})
	requestDuration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "request_duration_ns",
		Help: "request duration",
	}, []string{"service", "method", "code"})

	requestClientDuration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "request_duration_ns",
		Help: "request duration",
	}, []string{"method", "code"})
	requestClientTransportDuration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "clickaine",
		Subsystem: "grpc_client",
		Name:      "transport_duration_ns",
		Help:      "transport duration",
	}, []string{"method", "code"})
)

func init() {
	prometheus.MustRegister(requestsTotal, requestDuration, requestClientDuration, requestClientTransportDuration)
}

func getServiceAndMethod(fullMethodName string) (string, string) {
	parts := strings.Split(fullMethodName, "/")
	return parts[1], parts[2]
}

func MakeServerInstrumentingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, request interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		svc, method := getServiceAndMethod(info.FullMethod)
		logCtx := LogRequest(svc, method, request.(json.Marshaler))

		startTime := time.Now()
		response, err := handler(ctx, request)
		took := time.Since(startTime)

		logCtx.LogResponse(ctx, response.(json.Marshaler), err)

		var code int

		if err != nil {
			if own, ok := err.(errutil.Error); ok {
				code = own.GetCode()
			}
		}
		labels := prometheus.Labels{
			"service": svc,
			"method":  method,
			"code":    strconv.Itoa(code),
		}

		requestsTotal.With(labels).Inc()
		requestDuration.With(labels).Observe(float64(took.Nanoseconds()))

		return response, err
	}
}
