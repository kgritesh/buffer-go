package buffer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

const (
	libraryVersion = "0.1"
	baseURL        = "https://api.bufferapp.com/1/"
	userAgent      = "buffer-go/" + libraryVersion
	contentType    = "application/x-www-form-urlencoded"
)

// A Client manages communication with the Buffer API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	BaseURL   *url.URL
	UserAgent string

	UserService    *UserService
	ProfileService *ProfileService
	UpdateService  *UpdateService
}

/*
An ErrorResponse reports one or more errors caused by an API request.
*/
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"error"`
	Code     int            `json:"code"`
}

func (errorResponse *ErrorResponse) Error() string {
	return fmt.Sprintf("Error While Processing %v request at %v:\n"+
		"Status: %d, Error Code: %d, Error: %v",
		errorResponse.Response.Request.Method,
		sanitizeURL(errorResponse.Response.Request.URL),
		errorResponse.Response.StatusCode,
		errorResponse.Code, errorResponse.Message)

}

// NewClient returns a new GitHub API client.  If a nil httpClient is
// provided, http.DefaultClient will be used.  To use API methods which require
// authentication, provide an http.Client that will perform the authentication
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(baseURL)
	client := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	client.UserService = &UserService{client: client}
	client.ProfileService = &ProfileService{client: client}
	client.UpdateService = &UpdateService{client: client}
	return client
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (client *Client) NewRequest(method string, urlStr string, body interface{}) (*http.Request, error) {
	relative, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	reqURL := client.BaseURL.ResolveReference(relative)
	fmt.Println("Request Url ", reqURL)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, reqURL.String(), buf)
	if err != nil {
		return nil, err
	}
	if client.UserAgent != "" {
		req.Header.Add("User-Agent", client.UserAgent)
	}
	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (client *Client) Do(req *http.Request, result interface{}) (*http.Response, error) {
	resp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Response %v, Error %v", resp.Body, err)

	defer resp.Body.Close()

	err = CheckResponse(resp)

	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return resp, err
	}
	if result != nil {
		if w, ok := result.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(result)
		}
	}
	return resp, err

}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func CheckResponse(resp *http.Response) error {
	if c := resp.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: resp}
	data, err := ioutil.ReadAll(resp.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

// sanitizeURL redacts the client_id and client_secret tokens from the URL which
// may be exposed to the user, specifically in the ErrorResponse error message.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("client_secret")) > 0 {
		params.Set("client_secret", "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}

func main() {
	client := GetOauth2Client("1/020ba75e8e3c93b92516771bc42915ed")
	bufferClient := NewClient(client)
	user, resp, err := bufferClient.UserService.Get()
	fmt.Println(user, resp, err)
}
