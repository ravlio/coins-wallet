package account

import (
	"io"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/examples/addsvc/pkg/addendpoint"
	"github.com/go-kit/kit/examples/addsvc/pkg/addservice"
	"github.com/go-kit/kit/examples/addsvc/pkg/addtransport"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	dnssd "github.com/go-kit/kit/sd/dnssrv"
	"github.com/go-kit/kit/sd/lb"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/ravlio/wallet/account"
	"google.golang.org/grpc"
)

func New(r *mux.Router, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) {
	var (
		endpoints = addendpoint.Set{}
		instancer = dnssd.NewInstancer("account", time.Millisecond*10, logger)
	)
	{
		factory := factory(addendpoint.MakeSumEndpoint, otTracer, zipkinTracer, logger)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(*retryMax, *retryTimeout, balancer)
		endpoints.SumEndpoint = retry
	}
	{
		factory := factory(addendpoint.MakeConcatEndpoint, otTracer, zipkinTracer, logger)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(*retryMax, *retryTimeout, balancer)
		endpoints.ConcatEndpoint = retry
	}

	// Here we leverage the fact that addsvc comes with a constructor for an
	// HTTP handler, and just install it under a particular path prefix in
	// our router.

	r.PathPrefix("/addsvc").Handler(http.StripPrefix("/addsvc", addtransport.NewHTTPHandler(endpoints, tracer, zipkinTracer, logger)))
}

func factory(makeEndpoint func(account.Client) endpoint.Endpoint, tracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := addtransport.NewGRPCClient(conn, tracer, zipkinTracer, logger)
		endpoint := makeEndpoint(service)

		// Notice that the addsvc gRPC client converts the connection to a
		// complete addsvc, and we just throw away everything except the method
		// we're interested in. A smarter factory would mux multiple methods
		// over the same connection. But that would require more work to manage
		// the returned io.Closer, e.g. reference counting. Since this is for
		// the purposes of demonstration, we'll just keep it simple.

		return endpoint, conn, nil
	}
}
