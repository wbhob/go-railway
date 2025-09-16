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

### HTTP Headers

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/wbhob/go-railway"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // Extract Railway headers from the request
    headers := railway.HeadersFromRequest(r)

    fmt.Printf("Client IP: %s\n", headers.RealIP)
    fmt.Printf("Edge Region: %s\n", headers.RailwayEdge)
    fmt.Printf("Request ID: %s\n", headers.RailwayRequestID)

    w.WriteHeader(http.StatusOK)
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
    // Extract Railway headers
    headers := railway.HeadersFromRequest(r)

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

        // Use Railway's PORT variable for the server
        if railwayPort := os.Getenv("PORT"); railwayPort != "" {
            port = railwayPort
        }
    }

    http.HandleFunc("/", handler)

    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

### Middleware for Request Logging

```go
func railwayLoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        headers := railway.HeadersFromRequest(r)

        // Log with Railway-specific context
        log.Printf("[%s] %s %s - IP: %s, Edge: %s",
            headers.RailwayRequestID,
            r.Method,
            r.URL.Path,
            headers.RealIP,
            headers.RailwayEdge,
        )

        next.ServeHTTP(w, r)
    })
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
