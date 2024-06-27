// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	HTTPBasicScopes     = "HTTPBasic.Scopes"
	OpenIdConnectScopes = "OpenIdConnect.Scopes"
)

// AWSCredentials defines model for AWSCredentials.
type AWSCredentials struct {
	AccessKeyId     string    `json:"AccessKeyId"`
	Expiration      time.Time `json:"Expiration"`
	SecretAccessKey string    `json:"SecretAccessKey"`
	SessionToken    string    `json:"SessionToken"`
}

// HTTPValidationError defines model for HTTPValidationError.
type HTTPValidationError struct {
	Detail *[]ValidationError `json:"detail,omitempty"`
}

// User defines model for User.
type User struct {
	Email      openapi_types.Email `json:"email"`
	ExternalId string              `json:"external_id"`
	FirstName  string              `json:"first_name"`
	Id         openapi_types.UUID  `json:"id"`
	LastName   string              `json:"last_name"`
}

// ValidationError defines model for ValidationError.
type ValidationError struct {
	Loc  []ValidationError_Loc_Item `json:"loc"`
	Msg  string                     `json:"msg"`
	Type string                     `json:"type"`
}

// ValidationErrorLoc0 defines model for .
type ValidationErrorLoc0 = string

// ValidationErrorLoc1 defines model for .
type ValidationErrorLoc1 = int

// ValidationError_Loc_Item defines model for ValidationError.loc.Item.
type ValidationError_Loc_Item struct {
	union json.RawMessage
}

// GetAwsCredentialsApiAuthAwsCredentialsGetParams defines parameters for GetAwsCredentialsApiAuthAwsCredentialsGet.
type GetAwsCredentialsApiAuthAwsCredentialsGetParams struct {
	DataProductName string `form:"data_product_name" json:"data_product_name"`
	Environment     string `form:"environment" json:"environment"`
}

// GetDeviceTokenApiAuthDeviceDeviceTokenPostParams defines parameters for GetDeviceTokenApiAuthDeviceDeviceTokenPost.
type GetDeviceTokenApiAuthDeviceDeviceTokenPostParams struct {
	ClientId string  `form:"client_id" json:"client_id"`
	Scope    *string `form:"scope,omitempty" json:"scope,omitempty"`
}

// GetJwtTokenApiAuthDeviceJwtTokenPostParams defines parameters for GetJwtTokenApiAuthDeviceJwtTokenPost.
type GetJwtTokenApiAuthDeviceJwtTokenPostParams struct {
	ClientId   string `form:"client_id" json:"client_id"`
	DeviceCode string `form:"device_code" json:"device_code"`
	GrantType  string `form:"grant_type" json:"grant_type"`
}

// AsValidationErrorLoc0 returns the union data inside the ValidationError_Loc_Item as a ValidationErrorLoc0
func (t ValidationError_Loc_Item) AsValidationErrorLoc0() (ValidationErrorLoc0, error) {
	var body ValidationErrorLoc0
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromValidationErrorLoc0 overwrites any union data inside the ValidationError_Loc_Item as the provided ValidationErrorLoc0
func (t *ValidationError_Loc_Item) FromValidationErrorLoc0(v ValidationErrorLoc0) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeValidationErrorLoc0 performs a merge with any union data inside the ValidationError_Loc_Item, using the provided ValidationErrorLoc0
func (t *ValidationError_Loc_Item) MergeValidationErrorLoc0(v ValidationErrorLoc0) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsValidationErrorLoc1 returns the union data inside the ValidationError_Loc_Item as a ValidationErrorLoc1
func (t ValidationError_Loc_Item) AsValidationErrorLoc1() (ValidationErrorLoc1, error) {
	var body ValidationErrorLoc1
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromValidationErrorLoc1 overwrites any union data inside the ValidationError_Loc_Item as the provided ValidationErrorLoc1
func (t *ValidationError_Loc_Item) FromValidationErrorLoc1(v ValidationErrorLoc1) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeValidationErrorLoc1 performs a merge with any union data inside the ValidationError_Loc_Item, using the provided ValidationErrorLoc1
func (t *ValidationError_Loc_Item) MergeValidationErrorLoc1(v ValidationErrorLoc1) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t ValidationError_Loc_Item) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *ValidationError_Loc_Item) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetAwsCredentialsApiAuthAwsCredentialsGet request
	GetAwsCredentialsApiAuthAwsCredentialsGet(ctx context.Context, params *GetAwsCredentialsApiAuthAwsCredentialsGetParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetDeviceTokenApiAuthDeviceDeviceTokenPost request
	GetDeviceTokenApiAuthDeviceDeviceTokenPost(ctx context.Context, params *GetDeviceTokenApiAuthDeviceDeviceTokenPostParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetJwtTokenApiAuthDeviceJwtTokenPost request
	GetJwtTokenApiAuthDeviceJwtTokenPost(ctx context.Context, params *GetJwtTokenApiAuthDeviceJwtTokenPostParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AuthorizeApiAuthUserGet request
	AuthorizeApiAuthUserGet(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetAwsCredentialsApiAuthAwsCredentialsGet(ctx context.Context, params *GetAwsCredentialsApiAuthAwsCredentialsGetParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAwsCredentialsApiAuthAwsCredentialsGetRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetDeviceTokenApiAuthDeviceDeviceTokenPost(ctx context.Context, params *GetDeviceTokenApiAuthDeviceDeviceTokenPostParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetDeviceTokenApiAuthDeviceDeviceTokenPostRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetJwtTokenApiAuthDeviceJwtTokenPost(ctx context.Context, params *GetJwtTokenApiAuthDeviceJwtTokenPostParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetJwtTokenApiAuthDeviceJwtTokenPostRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AuthorizeApiAuthUserGet(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAuthorizeApiAuthUserGetRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetAwsCredentialsApiAuthAwsCredentialsGetRequest generates requests for GetAwsCredentialsApiAuthAwsCredentialsGet
func NewGetAwsCredentialsApiAuthAwsCredentialsGetRequest(server string, params *GetAwsCredentialsApiAuthAwsCredentialsGetParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/auth/aws_credentials")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "data_product_name", runtime.ParamLocationQuery, params.DataProductName); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "environment", runtime.ParamLocationQuery, params.Environment); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetDeviceTokenApiAuthDeviceDeviceTokenPostRequest generates requests for GetDeviceTokenApiAuthDeviceDeviceTokenPost
func NewGetDeviceTokenApiAuthDeviceDeviceTokenPostRequest(server string, params *GetDeviceTokenApiAuthDeviceDeviceTokenPostParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/auth/device/device_token")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "client_id", runtime.ParamLocationQuery, params.ClientId); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		if params.Scope != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "scope", runtime.ParamLocationQuery, *params.Scope); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetJwtTokenApiAuthDeviceJwtTokenPostRequest generates requests for GetJwtTokenApiAuthDeviceJwtTokenPost
func NewGetJwtTokenApiAuthDeviceJwtTokenPostRequest(server string, params *GetJwtTokenApiAuthDeviceJwtTokenPostParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/auth/device/jwt_token")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "client_id", runtime.ParamLocationQuery, params.ClientId); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "device_code", runtime.ParamLocationQuery, params.DeviceCode); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "grant_type", runtime.ParamLocationQuery, params.GrantType); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewAuthorizeApiAuthUserGetRequest generates requests for AuthorizeApiAuthUserGet
func NewAuthorizeApiAuthUserGetRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/auth/user")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetAwsCredentialsApiAuthAwsCredentialsGetWithResponse request
	GetAwsCredentialsApiAuthAwsCredentialsGetWithResponse(ctx context.Context, params *GetAwsCredentialsApiAuthAwsCredentialsGetParams, reqEditors ...RequestEditorFn) (*GetAwsCredentialsApiAuthAwsCredentialsGetResponse, error)

	// GetDeviceTokenApiAuthDeviceDeviceTokenPostWithResponse request
	GetDeviceTokenApiAuthDeviceDeviceTokenPostWithResponse(ctx context.Context, params *GetDeviceTokenApiAuthDeviceDeviceTokenPostParams, reqEditors ...RequestEditorFn) (*GetDeviceTokenApiAuthDeviceDeviceTokenPostResponse, error)

	// GetJwtTokenApiAuthDeviceJwtTokenPostWithResponse request
	GetJwtTokenApiAuthDeviceJwtTokenPostWithResponse(ctx context.Context, params *GetJwtTokenApiAuthDeviceJwtTokenPostParams, reqEditors ...RequestEditorFn) (*GetJwtTokenApiAuthDeviceJwtTokenPostResponse, error)

	// AuthorizeApiAuthUserGetWithResponse request
	AuthorizeApiAuthUserGetWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*AuthorizeApiAuthUserGetResponse, error)
}

type GetAwsCredentialsApiAuthAwsCredentialsGetResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *AWSCredentials
	JSON422      *HTTPValidationError
}

// Status returns HTTPResponse.Status
func (r GetAwsCredentialsApiAuthAwsCredentialsGetResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAwsCredentialsApiAuthAwsCredentialsGetResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetDeviceTokenApiAuthDeviceDeviceTokenPostResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *interface{}
	JSON422      *HTTPValidationError
}

// Status returns HTTPResponse.Status
func (r GetDeviceTokenApiAuthDeviceDeviceTokenPostResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetDeviceTokenApiAuthDeviceDeviceTokenPostResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetJwtTokenApiAuthDeviceJwtTokenPostResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *interface{}
	JSON422      *HTTPValidationError
}

// Status returns HTTPResponse.Status
func (r GetJwtTokenApiAuthDeviceJwtTokenPostResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetJwtTokenApiAuthDeviceJwtTokenPostResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AuthorizeApiAuthUserGetResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *User
}

// Status returns HTTPResponse.Status
func (r AuthorizeApiAuthUserGetResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AuthorizeApiAuthUserGetResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetAwsCredentialsApiAuthAwsCredentialsGetWithResponse request returning *GetAwsCredentialsApiAuthAwsCredentialsGetResponse
func (c *ClientWithResponses) GetAwsCredentialsApiAuthAwsCredentialsGetWithResponse(ctx context.Context, params *GetAwsCredentialsApiAuthAwsCredentialsGetParams, reqEditors ...RequestEditorFn) (*GetAwsCredentialsApiAuthAwsCredentialsGetResponse, error) {
	rsp, err := c.GetAwsCredentialsApiAuthAwsCredentialsGet(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAwsCredentialsApiAuthAwsCredentialsGetResponse(rsp)
}

// GetDeviceTokenApiAuthDeviceDeviceTokenPostWithResponse request returning *GetDeviceTokenApiAuthDeviceDeviceTokenPostResponse
func (c *ClientWithResponses) GetDeviceTokenApiAuthDeviceDeviceTokenPostWithResponse(ctx context.Context, params *GetDeviceTokenApiAuthDeviceDeviceTokenPostParams, reqEditors ...RequestEditorFn) (*GetDeviceTokenApiAuthDeviceDeviceTokenPostResponse, error) {
	rsp, err := c.GetDeviceTokenApiAuthDeviceDeviceTokenPost(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetDeviceTokenApiAuthDeviceDeviceTokenPostResponse(rsp)
}

// GetJwtTokenApiAuthDeviceJwtTokenPostWithResponse request returning *GetJwtTokenApiAuthDeviceJwtTokenPostResponse
func (c *ClientWithResponses) GetJwtTokenApiAuthDeviceJwtTokenPostWithResponse(ctx context.Context, params *GetJwtTokenApiAuthDeviceJwtTokenPostParams, reqEditors ...RequestEditorFn) (*GetJwtTokenApiAuthDeviceJwtTokenPostResponse, error) {
	rsp, err := c.GetJwtTokenApiAuthDeviceJwtTokenPost(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetJwtTokenApiAuthDeviceJwtTokenPostResponse(rsp)
}

// AuthorizeApiAuthUserGetWithResponse request returning *AuthorizeApiAuthUserGetResponse
func (c *ClientWithResponses) AuthorizeApiAuthUserGetWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*AuthorizeApiAuthUserGetResponse, error) {
	rsp, err := c.AuthorizeApiAuthUserGet(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAuthorizeApiAuthUserGetResponse(rsp)
}

// ParseGetAwsCredentialsApiAuthAwsCredentialsGetResponse parses an HTTP response from a GetAwsCredentialsApiAuthAwsCredentialsGetWithResponse call
func ParseGetAwsCredentialsApiAuthAwsCredentialsGetResponse(rsp *http.Response) (*GetAwsCredentialsApiAuthAwsCredentialsGetResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAwsCredentialsApiAuthAwsCredentialsGetResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest AWSCredentials
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 422:
		var dest HTTPValidationError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON422 = &dest

	}

	return response, nil
}

// ParseGetDeviceTokenApiAuthDeviceDeviceTokenPostResponse parses an HTTP response from a GetDeviceTokenApiAuthDeviceDeviceTokenPostWithResponse call
func ParseGetDeviceTokenApiAuthDeviceDeviceTokenPostResponse(rsp *http.Response) (*GetDeviceTokenApiAuthDeviceDeviceTokenPostResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetDeviceTokenApiAuthDeviceDeviceTokenPostResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest interface{}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 422:
		var dest HTTPValidationError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON422 = &dest

	}

	return response, nil
}

// ParseGetJwtTokenApiAuthDeviceJwtTokenPostResponse parses an HTTP response from a GetJwtTokenApiAuthDeviceJwtTokenPostWithResponse call
func ParseGetJwtTokenApiAuthDeviceJwtTokenPostResponse(rsp *http.Response) (*GetJwtTokenApiAuthDeviceJwtTokenPostResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetJwtTokenApiAuthDeviceJwtTokenPostResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest interface{}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 422:
		var dest HTTPValidationError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON422 = &dest

	}

	return response, nil
}

// ParseAuthorizeApiAuthUserGetResponse parses an HTTP response from a AuthorizeApiAuthUserGetWithResponse call
func ParseAuthorizeApiAuthUserGetResponse(rsp *http.Response) (*AuthorizeApiAuthUserGetResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AuthorizeApiAuthUserGetResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest User
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get Aws Credentials
	// (GET /api/auth/aws_credentials)
	GetAwsCredentialsApiAuthAwsCredentialsGet(ctx echo.Context, params GetAwsCredentialsApiAuthAwsCredentialsGetParams) error
	// Get Device Token
	// (POST /api/auth/device/device_token)
	GetDeviceTokenApiAuthDeviceDeviceTokenPost(ctx echo.Context, params GetDeviceTokenApiAuthDeviceDeviceTokenPostParams) error
	// Get Jwt Token
	// (POST /api/auth/device/jwt_token)
	GetJwtTokenApiAuthDeviceJwtTokenPost(ctx echo.Context, params GetJwtTokenApiAuthDeviceJwtTokenPostParams) error
	// Authorize
	// (GET /api/auth/user)
	AuthorizeApiAuthUserGet(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetAwsCredentialsApiAuthAwsCredentialsGet converts echo context to params.
func (w *ServerInterfaceWrapper) GetAwsCredentialsApiAuthAwsCredentialsGet(ctx echo.Context) error {
	var err error

	ctx.Set(OpenIdConnectScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAwsCredentialsApiAuthAwsCredentialsGetParams
	// ------------- Required query parameter "data_product_name" -------------

	err = runtime.BindQueryParameter("form", true, true, "data_product_name", ctx.QueryParams(), &params.DataProductName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter data_product_name: %s", err))
	}

	// ------------- Required query parameter "environment" -------------

	err = runtime.BindQueryParameter("form", true, true, "environment", ctx.QueryParams(), &params.Environment)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter environment: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAwsCredentialsApiAuthAwsCredentialsGet(ctx, params)
	return err
}

// GetDeviceTokenApiAuthDeviceDeviceTokenPost converts echo context to params.
func (w *ServerInterfaceWrapper) GetDeviceTokenApiAuthDeviceDeviceTokenPost(ctx echo.Context) error {
	var err error

	ctx.Set(HTTPBasicScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetDeviceTokenApiAuthDeviceDeviceTokenPostParams
	// ------------- Required query parameter "client_id" -------------

	err = runtime.BindQueryParameter("form", true, true, "client_id", ctx.QueryParams(), &params.ClientId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter client_id: %s", err))
	}

	// ------------- Optional query parameter "scope" -------------

	err = runtime.BindQueryParameter("form", true, false, "scope", ctx.QueryParams(), &params.Scope)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scope: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetDeviceTokenApiAuthDeviceDeviceTokenPost(ctx, params)
	return err
}

// GetJwtTokenApiAuthDeviceJwtTokenPost converts echo context to params.
func (w *ServerInterfaceWrapper) GetJwtTokenApiAuthDeviceJwtTokenPost(ctx echo.Context) error {
	var err error

	ctx.Set(HTTPBasicScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetJwtTokenApiAuthDeviceJwtTokenPostParams
	// ------------- Required query parameter "client_id" -------------

	err = runtime.BindQueryParameter("form", true, true, "client_id", ctx.QueryParams(), &params.ClientId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter client_id: %s", err))
	}

	// ------------- Required query parameter "device_code" -------------

	err = runtime.BindQueryParameter("form", true, true, "device_code", ctx.QueryParams(), &params.DeviceCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter device_code: %s", err))
	}

	// ------------- Required query parameter "grant_type" -------------

	err = runtime.BindQueryParameter("form", true, true, "grant_type", ctx.QueryParams(), &params.GrantType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter grant_type: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetJwtTokenApiAuthDeviceJwtTokenPost(ctx, params)
	return err
}

// AuthorizeApiAuthUserGet converts echo context to params.
func (w *ServerInterfaceWrapper) AuthorizeApiAuthUserGet(ctx echo.Context) error {
	var err error

	ctx.Set(OpenIdConnectScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AuthorizeApiAuthUserGet(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/api/auth/aws_credentials", wrapper.GetAwsCredentialsApiAuthAwsCredentialsGet)
	router.POST(baseURL+"/api/auth/device/device_token", wrapper.GetDeviceTokenApiAuthDeviceDeviceTokenPost)
	router.POST(baseURL+"/api/auth/device/jwt_token", wrapper.GetJwtTokenApiAuthDeviceJwtTokenPost)
	router.GET(baseURL+"/api/auth/user", wrapper.AuthorizeApiAuthUserGet)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xY224bNxN+lQX//1LSKmmvdFXFdlwf0hiV0wI1DIHaHcm0dskNOStlI+jdiyG5J2l1",
	"SJECLdCbRBrO4ZuZjzOUNyxSaaYkSDRstGEmeoGU24/j3ycXGmKQKHhiJZlWGWgU4M6jCIy5g+Impq8o",
	"MAE28uIlFCJmPYZFRkKDWsgF2/bY1ZdMaI5CSTKaK51yZCMWc4Q+ihTIxntqqHY4mkCkASsQTQjuiJdA",
	"uq2NEUo+qiXItqmVo5Xv2W17TMPnXGiI2eipVYF9QDtBWqk/11nulLmKqWavECFh/fnx8eE3nojY2l5p",
	"rfR+N2JALhL6JBBSK/q/hjkbsf+FdYtD399w19+2BnTpPFVAuNa8sLmXGl2AOnB/MtABFFKPs+q9kzT6",
	"nrbi102DLwha8mQqWoy78uLgppNxc6ENTiVPoWn0nqTBL9wxbtfGBagQ5rkjs7ftjpPwjjD3/GCUHTKV",
	"VWgm2QLfjGARNkhkS93RgpO0SVTU4gyXxcc5Gz1t9vKrJEIiLECz7XODEvcq2rmpnjY9lppFsyIfwBi+",
	"6Ky6EzQaS5iDR5Keqh7l4UJ5zUZxTlJ122MGolwLLCZ0P1xpiOXvuBFRNRfJZmYllYsXxIygf8xA3sQX",
	"SkryONow1RR80onXNaMwjNRCClR9EWcDyPtrMNh/M+Ap/6okX5tBpNKwkk8/3KV/3N2/v5iFgzUkSX8p",
	"1VqG5F7E/UjJuVjku1OyFZxtKUEh54pwRUoidxD9TWQGxascvHJpDEjzU8yRp0LGEBMS1mOO02xCasGt",
	"V6OkTZ6mXBdsxN7xaAkyDsYPN4FIswRSkGhBBXOlg0uOPMi0ivMIg0xp5EkvKFSuA5AIOtPCQF8Dj4uA",
	"ogdudAeoVHMudHhhPbYCbew2YcPBcDAkXJQ+zwQbsR8GbwZD1mMZxxfb1JCNNtseC3kmQp7jS8jXZhq1",
	"19wCyg66utJ+Y9eA47VpTOpxJsY5vrSF14A2muYpIGhjb5IgcJ9z0EVdTEpz6nMpr3TNZ9Q59Pwubl4I",
	"W4EHX4FDU6U7IMiV0EpSX84KddXS3w3yTC5MpqRxd+XtcFiSiwxokGRZItxMCF+N2/Z1lGOraWchWvbG",
	"YCItMvdsYJPc8mOeJ8GvHgW1/ce3b78biK4V14GkVglKnXqW2ObvDIanZ6pdfXGuAYPx2gQ7TwC+IOow",
	"Yih73rYIG8NKROD/m2L5gMmU6WbtpVW0bxBPWSdpyB/I9izaRokAiW41nebQhdXuXMuHaGoiZad97SyG",
	"Oc8T9FOttYcnXvk7s/NfzrjG3upim+t7UD5KW1SjYPZ0nqj1AeK9rvEM1t2ucZ9ypfAfxDd/jSIVnzmB",
	"XX0unP6ZQRaaS5yie8ecjnFN6oeePf9x+wi3b9f4F4id+98pnXuf2Ku0+AqeyfTSdlv+b1uA9jF/bqO+",
	"cd9U6RzeMrPcCAnGTLkG+5eA8qD1ZqEvBnCaCLk0Ic8yrVYQbkS8Pc8iBll8g7qG9JR/uimH4B6U2+aH",
	"G/p3eti3ORbYHYZ8pnI8rhIpuYJCUSSY5vSb4CxtqRBmSi1Pm/iahZuyeCdRa0Uj55iGEQsp5OnYRwtp",
	"YK8FJDpt1FX4St5Z8/r0oHOQqyYa0tj7XgXebv8MAAD//8P1NmwkEwAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
