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
		Resolve(InsertBook).
		OnSuccess(OnBookInserted).
		OnError(OnFailedToInsertBook).
		Monitor(MonitorBookInsertion).
		Audit(AuditBookInsertion))
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

func InsertBook(context.Context, chain.ParsedRequest[Book]) (chain.SuccessResponse[chain.EmptyResponse], error) {

}

func OnBookInserted(context.Context, http.ResponseWriter, chain.SuccessResponse[chain.EmptyResponse]) {

}

func OnFailedToInsertBook(context.Context, http.ResponseWriter, error) {

}

func MonitorBookInsertion(context.Context, chain.ParsedRequest[Book], error) {

}

func AuditBookInsertion(context.Context, chain.ParsedRequest[Book], error) {

}
