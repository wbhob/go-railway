package railway

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

type contextKey struct{}

var headersKey contextKey

const (
	// HeaderRealIP for identifying client's remote IP.
	HeaderRealIP = "X-Real-IP"
	// HeaderForwardedProto always indicates https.
	HeaderForwardedProto = "X-Forwarded-Proto"
	// HeaderForwardedHost for identifying the original host header.
	HeaderForwardedHost = "X-Forwarded-Host"
	// HeaderRailwayEdge for identifying the edge region that handled the request.
	HeaderRailwayEdge = "X-Railway-Edge"
	// HeaderRequestStart for identifying the time the request was received (Unix milliseconds timestamp).
	HeaderRequestStart = "X-Request-Start"
	// HeaderRailwayRequestID for correlating requests against network logs.
	HeaderRailwayRequestID = "X-Railway-Request-Id"
)

// Headers are the HTTP request headers that are set by Railway. See https://docs.railway.com/reference/public-networking for more details.
type Headers struct {
	// X-Real-IP for identifying client's remote IP.
	RealIP string
	// X-Forwarded-Proto always indicates https.
	ForwardedProto string
	// X-Forwarded-Host for identifying the original host header.
	ForwardedHost string
	// X-Railway-Edge for identifying the edge region that handled the request.
	RailwayEdge string
	// X-Request-Start for identifying the time the request was received.
	RequestStart time.Time
	// X-Railway-Request-Id for correlating requests against network logs.
	RailwayRequestID string
}

// HeadersFromRequest parses the HTTP request headers and returns a Headers struct.
func HeadersFromRequest(r *http.Request) Headers {
	var requestStart time.Time
	if h := r.Header.Get(HeaderRequestStart); h != "" {
		requestStartInt, err := strconv.ParseInt(h, 10, 64)
		if err != nil {
			requestStart = time.Time{}
		}
		requestStart = time.UnixMilli(requestStartInt)
	}

	return Headers{
		RealIP:           r.Header.Get(HeaderRealIP),
		ForwardedProto:   r.Header.Get(HeaderForwardedProto),
		ForwardedHost:    r.Header.Get(HeaderForwardedHost),
		RailwayEdge:      r.Header.Get(HeaderRailwayEdge),
		RequestStart:     requestStart,
		RailwayRequestID: r.Header.Get(HeaderRailwayRequestID),
	}
}

// Handler returns an HTTP middleware that extracts Railway headers from the request
// and stores them in the request context for access by downstream handlers.
func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := HeadersFromRequest(r)
		ctx := context.WithValue(r.Context(), headersKey, headers)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HeadersFromContext extracts the Railway headers from the request context.
// Returns zero-value Headers if not found.
func HeadersFromContext(ctx context.Context) (Headers, bool) {
	if headers, ok := ctx.Value(headersKey).(Headers); ok {
		return headers, true
	}
	return Headers{}, false
}
