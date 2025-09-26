package versions

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

func TestBackoffHandler(t *testing.T) {
	tests := []struct {
		name           string
		min            time.Duration
		max            time.Duration
		attemptNum     int
		statusCode     int
		rateLimitRem   string
		rateLimitReset string
		wantSleep      time.Duration
	}{
		{
			name:       "No response, fallback to default backoff",
			min:        1 * time.Second,
			max:        10 * time.Second,
			attemptNum: 2,
			statusCode: 0,
			wantSleep:  retryablehttp.DefaultBackoff(1*time.Second, 10*time.Second, 2, nil),
		},
		{
			name:       "No rate limit headers, fallback to default backoff",
			min:        1 * time.Second,
			max:        10 * time.Second,
			attemptNum: 1,
			statusCode: http.StatusOK,
			wantSleep:  retryablehttp.DefaultBackoff(1*time.Second, 10*time.Second, 1, &http.Response{StatusCode: http.StatusOK}),
		},
		{
			name:           "Rate limit exceeded, valid reset header in future",
			min:            1 * time.Second,
			max:            10 * time.Second,
			attemptNum:     1,
			statusCode:     http.StatusTooManyRequests,
			rateLimitRem:   "0",
			rateLimitReset: strconv.FormatInt(time.Now().UTC().Unix()+5, 10),
			wantSleep:      5 * time.Second,
		},
		{
			name:           "Rate limit exceeded, reset header in past",
			min:            1 * time.Second,
			max:            10 * time.Second,
			attemptNum:     1,
			statusCode:     http.StatusTooManyRequests,
			rateLimitRem:   "0",
			rateLimitReset: strconv.FormatInt(time.Now().UTC().Unix()-10, 10),
			wantSleep: retryablehttp.DefaultBackoff(1*time.Second, 10*time.Second, 1, &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Header:     http.Header{"X-Ratelimit-Remaining": []string{"0"}, "X-Ratelimit-Reset": []string{strconv.FormatInt(time.Now().UTC().Unix()-10, 10)}},
			}),
		},
		{
			name:           "Rate limit exceeded, invalid reset header",
			min:            1 * time.Second,
			max:            10 * time.Second,
			attemptNum:     1,
			statusCode:     http.StatusTooManyRequests,
			rateLimitRem:   "0",
			rateLimitReset: "invalid",
			wantSleep: retryablehttp.DefaultBackoff(1*time.Second, 10*time.Second, 1, &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Header:     http.Header{"X-Ratelimit-Remaining": []string{"0"}, "X-Ratelimit-Reset": []string{"invalid"}},
			}),
		},
		{
			name:           "Rate limit not exceeded, fallback to default backoff",
			min:            1 * time.Second,
			max:            10 * time.Second,
			attemptNum:     1,
			statusCode:     http.StatusTooManyRequests,
			rateLimitRem:   "5",
			rateLimitReset: strconv.FormatInt(time.Now().UTC().Unix()+10, 10),
			wantSleep: retryablehttp.DefaultBackoff(1*time.Second, 10*time.Second, 1, &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Header:     http.Header{"X-Ratelimit-Remaining": []string{"5"}, "X-Ratelimit-Reset": []string{strconv.FormatInt(time.Now().UTC().Unix()+10, 10)}},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *http.Response
			if tt.statusCode != 0 {
				header := http.Header{}
				if tt.rateLimitRem != "" {
					header.Set("X-Ratelimit-Remaining", tt.rateLimitRem)
				}
				if tt.rateLimitReset != "" {
					header.Set("X-Ratelimit-Reset", tt.rateLimitReset)
				}
				resp = &http.Response{
					StatusCode: tt.statusCode,
					Header:     header,
				}
			}
			got := backoffHandler(tt.min, tt.max, tt.attemptNum, resp)
			// Allow a small delta for timing differences
			if tt.wantSleep > 0 && tt.wantSleep < 10*time.Second {
				assert.InDelta(t, tt.wantSleep.Seconds(), got.Seconds(), 1.0)
			} else {
				assert.Equal(t, tt.wantSleep, got)
			}
		})
	}
}

func TestRetryPolicy_Table(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		rateLimitRem string
		expectRetry  bool
		expectErr    bool
	}{
		{
			name:         "OK status, no rate limit header",
			statusCode:   http.StatusOK,
			rateLimitRem: "",
			expectRetry:  false,
			expectErr:    false,
		},
		{
			name:         "Forbidden, rate limit remaining 0",
			statusCode:   http.StatusForbidden,
			rateLimitRem: "0",
			expectRetry:  true,
			expectErr:    false,
		},
		{
			name:         "Forbidden, rate limit remaining 1",
			statusCode:   http.StatusForbidden,
			rateLimitRem: "1",
			expectRetry:  false,
			expectErr:    false,
		},
		{
			name:         "TooManyRequests, rate limit remaining 0",
			statusCode:   http.StatusTooManyRequests,
			rateLimitRem: "0",
			expectRetry:  true,
			expectErr:    false,
		},
		{
			name:         "TooManyRequests, rate limit remaining 5",
			statusCode:   http.StatusTooManyRequests,
			rateLimitRem: "5",
			expectRetry:  false,
			expectErr:    false,
		},
		{
			name:         "InternalServerError, no rate limit header",
			statusCode:   http.StatusInternalServerError,
			rateLimitRem: "",
			expectRetry:  true, // DefaultRetryPolicy retries 500
			expectErr:    false,
		},
		{
			name:         "BadGateway, no rate limit header",
			statusCode:   http.StatusBadGateway,
			rateLimitRem: "",
			expectRetry:  true, // DefaultRetryPolicy retries 502
			expectErr:    false,
		},
		{
			name:         "Forbidden, no rate limit header",
			statusCode:   http.StatusForbidden,
			rateLimitRem: "",
			expectRetry:  false,
			expectErr:    false,
		},
		{
			name:         "TooManyRequests, no rate limit header",
			statusCode:   http.StatusTooManyRequests,
			rateLimitRem: "",
			expectRetry:  false,
			expectErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := http.Header{}
			if tt.rateLimitRem != "" {
				header.Set("X-Ratelimit-Remaining", tt.rateLimitRem)
			}
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Header:     header,
			}
			shouldRetry, err := retryPolicy(context.Background(), resp, nil)
			assert.Equal(t, tt.expectRetry, shouldRetry)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
