package lib

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"io"
	"net/http"
)

// Ping a port, return span if ctx contains a Span.
func Ping(ctx context.Context, port string) (string, error) {
	url := fmt.Sprintf("http://%s/ping", port)

	// Create a Span using the given Context.
	span, _ := opentracing.StartSpanFromContext(ctx, "ping-send")
	defer span.Finish()

	// Generate a new Get Request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	// Inject the Span into the Request.
	if err := Inject(span, req); err != nil {
		return "", err
	}

	// Call Service-Two Ping HTTP Handler
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Assert http.StatusOk
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("StatusCode: %d, Body: %s", resp.StatusCode, body)
	}

	return string(body), nil
}


