package httpclient

import (
	"context"
	"net/http"
	"time"
)

const (
	// DefaultRetryTimes if request failed, can retry max times
	DefaultRetryTimes = 3
	// DefaultRetryDelay before retry, should wait delay time
	DefaultRetryDelay = time.Millisecond * 100
)

// RetryVerify parse the body and verify if it's correct
type RetryVerify func(body []byte) (shouldRetry bool)

func shouldRetry(ctx context.Context, httpCode int) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}
	switch httpCode {
	case
		_StatusReadRespErr,
		_StatusDoReqErr,

		http.StatusRequestTimeout,
		http.StatusLocked,
		http.StatusTooEarly,
		http.StatusTooManyRequests,

		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}
