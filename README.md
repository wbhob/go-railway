# go-railway

A Go package that provides type-safe access to Railway platform environment variables and HTTP request headers.

## Features

- **Type-safe**: Automatically converts environment variables to appropriate Go types (`int`, `int64`, `string`) and HTTP request headers to a `Headers` struct
- **Error handling**: Proper error handling for invalid values and non-Railway environments
- **Zero dependencies**: Uses only Go standard library
- **Complete coverage**: Includes all Railway-provided environment variables and HTTP request headers
- **Railway detection**: Built-in detection of Railway environment

## Installation

```bash
go get github.com/wbhob/go-railway
```

## Quick Start

### Environment Variables

```go
package main

import (
    "fmt"
    "log"

    "github.com/wbhob/go-railway"
)

func main() {
    // Check if running on Railway
    if !railway.IsRailway() {
        log.Fatal("Not running on Railway")
    }

    // Load environment variables with error handling
    env, err := railway.Load()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Project: %s\n", env.ProjectName)
    fmt.Printf("Service: %s\n", env.ServiceName)
    fmt.Printf("Environment: %s\n", env.EnvironmentName)

    // Access typed values directly
    if env.TCPProxyPort > 0 {
        fmt.Printf("TCP Proxy Port: %d\n", env.TCPProxyPort)
    }
}
```

### HTTP Middleware

```go
package main

import (
    "fmt"
    "net/http"
    "log"

    "github.com/wbhob/go-railway"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // Get Railway headers from context (set by middleware)
    headers, ok := railway.HeadersFromContext(r.Context())
    if !ok {
        log.Printf("Railway headers not found in context")
        return
    }

    fmt.Printf("Client IP: %s\n", headers.RealIP)
    fmt.Printf("Edge Region: %s\n", headers.RailwayEdge)
    fmt.Printf("Request ID: %s\n", headers.RailwayRequestID)

    w.WriteHeader(http.StatusOK)
}

func main() {
    // Wrap your handler with Railway middleware
    http.Handle("/", railway.Handler(http.HandlerFunc(handler)))

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Direct Header Parsing

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Extract Railway headers directly from the request
    headers := railway.HeadersFromRequest(r)

    fmt.Printf("Client IP: %s\n", headers.RealIP)
    // ... rest stays the same
}
```

## Must Functions

For scenarios where you want to panic on errors (e.g., during application startup):

```go
// Panics if not running on Railway or if parsing fails
env := railway.MustLoad()

// Or use with package-level variables
var env = railway.MustLoad()
```

## Examples

### Web Server with Railway Detection and Headers

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/wbhob/go-railway"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // Get Railway headers from context
    headers, _ := railway.HeadersFromContext(r.Context())

    // Log request details
    log.Printf("Request from %s via edge %s (ID: %s)",
        headers.RealIP, headers.RailwayEdge, headers.RailwayRequestID)

    fmt.Fprintf(w, "Hello from Railway! Your IP: %s\n", headers.RealIP)
}

func main() {
    var port string = "8080" // default

    if railway.IsRailway() {
        env, err := railway.Load()
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("Running on Railway: %s/%s\n", env.ProjectName, env.ServiceName)

        if railwayPort := os.Getenv("PORT"); railwayPort != "" {
            port = railwayPort
        }
    }

    // Use Railway middleware to automatically parse headers into context
    http.Handle("/", railway.Handler(http.HandlerFunc(handler)))

    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

### Middleware for Request Logging

```go
func logRequest(ctx context.Context) {
    headers, ok := railway.HeadersFromContext(ctx)
    if !ok {
        log.Printf("No Railway headers found")
        return
    }

    log.Printf("Processing request %s from %s",
        headers.RailwayRequestID, headers.RealIP)
}

func handler(w http.ResponseWriter, r *http.Request) {
    // Headers are available throughout the request lifecycle
    logRequest(r.Context())

    // Your business logic here
    w.WriteHeader(http.StatusOK)
}

func main() {
    http.Handle("/", railway.Handler(http.HandlerFunc(handler)))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Related

- [Railway Documentation](https://docs.railway.com/reference/variables#railway-provided-variables)
- [Railway Public Networking](https://docs.railway.com/reference/public-networking)
- [Railway Platform](https://railway.app)
