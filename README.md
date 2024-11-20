# handle-chain

This repository provides an utility to chain functions that ultimately create a http.Handler.

This handler can be used with any of the router libraries available in Go.

## Example

```go
import chain "github.com/go-dev-tools/handle-chain"


func () {

...

// Create Book Endpoint
http.Handle("/books", chain.New[Book, chain.EmptyResponse]().
    Parse(ReadCreateBookRequest).
    AuthZ(AllowBookWrite).
    Resolve(CreateBook).
    OnSuccess(OnBookCreated).
    OnError(OnFailedToCreateBook).
    Monitor(MonitorBookCreation).
    Audit(AuditBookCreation))
...

}
```

Read more [examples](./examples)