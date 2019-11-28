package influx

import (
	"crypto/tls"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
	"net/http"
	"net/url"
	"time"
)

type ClientConfig struct {
	Addr               string
	Username           string
	Password           string
	Database           string
	RetentionPolicy    string
	UserAgent          string
	Timeout            time.Duration
	InsecureSkipVerify bool
	TLSConfig          *tls.Config
	Proxy              func(req *http.Request) (*url.URL, error)
}

type Client struct {
	client          client.Client
	dbName          string
	retentionPolicy string
	timeout         time.Duration
}

type BatchPointsWriter struct {
	batchPoints client.BatchPoints
	client      client.Client
}

type Tags map[string]string
type Fields map[string]interface{}

func NewClient(config ClientConfig) (c *Client, err error) {
	httpClient, err := newHTTPClient(config)
	if err != nil {
		return
	}

	c = &Client{
		client:          httpClient,
		dbName:          config.Database,
		retentionPolicy: config.RetentionPolicy,
		timeout:         config.Timeout,
	}

	return
}

func newHTTPClient(config ClientConfig) (c client.Client, err error) {
	httpConfig := client.HTTPConfig{
		Addr:               config.Addr,
		Username:           config.Username,
		Password:           config.Password,
		UserAgent:          config.UserAgent,
		Timeout:            config.Timeout,
		InsecureSkipVerify: config.InsecureSkipVerify,
		TLSConfig:          config.TLSConfig,
		Proxy:              config.Proxy,
	}

	c, err = client.NewHTTPClient(httpConfig)
	if err != nil {
		return
	}

	return
}

func (c *Client) Ping() (err error) {
	_, _, err = c.client.Ping(c.timeout)
	return
}

func (c *Client) WritePoints(name string, tags Tags, fields Fields, precision string) (err error) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision:        precision,
		Database:         c.dbName,
		RetentionPolicy:  c.retentionPolicy,
		WriteConsistency: "one",
	})

	if err != nil {
		return
	}

	pt, err := client.NewPoint(name, tags, fields, time.Now())
	if err != nil {
		return
	}

	bp.AddPoint(pt)

	if err = c.client.Write(bp); err != nil {
		return
	}

	return
}

func (c *Client) NewBatchPointsWriter(precision string) (bpw BatchPointsWriter, err error) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision:        precision,
		Database:         c.dbName,
		RetentionPolicy:  c.retentionPolicy,
		WriteConsistency: "one",
	})

	if err != nil {
		return
	}

	bpw.batchPoints = bp
	bpw.client = c.client

	return
}

func (b BatchPointsWriter) AddPoints(name string, tags Tags, fields Fields) {
	pt, err := client.NewPoint(name, tags, fields, time.Now())
	if err != nil {
		return
	}

	b.batchPoints.AddPoint(pt)
}

func (b BatchPointsWriter) Write() (err error) {
	err = b.client.Write(b.batchPoints)
	if err != nil {
		return
	}

	return
}
