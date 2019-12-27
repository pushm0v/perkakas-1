package httpclient

import (
	"time"

	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/httpclient"
)

type HttpClient struct {
	Client *httpclient.Client
}

type HttpClientConf struct {
	BackoffInterval       time.Duration
	MaximumJitterInterval time.Duration
	Timeout               time.Duration
	RetryCount            int
}

func NewHttpClient(conf *HttpClientConf) *HttpClient {

	if conf == nil {
		// Default configuration
		conf = new(HttpClientConf)
		conf.BackoffInterval = 2 * time.Millisecond       // 2ms
		conf.MaximumJitterInterval = 5 * time.Millisecond // 5ms
		conf.Timeout = 15000 * time.Millisecond           // 15s
		conf.RetryCount = 3                               // 3 times
	}

	backoff := heimdall.NewConstantBackoff(conf.BackoffInterval, conf.MaximumJitterInterval)
	retrier := heimdall.NewRetrier(backoff)

	newClient := httpclient.NewClient(
		httpclient.WithHTTPTimeout(conf.Timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(conf.RetryCount),
	)

	return &HttpClient{
		Client: newClient,
	}
}
