package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func GetRetriableClient(retryMaxCount int, timeout time.Duration, checkRetry retryablehttp.CheckRetry) *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 30 * time.Second
	retryClient.RetryMax = retryMaxCount
	retryClient.HTTPClient.Timeout = timeout
	if checkRetry != nil {
		retryClient.CheckRetry = checkRetry
	} else {
		retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
			// skip retries for non 5xx errors
			if resp != nil && resp.StatusCode < http.StatusInternalServerError {
				return false, err
			}
			return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
		}
	}
	return retryClient.StandardClient()
}
