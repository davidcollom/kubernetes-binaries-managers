package versions

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/little-angry-clouds/kubernetes-binaries-managers/internal/logging"
)

// backoffHandler implements backoff logic that respects rate limit headers.
func backoffHandler(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	// If a rate limit header is present, prioritize it over standard backoff logic.
	if resp != nil {
		if remainingHeader := resp.Header.Get("X-Ratelimit-Remaining"); remainingHeader == "0" {
			if resetHeader := resp.Header.Get("X-Ratelimit-Reset"); resetHeader != "" {
				resetTimestamp, err := strconv.ParseInt(resetHeader, 10, 64)
				if err == nil {
					// Calculate how long to wait until the reset time.
					now := time.Now().UTC().Unix()

					sleepDuration := time.Duration(resetTimestamp-now) * time.Second
					if sleepDuration > 0 {
						log.Printf("Rate limit exceeded. Sleeping for %v until reset.", sleepDuration)
						return sleepDuration
					}
				}
			}
		}
	}

	// Fall back to the default exponential backoff if no rate limit headers are found.
	return retryablehttp.DefaultBackoff(min, max, attemptNum, resp)
}

// retryPolicy extends the default policy to retry on specific status codes.
func retryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// Use the default policy as a base.
	shouldRetry, retryErr := retryablehttp.DefaultRetryPolicy(ctx, resp, err)

	// Add custom logic for rate limiting status codes.
	if resp != nil {
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
			// Only retry if the rate limit headers suggests a future reset.
			if resp.Header.Get("X-Ratelimit-Remaining") == "0" {
				return true, nil
			}
		}
	}

	return shouldRetry, retryErr
}

// AuthRoundTripper adds an Authorization header to all requests.
type AuthRoundTripper struct {
	token            string
	nextRoundTripper http.RoundTripper
}

func (art *AuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// If we have a token and one is not already set, lets set it!
	if art.token != "" && req.Header.Get("Authorization") == "" {
		logging.Debug("adding Authorization header to request")
		// Add the Authorization header to the request.
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", art.token))
	}

	// Execute the request using the next RoundTripper in the chain.
	return art.nextRoundTripper.RoundTrip(req)
}
