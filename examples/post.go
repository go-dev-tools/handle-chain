package examples

import (
	"context"
	"net/http"
	"time"

	chain "github.com/go-dev-tools/handle-chain"
)

func init() {

	// Create book
	http.Handle("/books", chain.New[Book, chain.EmptyResponse]().
		Parse(ReadCreateBookRequest).
		AuthZ(AllowBookWrite).
		Resolve(CreateBook).
		OnSuccess(OnBookCreated).
		OnError(OnFailedToCreateBook).
		Monitor(MonitorBookCreation).
		Audit(AuditBookCreation))
}

type Book struct {
	Id          string
	Title       string
	Pages       int
	Language    string
	Author      string
	PublishedOn time.Time
}

func ReadCreateBookRequest(r *http.Request) (chain.ParsedRequest[Book], error) {

}

func AllowBookWrite(context.Context, chain.ParsedRequest[Book]) error {

}

func CreateBook(context.Context, chain.ParsedRequest[Book]) (chain.SuccessResponse[chain.EmptyResponse], error) {

}

func OnBookCreated(context.Context, http.ResponseWriter, chain.SuccessResponse[chain.EmptyResponse]) {

}

func OnFailedToCreateBook(context.Context, http.ResponseWriter, error) {

}

func MonitorBookCreation(context.Context, chain.ParsedRequest[Book], error) {

}

func AuditBookCreation(context.Context, chain.ParsedRequest[Book], error) {

}
