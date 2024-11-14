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

	Parse(func(*http.Request) (Request, error)) Handler[Request, Response]

	AuthZ(func(context.Context, Request) error) Handler[Request, Response]

	Resolve(func(context.Context, Request) (Response, error)) Handler[Request, Response]

	OnSuccess(func(context.Context, Response)) Handler[Request, Response]

	OnError(func(context.Context, error)) Handler[Request, Response]

	Monitor(func(context.Context)) Handler[Request, Response]

	Audit(func(context.Context)) Handler[Request, Response]
}

type handler[Request, Response any] struct {
}

func (h *handler[Request, Response]) Parse(f func(*http.Request) (Request, error)) Handler[Request, Response] {

}

func (h *handler[Request, Response]) AuthZ(f func(context.Context, Request) error) Handler[Request, Response] {

}

func (h *handler[Request, Response]) Resolve(f func(context.Context, Request) (Response, error)) Handler[Request, Response] {

}

func (h *handler[Request, Response]) OnSuccess(f func(context.Context, Response)) Handler[Request, Response] {

}

func (h *handler[Request, Response]) OnError(f func(context.Context, error)) Handler[Request, Response] {

}

func (h *handler[Request, Response]) Monitor(f func(context.Context)) Handler[Request, Response] {

}

func (h *handler[Request, Response]) Audit(f func(context.Context)) Handler[Request, Response] {

}

func (h *handler[Request, Response]) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
