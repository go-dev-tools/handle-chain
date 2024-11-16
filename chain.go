package chain

import (
	"context"
	"net/http"
)

func New[Request, Response any]() Handler[Request, Response] {
	return &handler[Request, Response]{}
}

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
	Headers     map[string]string
	PathParams  map[string]string
	QueryParams map[string]string
}

type ParseRequestFunc[Request any] func(*http.Request) (ParsedRequest[Request], error)

type AuthorizeFunc[Request any] func(context.Context, ParsedRequest[Request]) error

type ResolveRequestFunc[Request, Response any] func(context.Context, ParsedRequest[Request]) (Response, error)

type OnSuccessFunc[Response any] func(context.Context, http.ResponseWriter, Response)

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
	
  req, err := h.parseReqF(r)
  if err != nil {
	h.handleError(ctx, w, req, err)
    return
  }

  err = h.authF(ctx, req)
  if err != nil {
    h.handleError(ctx, w, req, err)
	return
  }

  resp, err := h.resolveF(ctx, req)
  if err != nil {
	h.handleError(ctx, w, req, err)
	return
  }

  h.onSuccessF(ctx, w, resp)
  h.metricF(ctx, req, err)
  h.auditF(ctx, req, err)
}

func (h *handler[Request, Response]) handleError(ctx context.Context, w http.ResponseWriter, req ParsedRequest[Request], err error) {
	h.onErrorF(ctx, w,err)
	h.metricF(ctx,req, err)
	h.auditF(ctx, req, err)
}