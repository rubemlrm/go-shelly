package contracts

import (
	"context"
	"net/http"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

type ShellyClient interface {
	RetryHTTPCheck(ctx context.Context, resp *http.Response, err error) (bool, error)
	NewRequest(method, endpoint string, opts interface{}) (*retryablehttp.Request, error)
	ParseUrl(method, endpoint string, opts interface{}) (string, error)
	SetAdditionalHeaders(request *retryablehttp.Request, headers http.Header) error
	SetBasicAuth(request *retryablehttp.Request) error
	Do(req *retryablehttp.Request, v interface{}) (*http.Response, error)
}
