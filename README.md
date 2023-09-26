# go-altair
Altair GraphQL Client Go HTTP Handler

[![Lint](https://github.com/hamidfzm/go-altair/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/hamidfzm/go-altair/actions/workflows/golangci-lint.yml)
[![Test](https://github.com/hamidfzm/go-altair/actions/workflows/test.yml/badge.svg)](https://github.com/hamidfzm/go-altair/actions/workflows/test.yml)

## Installation

```bash
go get github.com/hamidfzm/go-altair
```

## Usage

```go
package main

import (
	"net/http"
)

func main() {
	config := &altair.Config{
		DefaultWindowTitle: "Mafia Altair",
		Endpoint:           endpoint,
		Force:              false,
		Headers: []altair.Header{
			{
				Key:   echo.HeaderAuthorization,
				Value: "Bearer <token>",
			},
		},
	}
	http.Handle("/altair", altair.Handler(config))
	http.ListenAndServe(":8080", nil)
}

```

## License

Check [LICENSE](LICENSE) file for more information.
