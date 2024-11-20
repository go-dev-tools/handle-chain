package chain

import (
	"context"
	"net/http"
)

// New creates an instance of Handler.
// Requires the Request and Response types to be defined for the Handler.
// For empty request body in the case of GET method use EmptyRequest as the Request type.
// For empty response body in the case of POST, PUT, PATCH and DELETE methods use EmptyResponse as the Response type.
func New[Request, Response any]() Handler[Request, Response] {
	return &handler[Request, Response]{}
}

// Handler composes http.Handler and can be used with any http router.
//
// The order of execution is:
// ParseRequestFunc -> AuthorizeFunc -> ResolveRequestFunc
//
// OnSuccessFunc is called when ResolveRequestFunc doesn't return an error.
// OnErrorFunc is called when ResolveRequestFunc returns an error.
//
// CollectMetricFunc and AuditLogFunc get called whether or not an error occurs.
type Handler[Request, Response any] interface {
	http.Handler

	Parse(ParseRequestFunc[Request]) Handler[Request, Response]

	AuthZ(AuthorizeFunc[Request]) Handler[Request, Response]

	Resolve(ResolveRequestFunc[Request, Response]) Handler[Request, Response]

	OnSuccess(OnSuccessFunc[Response]) Handler[Request, Response]

	OnError(OnErrorFunc) Handler[Request, Response]

	Monitor(CollectMetricFunc[Request]) Handler[Request, Response]

	Audit(AuditLogFunc[Request]) Handler[Request, Response]
}

type EmptyRequest any
type EmptyResponse any

type ParsedRequest[Request any] struct {
	Body        Request
	Headers     http.Header
	PathParams  map[string]string
	QueryParams map[string]string
}

type SuccessResponse[Response any] struct {
	HttpStatusCode int
	Body           Response
	Headers        http.Header
}

type ParseRequestFunc[Request any] func(*http.Request) (ParsedRequest[Request], error)

type AuthorizeFunc[Request any] func(context.Context, ParsedRequest[Request]) error

type ResolveRequestFunc[Request, Response any] func(context.Context, ParsedRequest[Request]) (SuccessResponse[Response], error)

type OnSuccessFunc[Response any] func(context.Context, http.ResponseWriter, SuccessResponse[Response])

type OnErrorFunc func(context.Context, http.ResponseWriter, error)

type CollectMetricFunc[Request any] func(context.Context, ParsedRequest[Request], error)

type AuditLogFunc[Request any] func(context.Context, ParsedRequest[Request], error)

type handler[Request, Response any] struct {
	parseReqF  ParseRequestFunc[Request]
	authF      AuthorizeFunc[Request]
	resolveF   ResolveRequestFunc[Request, Response]
	onSuccessF OnSuccessFunc[Response]
	onErrorF   OnErrorFunc
	metricF    CollectMetricFunc[Request]
	auditF     AuditLogFunc[Request]
}

func (h *handler[Request, Response]) Parse(f ParseRequestFunc[Request]) Handler[Request, Response] {
	h.parseReqF = f
	return h
}

func (h *handler[Request, Response]) AuthZ(f AuthorizeFunc[Request]) Handler[Request, Response] {
	h.authF = f
	return h
}

func (h *handler[Request, Response]) Resolve(f ResolveRequestFunc[Request, Response]) Handler[Request, Response] {
	h.resolveF = f
	return h
}

func (h *handler[Request, Response]) OnSuccess(f OnSuccessFunc[Response]) Handler[Request, Response] {
	h.onSuccessF = f
	return h
}

func (h *handler[Request, Response]) OnError(f OnErrorFunc) Handler[Request, Response] {
	h.onErrorF = f
	return h
}

func (h *handler[Request, Response]) Monitor(f CollectMetricFunc[Request]) Handler[Request, Response] {
	h.metricF = f
	return h
}

func (h *handler[Request, Response]) Audit(f AuditLogFunc[Request]) Handler[Request, Response] {
	h.auditF = f
	return h
}

func (h *handler[Request, Response]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req ParsedRequest[Request]
	var err error

	// Parse the request
	if h.parseReqF != nil {
		req, err = h.parseReqF(r)
		if err != nil {
			h.handleError(ctx, w, req, err)
			return
		}
	}

	// Authorize the request
	if h.authF != nil {
		err = h.authF(ctx, req)
		if err != nil {
			h.handleError(ctx, w, req, err)
			return
		}
	}

	var resp SuccessResponse[Response]

	// Resolve the request and get the response
	if h.resolveF != nil {
		resp, err = h.resolveF(ctx, req)
		if err != nil {
			h.handleError(ctx, w, req, err)
			return
		}
	}

	// Send the success response
	if h.onSuccessF != nil {
		h.onSuccessF(ctx, w, resp)
	}

	h.postProcess(ctx, req, err)
}

func (h *handler[Request, Response]) handleError(ctx context.Context, w http.ResponseWriter, req ParsedRequest[Request], err error) {
	// Send the error response
	if h.onErrorF != nil {
		h.onErrorF(ctx, w, err)
	}

	h.postProcess(ctx, req, err)
}

func (h *handler[Request, Response]) postProcess(ctx context.Context, req ParsedRequest[Request], err error) {
	if h.metricF != nil {
		h.metricF(ctx, req, err)
	}

	if h.auditF != nil {
		h.auditF(ctx, req, err)
	}
}
