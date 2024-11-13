package chain

import (
	"context"
	"net/http"
)

func New [Request, Response any] () Handler[Request, Response] {

}

type Handler [Request, Response any] interface {
	http.Handler

    Parse(func(*http.Request) (Request, error)) Handler[Request, Response]

	AuthZ(func(context.Context, Request) error) Handler[Request, Response]

    Resolve(func(context.Context, Request) (Response, error)) Handler[Request, Response]

    OnSuccess(func (context.Context, Response)) Handler[Request, Response]

	OnError(func (context.Context, error)) Handler[Request, Response]
}