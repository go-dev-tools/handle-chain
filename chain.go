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

type ParseRequestFunc [Request any] func (*http.Request) (ParsedRequest[Request], error)

type AuthorizeFunc [Request any] func (context.Context, ParsedRequest[Request]) error

type ResolveRequestFunc [Request, Response any] func (context.Context, ParsedRequest[Request]) (Response, error)

type OnSuccessFunc [Response any] func (context.Context, http.ResponseWriter, Response)

type OnErrorFunc func (context.Context, http.ResponseWriter, error)

type CollectMetricFunc [Request any] func (context.Context, ParsedRequest[Request], error)

type AuditLogFunc [Request any] func (context.Context, ParsedRequest[Request], error)

type handler[Request, Response any] struct {
}

func (h *handler[Request, Response]) Parse(f ParseRequestFunc[Request]) Handler[Request, Response] {

}

func (h *handler[Request, Response]) AuthZ(f AuthorizeFunc[Request]) Handler[Request, Response] {

}

func (h *handler[Request, Response]) Resolve(f ResolveRequestFunc[Request, Response]) Handler[Request, Response] {

}

func (h *handler[Request, Response]) OnSuccess(f OnSuccessFunc[Response]) Handler[Request, Response] {

}

func (h *handler[Request, Response]) OnError(f OnErrorFunc) Handler[Request, Response] {

}

func (h *handler[Request, Response]) Monitor(f CollectMetricFunc[Request]) Handler[Request, Response] {

}

func (h *handler[Request, Response]) Audit(f AuditLogFunc[Request]) Handler[Request, Response] {

}

func (h *handler[Request, Response]) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
