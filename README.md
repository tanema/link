# Link 🧝

[![Go Reference](https://pkg.go.dev/badge/github.com/tanema/link.svg)](https://pkg.go.dev/github.com/tanema/link)

Easily parse and create Link headers for pagination.

### Create Link Headers
Easily create link headers for responses to your paginated endpoints.

```go
func handler(w http.ResponseWriter, req http.Request) {
    // ...
    linkHeader := link.NewHeader(map[string]*url.URL{
        "first": firstURL,
        "last": lastURL,
        // if these are nil then they will not be included in the header
        "next": nextURL,
        "prev": prevURL,
    })
    w.WriteHeader("link", linkHeader.String())
    // ...
}
```

### Parse Link Headers for integrating with services.
Easily parse link headers from external services to paginate through results.

```go
func loadUsers(path string) {
    resp, err := http.DefaultClient.Get(path)
    linkHeader, err := link.Parse(resp)
    if next := linkHeader.Next(); next != nil {
        loadUsers(next.URL.String())
    }
}

func main() {
    loadUsers("http://service.com/users")
}
```
