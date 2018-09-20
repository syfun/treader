package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

var (
	apiBase     = "http://api.zhuishushenqi.com"
	chapterBase = "http://chapter2.zhuishushenqi.com"
)

// Client with host.
type Client struct {
	APIBaseURL     *url.URL
	ChapterBaseURL *url.URL
	UserAgent      string
	client         *http.Client
}

// New create new client.
func New() *Client {
	apiBaseURL, _ := url.Parse(apiBase)
	chapterBaseURL, _ := url.Parse(chapterBase)
	return &Client{apiBaseURL, chapterBaseURL, "go-client", http.DefaultClient}
}

// DefaultClient represents default client.
var DefaultClient = New()

// NewAPIRequest creates an get method request with param values.
// urlStr should be relative url with a preceding slash.
func (c *Client) NewAPIRequest(urlStr string, query url.Values) (*http.Request, error) {
	u, err := c.APIBaseURL.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot new api request")
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "cannot new api request")
	}

	req.Header.Set("Accept", "application/json")
	// req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

// NewChatperRequest creates an get method request with param values.
// urlStr should be relative url with a preceding slash.
func (c *Client) NewChatperRequest(urlStr string, query url.Values) (*http.Request, error) {
	u, err := c.ChapterBaseURL.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot new chapter request")
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "cannot new chapter request")
	}

	// req.Header.Set("Accept", "application/json")
	// req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, errors.Wrap(ctx.Err(), "cannot do request")
		default:
		}

		return nil, errors.Wrap(err, "cannot do request")
	}

	defer resp.Body.Close()
	response := newResponse(resp)

	if err := CheckResponse(resp); err != nil {
		return response, err
	}
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = errors.Wrap(err, "cannot do request")
			}
		}
	}

	return response, err
}

// Response wraps the standard http.Response.
type Response struct {
	*http.Response
}

// newResponse creates a new Response for the provided http.Response.
// r must not be nil.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// ErrorResponse reports one or more errors casued by api request.
type ErrorResponse struct {
	Response *http.Response
	Code     int    `json:"code"`
	Message  string `json:"msg"`
	Success  bool   `json:"ok"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Code, r.Message)
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range or equal to 202 Accepted.
// API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other
// response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}
