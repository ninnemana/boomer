package api

import (
	"fmt"
	"github.com/ninnemana/boom/boomer"
	"net/http"
	"net/url"
	"regexp"
	"runtime"
	"strings"
)

const (
	headerRegexp = "^([\\w-]+):\\s*(.+)"
	authRegexp   = "^([\\w-\\.]+):(.+)"
)

// Request ...
type Request struct {
	// url the make the requests against
	Endpoint string `json:"endpoint"`

	// number of requests to run
	Requests int `json:"requests"`

	// number of requests to run concurrently. Total number
	// of requests cannot be smaler than the concurrency level.
	ConcurrentRequests int `json:"concurrent_requests"`

	// Rate limit, in seconds (QPS)
	RateLimit int `json:"rate_limit"`

	// HTTP Method - GET POST PUT DELETE HEAD OPTIONS
	Method string `json:"method"`

	// HTTP Headers
	Headers []HTTPHeader `json:"headers"`

	// Timeout - in millieseconds
	Timeout int `json:"timeout"`

	// HTTP Accept Header
	AcceptHeader string `json:"accept"`

	// HTTP request body
	Body string `json:"body"`

	// Content-Type - defaults to "text/html"
	ContentType string `json:"content_type"`

	// Basic Authentication - username:password
	BasicAuth string `json:"basic_auth"`

	// HTTP Proxy address - host:port
	Proxy string `json:"proxy"`

	// Consumes the entire request body
	ReadAll bool `json:"read_all"`

	// Allow bad/expires TLS/SSL certificates
	AllowInsecure bool `json:"allow_insecure"`

	// Disable compression
	DisableCompression bool `json:"disable_compression"`

	// Disable keep-alive, prevents re-use of TCP connections
	// between different HTTP requests.
	DisableKeepAlive bool `json:"disable_keep_alive"`

	// NUmber of used CPU cores. (default for current machine is %d cores)
	CPUs int `json:"cpus"`

	Report *boomer.Report `json:"report"`

	username string
	password string
	proxyURL *url.URL
}

// HTTPHeader ...
type HTTPHeader struct {
	// Name - "Content-Type"
	Name string `json:"name"`

	// Value - "application/json"
	Value string `json:"value"`
}

func (r *Request) Do() error {
	req, err := http.NewRequest(r.Method, r.Endpoint, nil)
	if err != nil {
		return err
	}

	for _, h := range r.Headers {
		req.Header.Add(h.Name, h.Value)
	}
	if r.username != "" && r.password != "" {
		req.SetBasicAuth(r.username, r.password)
	}
	bmr := &boomer.Boomer{
		Request:            req,
		N:                  r.Requests,
		C:                  r.ConcurrentRequests,
		RequestBody:        r.Body,
		Qps:                r.RateLimit,
		Timeout:            r.Timeout,
		AllowInsecure:      r.AllowInsecure,
		DisableCompression: r.DisableCompression,
		DisableKeepAlives:  r.DisableKeepAlive,
		ProxyAddr:          r.proxyURL,
		Output:             "csv",
		ReadAll:            r.ReadAll,
	}
	r.Report = bmr.Run()

	return nil
}

func (r *Request) validate() error {

	if r.CPUs == 0 {
		r.CPUs = runtime.GOMAXPROCS(-1)
	}
	runtime.GOMAXPROCS(r.CPUs)

	if r.Requests <= 0 || r.ConcurrentRequests <= 0 {
		return fmt.Errorf("%s", "requests and concurrent requests cannot be smaller than 1.")
	}

	header := make(http.Header)

	r.Method = strings.ToUpper(r.Method)

	// set content-type
	if r.ContentType == "" {
		r.ContentType = "text/html"
	}
	header.Set("Content-Type", r.ContentType)

	// set any other additional headers
	for _, h := range r.Headers {
		if parseInputWithRegexp(h.Name, headerRegexp) != nil || parseInputWithRegexp(h.Value, headerRegexp) != nil {
			continue
		}

		header.Set(h.Name, h.Value)
	}

	if r.AcceptHeader != "" {
		header.Set("Accept", r.AcceptHeader)
	}

	// set basic auth if set
	if r.BasicAuth != "" {
		err := parseInputWithRegexp(r.BasicAuth, authRegexp)
		if err != nil {
			return err
		}
		match := strings.Split(r.BasicAuth, ":")
		r.username, r.password = match[0], match[1]
	}

	if r.Proxy != "" {
		var err error
		r.proxyURL, err = url.Parse(r.Proxy)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseInputWithRegexp(input, regx string) error {
	re := regexp.MustCompile(regx)
	matches := re.FindStringSubmatch(input)

	if len(matches) < 1 {
		return fmt.Errorf("%s", "Could not parse provided input")
	}

	return nil
}
