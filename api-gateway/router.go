package api_gateway

import (
	"bytes"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/naoina/denco"
	"github.com/valyala/fasthttp"
)

var (
	methodOptions = []byte(http.MethodOptions)
	router        = NewRouter()
)

func Listen(addr string) error {
	router.Build()
	server := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			if bytes.Equal(methodOptions, ctx.Method()) {
				header := &ctx.Response.Header
				header.Set("Access-Control-Allow-Origin", "*")
				header.Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")
				header.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
				header.Set("Access-Control-Expose-Headers", "X-Count")
				header.Set("Access-Control-Max-Age", strconv.FormatInt(int64(time.Hour/time.Second*24), 10))
				return
			}
			router.Handler(ctx)
			header := &ctx.Response.Header
			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Max-Age", strconv.FormatInt(int64(time.Hour/time.Second*24), 10))
		},
		ReadBufferSize: 4096 * 4,
		//MaxKeepaliveDuration: time.Second * 10,
	}
	return server.ListenAndServe(addr)
}

type Router struct {
	routes  map[string][]denco.Record
	routers map[string]*denco.Router
	mu      *sync.Mutex
}

func (r *Router) Handler(ctx *fasthttp.RequestCtx) {
	r.mu.Lock()
	router, ok := r.routers[string(ctx.Method())]
	r.mu.Unlock()
	if !ok {
		ctx.NotFound()
		return
	}
	handler, params, ok := router.Lookup(string(ctx.URI().Path()))
	if !ok {
		ctx.NotFound()
		return
	}
	for _, param := range params {
		ctx.SetUserValue(param.Name, param.Value)
	}
	handler.(fasthttp.RequestHandler)(ctx)
}

func (r *Router) Build() {
	for method, routes := range r.routes {
		r.routers[method] = denco.New()
		r.routers[method].Build(routes)
	}
}

func (r *Router) add(method string, path string, handler fasthttp.RequestHandler) {
	r.mu.Lock()
	r.routes[method] = append(r.routes[method], denco.NewRecord(path, handler))
	r.mu.Unlock()
}

func (r *Router) HEAD(path string, handler fasthttp.RequestHandler) {
	r.add(http.MethodHead, path, handler)
}

func (r *Router) GET(path string, handler fasthttp.RequestHandler) {
	r.add(http.MethodGet, path, handler)
}

func (r *Router) PUT(path string, handler fasthttp.RequestHandler) {
	r.add(http.MethodPut, path, handler)
}

func (r *Router) POST(path string, handler fasthttp.RequestHandler) {
	r.add(http.MethodPost, path, handler)
}

func (r *Router) DELETE(path string, handler fasthttp.RequestHandler) {
	r.add(http.MethodDelete, path, handler)
}

func NewRouter() *Router {
	return &Router{
		routes: map[string][]denco.Record{
			http.MethodHead:   {},
			http.MethodGet:    {},
			http.MethodPut:    {},
			http.MethodPost:   {},
			http.MethodDelete: {},
		},
		routers: map[string]*denco.Router{},
		mu:      &sync.Mutex{},
	}
}
