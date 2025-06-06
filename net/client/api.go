// Package netclient is the package that can use to send request
package netclient

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

// AIRequest is a struct that represents an AI request. It includes all the necessary fields to make a POST request to an AI API. The `SendRequest` method sends the request and returns the response or an error.
type AIRequest struct {
	Method      string
	URL         string
	Headers     map[string]string // 可放 Authorization、X-API-Key
	Cookies     map[string]string // 可放 session cookie 等
	QueryParams map[string]string //
	Body        string
	ContentType string
	TimeoutMs   time.Duration
	MaxRetries  int
	RetryDelay  time.Duration
}

// SendRequest sends the AI request and returns the response or an error.
func (r *AIRequest) SendRequest() (*http.Response, error) {

	client := &http.Client{
		Timeout: r.TimeoutMs,
	}

	// rawUrl can have query params in it. We need to parse them and add them to the request URL.
	parsedURL, err := url.Parse(r.URL)
	if err != nil {
		return nil, err
	}

	query := parsedURL.Query()
	for k, v := range r.QueryParams {
		query.Set(k, v)
	}

	parsedURL.RawQuery = query.Encode()

	targetURL := parsedURL.String()

	req, err := http.NewRequest(r.Method, targetURL, strings.NewReader(r.Body))

	if err != nil {
		return nil, err
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	// Set content type seperately if provided
	if r.ContentType != "" {
		req.Header.Set("Content-Type", r.ContentType)
	}

	for key, value := range r.Cookies {
		req.AddCookie(&http.Cookie{Name: key, Value: value})
	}

	// Send the request and get the response
	var resp *http.Response

	for range r.MaxRetries {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		time.Sleep(r.RetryDelay)
	}

	return resp, err
}
