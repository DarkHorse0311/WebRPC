// example-api-service v1.0.0 f1b366018b1650b6a0a4f09ff793089ec62830a6
// --
// Code generated by webrpc-gen@v0.10.x-dev with golang generator. DO NOT EDIT.
//
// webrpc-gen -schema=./proto/api.ridl -target=golang -pkg=main -server -client -out=./api.gen.go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// WebRPC description and code-gen version
func WebRPCVersion() string {
	return "v1"
}

// Schema version of your RIDL schema
func WebRPCSchemaVersion() string {
	return "v1.0.0"
}

// Schema hash generated from your RIDL schema
func WebRPCSchemaHash() string {
	return "f1b366018b1650b6a0a4f09ff793089ec62830a6"
}

//
// Types
//

type User struct {
	Username string `json:"username"`
	Age      uint32 `json:"age"`
}

type Location uint32

const (
	Location_TORONTO  Location = 0
	Location_NEW_YORK Location = 1
)

var Location_name = map[uint32]string{
	0: "TORONTO",
	1: "NEW_YORK",
}

var Location_value = map[string]uint32{
	"TORONTO":  0,
	"NEW_YORK": 1,
}

func (x Location) String() string {
	return Location_name[uint32(x)]
}

func (x Location) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	buf.WriteString(Location_name[uint32(x)])
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

func (x *Location) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*x = Location(Location_value[j])
	return nil
}

type ExampleAPI interface {
	Ping(ctx context.Context) error
	Status(ctx context.Context) (bool, error)
	GetUsers(ctx context.Context) ([]*User, *Location, error)
}

var WebRPCServices = map[string][]string{
	"ExampleAPI": {
		"Ping",
		"Status",
		"GetUsers",
	},
}

//
// Server
//

type WebRPCServer interface {
	http.Handler
}

type exampleAPIServer struct {
	ExampleAPI
}

func NewExampleAPIServer(svc ExampleAPI) WebRPCServer {
	return &exampleAPIServer{
		ExampleAPI: svc,
	}
}

func (s *exampleAPIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, HTTPResponseWriterCtxKey, w)
	ctx = context.WithValue(ctx, HTTPRequestCtxKey, r)
	ctx = context.WithValue(ctx, ServiceNameCtxKey, "ExampleAPI")

	if r.Method != "POST" {
		err := Errorf(ErrBadRoute, "unsupported method %q (only POST is allowed)", r.Method)
		RespondWithError(w, err)
		return
	}

	switch r.URL.Path {
	case "/rpc/ExampleAPI/Ping":
		s.servePing(ctx, w, r)
		return
	case "/rpc/ExampleAPI/Status":
		s.serveStatus(ctx, w, r)
		return
	case "/rpc/ExampleAPI/GetUsers":
		s.serveGetUsers(ctx, w, r)
		return
	default:
		err := Errorf(ErrBadRoute, "no handler for path %q", r.URL.Path)
		RespondWithError(w, err)
		return
	}
}

func (s *exampleAPIServer) servePing(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}

	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.servePingJSON(ctx, w, r)
	default:
		err := Errorf(ErrBadRoute, "unexpected Content-Type: %q", r.Header.Get("Content-Type"))
		RespondWithError(w, err)
	}
}

func (s *exampleAPIServer) servePingJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error
	ctx = context.WithValue(ctx, MethodNameCtxKey, "Ping")

	// Call service method
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if rr := recover(); rr != nil {
				RespondWithError(w, ErrorInternal("internal service panic"))
				panic(rr)
			}
		}()
		err = s.ExampleAPI.Ping(ctx)
	}()

	if err != nil {
		RespondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *exampleAPIServer) serveStatus(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}

	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveStatusJSON(ctx, w, r)
	default:
		err := Errorf(ErrBadRoute, "unexpected Content-Type: %q", r.Header.Get("Content-Type"))
		RespondWithError(w, err)
	}
}

func (s *exampleAPIServer) serveStatusJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error
	ctx = context.WithValue(ctx, MethodNameCtxKey, "Status")

	// Call service method
	var ret0 bool
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if rr := recover(); rr != nil {
				RespondWithError(w, ErrorInternal("internal service panic"))
				panic(rr)
			}
		}()
		ret0, err = s.ExampleAPI.Status(ctx)
	}()
	respContent := struct {
		Ret0 bool `json:"status"`
	}{ret0}

	if err != nil {
		RespondWithError(w, err)
		return
	}
	respBody, err := json.Marshal(respContent)
	if err != nil {
		err = WrapError(ErrInternal, err, "failed to marshal json response")
		RespondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (s *exampleAPIServer) serveGetUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}

	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveGetUsersJSON(ctx, w, r)
	default:
		err := Errorf(ErrBadRoute, "unexpected Content-Type: %q", r.Header.Get("Content-Type"))
		RespondWithError(w, err)
	}
}

func (s *exampleAPIServer) serveGetUsersJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error
	ctx = context.WithValue(ctx, MethodNameCtxKey, "GetUsers")

	// Call service method
	var ret0 []*User
	var ret1 *Location
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if rr := recover(); rr != nil {
				RespondWithError(w, ErrorInternal("internal service panic"))
				panic(rr)
			}
		}()
		ret0, ret1, err = s.ExampleAPI.GetUsers(ctx)
	}()
	respContent := struct {
		Ret0 []*User   `json:"users"`
		Ret1 *Location `json:"location"`
	}{ret0, ret1}

	if err != nil {
		RespondWithError(w, err)
		return
	}
	respBody, err := json.Marshal(respContent)
	if err != nil {
		err = WrapError(ErrInternal, err, "failed to marshal json response")
		RespondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func RespondWithError(w http.ResponseWriter, err error) {
	rpcErr, ok := err.(Error)
	if !ok {
		rpcErr = WrapError(ErrInternal, err, "webrpc error")
	}

	statusCode := HTTPStatusFromErrorCode(rpcErr.Code())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	respBody, _ := json.Marshal(rpcErr.Payload())
	w.Write(respBody)
}

//
// Client
//

const ExampleAPIPathPrefix = "/rpc/ExampleAPI/"

type exampleAPIClient struct {
	client HTTPClient
	urls   [3]string
}

func NewExampleAPIClient(addr string, client HTTPClient) ExampleAPI {
	prefix := urlBase(addr) + ExampleAPIPathPrefix
	urls := [3]string{
		prefix + "Ping",
		prefix + "Status",
		prefix + "GetUsers",
	}
	return &exampleAPIClient{
		client: client,
		urls:   urls,
	}
}

func (c *exampleAPIClient) Ping(ctx context.Context) error {

	err := doJSONRequest(ctx, c.client, c.urls[0], nil, nil)
	return err
}

func (c *exampleAPIClient) Status(ctx context.Context) (bool, error) {
	out := struct {
		Ret0 bool `json:"status"`
	}{}

	err := doJSONRequest(ctx, c.client, c.urls[1], nil, &out)
	return out.Ret0, err
}

func (c *exampleAPIClient) GetUsers(ctx context.Context) ([]*User, *Location, error) {
	out := struct {
		Ret0 []*User   `json:"users"`
		Ret1 *Location `json:"location"`
	}{}

	err := doJSONRequest(ctx, c.client, c.urls[2], nil, &out)
	return out.Ret0, out.Ret1, err
}

// HTTPClient is the interface used by generated clients to send HTTP requests.
// It is fulfilled by *(net/http).Client, which is sufficient for most users.
// Users can provide their own implementation for special retry policies.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// urlBase helps ensure that addr specifies a scheme. If it is unparsable
// as a URL, it returns addr unchanged.
func urlBase(addr string) string {
	// If the addr specifies a scheme, use it. If not, default to
	// http. If url.Parse fails on it, return it unchanged.
	url, err := url.Parse(addr)
	if err != nil {
		return addr
	}
	if url.Scheme == "" {
		url.Scheme = "http"
	}
	return url.String()
}

// newRequest makes an http.Request from a client, adding common headers.
func newRequest(ctx context.Context, url string, reqBody io.Reader, contentType string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", contentType)
	req.Header.Set("Content-Type", contentType)
	if headers, ok := HTTPRequestHeaders(ctx); ok {
		for k := range headers {
			for _, v := range headers[k] {
				req.Header.Add(k, v)
			}
		}
	}
	return req, nil
}

// doJSONRequest is common code to make a request to the remote service.
func doJSONRequest(ctx context.Context, client HTTPClient, url string, in, out interface{}) error {
	reqBody, err := json.Marshal(in)
	if err != nil {
		return clientError("failed to marshal json request", err)
	}
	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}

	req, err := newRequest(ctx, url, bytes.NewBuffer(reqBody), "application/json")
	if err != nil {
		return clientError("could not build request", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return clientError("request failed", err)
	}

	defer func() {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = clientError("failed to close response body", cerr)
		}
	}()

	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}

	if resp.StatusCode != 200 {
		return errorFromResponse(resp)
	}

	if out != nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return clientError("failed to read response body", err)
		}

		err = json.Unmarshal(respBody, &out)
		if err != nil {
			return clientError("failed to unmarshal json response body", err)
		}
		if err = ctx.Err(); err != nil {
			return clientError("aborted because context was done", err)
		}
	}

	return nil
}

// errorFromResponse builds a webrpc Error from a non-200 HTTP response.
func errorFromResponse(resp *http.Response) Error {
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return clientError("failed to read server error response body", err)
	}

	var respErr ErrorPayload
	if err := json.Unmarshal(respBody, &respErr); err != nil {
		return clientError("failed unmarshal error response", err)
	}

	errCode := ErrorCode(respErr.Code)

	if HTTPStatusFromErrorCode(errCode) == 0 {
		return ErrorInternal("invalid code returned from server error response: %s", respErr.Code)
	}

	return &rpcErr{
		code:  errCode,
		msg:   respErr.Msg,
		cause: errors.New(respErr.Cause),
	}
}

func clientError(desc string, err error) Error {
	return WrapError(ErrInternal, err, desc)
}

func WithHTTPRequestHeaders(ctx context.Context, h http.Header) (context.Context, error) {
	if _, ok := h["Accept"]; ok {
		return nil, errors.New("provided header cannot set Accept")
	}
	if _, ok := h["Content-Type"]; ok {
		return nil, errors.New("provided header cannot set Content-Type")
	}

	copied := make(http.Header, len(h))
	for k, vv := range h {
		if vv == nil {
			copied[k] = nil
			continue
		}
		copied[k] = make([]string, len(vv))
		copy(copied[k], vv)
	}

	return context.WithValue(ctx, HTTPClientRequestHeadersCtxKey, copied), nil
}

func HTTPRequestHeaders(ctx context.Context) (http.Header, bool) {
	h, ok := ctx.Value(HTTPClientRequestHeadersCtxKey).(http.Header)
	return h, ok
}

//
// Helpers
//

type ErrorPayload struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Cause  string `json:"cause,omitempty"`
	Msg    string `json:"msg"`
	Error  string `json:"error"`
}

type Error interface {
	// Code is of the valid error codes
	Code() ErrorCode

	// Msg returns a human-readable, unstructured messages describing the error
	Msg() string

	// Cause is reason for the error
	Cause() error

	// Error returns a string of the form "webrpc error <Code>: <Msg>"
	Error() string

	// Error response payload
	Payload() ErrorPayload
}

func Errorf(code ErrorCode, msgf string, args ...interface{}) Error {
	msg := fmt.Sprintf(msgf, args...)
	if IsValidErrorCode(code) {
		return &rpcErr{code: code, msg: msg}
	}
	return &rpcErr{code: ErrInternal, msg: "invalid error type " + string(code)}
}

func WrapError(code ErrorCode, cause error, format string, args ...interface{}) Error {
	msg := fmt.Sprintf(format, args...)
	if IsValidErrorCode(code) {
		return &rpcErr{code: code, msg: msg, cause: cause}
	}
	return &rpcErr{code: ErrInternal, msg: "invalid error type " + string(code), cause: cause}
}

func Failf(format string, args ...interface{}) Error {
	return Errorf(ErrFail, format, args...)
}

func WrapFailf(cause error, format string, args ...interface{}) Error {
	return WrapError(ErrFail, cause, format, args...)
}

func ErrorNotFound(format string, args ...interface{}) Error {
	return Errorf(ErrNotFound, format, args...)
}

func ErrorInvalidArgument(argument string, validationMsg string) Error {
	return Errorf(ErrInvalidArgument, argument+" "+validationMsg)
}

func ErrorRequiredArgument(argument string) Error {
	return ErrorInvalidArgument(argument, "is required")
}

func ErrorInternal(format string, args ...interface{}) Error {
	return Errorf(ErrInternal, format, args...)
}

type ErrorCode string

const (
	// Unknown error. For example when handling errors raised by APIs that do not
	// return enough error information.
	ErrUnknown ErrorCode = "unknown"

	// Fail error. General failure error type.
	ErrFail ErrorCode = "fail"

	// Canceled indicates the operation was cancelled (typically by the caller).
	ErrCanceled ErrorCode = "canceled"

	// InvalidArgument indicates client specified an invalid argument. It
	// indicates arguments that are problematic regardless of the state of the
	// system (i.e. a malformed file name, required argument, number out of range,
	// etc.).
	ErrInvalidArgument ErrorCode = "invalid argument"

	// DeadlineExceeded means operation expired before completion. For operations
	// that change the state of the system, this error may be returned even if the
	// operation has completed successfully (timeout).
	ErrDeadlineExceeded ErrorCode = "deadline exceeded"

	// NotFound means some requested entity was not found.
	ErrNotFound ErrorCode = "not found"

	// BadRoute means that the requested URL path wasn't routable to a webrpc
	// service and method. This is returned by the generated server, and usually
	// shouldn't be returned by applications. Instead, applications should use
	// NotFound or Unimplemented.
	ErrBadRoute ErrorCode = "bad route"

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	ErrAlreadyExists ErrorCode = "already exists"

	// PermissionDenied indicates the caller does not have permission to execute
	// the specified operation. It must not be used if the caller cannot be
	// identified (Unauthenticated).
	ErrPermissionDenied ErrorCode = "permission denied"

	// Unauthenticated indicates the request does not have valid authentication
	// credentials for the operation.
	ErrUnauthenticated ErrorCode = "unauthenticated"

	// ResourceExhausted indicates some resource has been exhausted, perhaps a
	// per-user quota, or perhaps the entire file system is out of space.
	ErrResourceExhausted ErrorCode = "resource exhausted"

	// FailedPrecondition indicates operation was rejected because the system is
	// not in a state required for the operation's execution. For example, doing
	// an rmdir operation on a directory that is non-empty, or on a non-directory
	// object, or when having conflicting read-modify-write on the same resource.
	ErrFailedPrecondition ErrorCode = "failed precondition"

	// Aborted indicates the operation was aborted, typically due to a concurrency
	// issue like sequencer check failures, transaction aborts, etc.
	ErrAborted ErrorCode = "aborted"

	// OutOfRange means operation was attempted past the valid range. For example,
	// seeking or reading past end of a paginated collection.
	//
	// Unlike InvalidArgument, this error indicates a problem that may be fixed if
	// the system state changes (i.e. adding more items to the collection).
	//
	// There is a fair bit of overlap between FailedPrecondition and OutOfRange.
	// We recommend using OutOfRange (the more specific error) when it applies so
	// that callers who are iterating through a space can easily look for an
	// OutOfRange error to detect when they are done.
	ErrOutOfRange ErrorCode = "out of range"

	// Unimplemented indicates operation is not implemented or not
	// supported/enabled in this service.
	ErrUnimplemented ErrorCode = "unimplemented"

	// Internal errors. When some invariants expected by the underlying system
	// have been broken. In other words, something bad happened in the library or
	// backend service. Do not confuse with HTTP Internal Server Error; an
	// Internal error could also happen on the client code, i.e. when parsing a
	// server response.
	ErrInternal ErrorCode = "internal"

	// Unavailable indicates the service is currently unavailable. This is a most
	// likely a transient condition and may be corrected by retrying with a
	// backoff.
	ErrUnavailable ErrorCode = "unavailable"

	// DataLoss indicates unrecoverable data loss or corruption.
	ErrDataLoss ErrorCode = "data loss"

	// ErrNone is the zero-value, is considered an empty error and should not be
	// used.
	ErrNone ErrorCode = ""
)

func HTTPStatusFromErrorCode(code ErrorCode) int {
	switch code {
	case ErrCanceled:
		return 408 // RequestTimeout
	case ErrUnknown:
		return 400 // Bad Request
	case ErrFail:
		return 422 // Unprocessable Entity
	case ErrInvalidArgument:
		return 400 // BadRequest
	case ErrDeadlineExceeded:
		return 408 // RequestTimeout
	case ErrNotFound:
		return 404 // Not Found
	case ErrBadRoute:
		return 404 // Not Found
	case ErrAlreadyExists:
		return 409 // Conflict
	case ErrPermissionDenied:
		return 403 // Forbidden
	case ErrUnauthenticated:
		return 401 // Unauthorized
	case ErrResourceExhausted:
		return 403 // Forbidden
	case ErrFailedPrecondition:
		return 412 // Precondition Failed
	case ErrAborted:
		return 409 // Conflict
	case ErrOutOfRange:
		return 400 // Bad Request
	case ErrUnimplemented:
		return 501 // Not Implemented
	case ErrInternal:
		return 500 // Internal Server Error
	case ErrUnavailable:
		return 503 // Service Unavailable
	case ErrDataLoss:
		return 500 // Internal Server Error
	case ErrNone:
		return 200 // OK
	default:
		return 0 // Invalid!
	}
}

func IsErrorCode(err error, code ErrorCode) bool {
	if rpcErr, ok := err.(Error); ok {
		if rpcErr.Code() == code {
			return true
		}
	}
	return false
}

func IsValidErrorCode(code ErrorCode) bool {
	return HTTPStatusFromErrorCode(code) != 0
}

type rpcErr struct {
	code  ErrorCode
	msg   string
	cause error
}

func (e *rpcErr) Code() ErrorCode {
	return e.code
}

func (e *rpcErr) Msg() string {
	return e.msg
}

func (e *rpcErr) Cause() error {
	return e.cause
}

func (e *rpcErr) Error() string {
	if e.cause != nil && e.cause.Error() != "" {
		if e.msg != "" {
			return fmt.Sprintf("webrpc %s error: %s -- %s", e.code, e.cause.Error(), e.msg)
		} else {
			return fmt.Sprintf("webrpc %s error: %s", e.code, e.cause.Error())
		}
	} else {
		return fmt.Sprintf("webrpc %s error: %s", e.code, e.msg)
	}
}

func (e *rpcErr) Payload() ErrorPayload {
	statusCode := HTTPStatusFromErrorCode(e.Code())
	errPayload := ErrorPayload{
		Status: statusCode,
		Code:   string(e.Code()),
		Msg:    e.Msg(),
		Error:  e.Error(),
	}
	if e.Cause() != nil {
		errPayload.Cause = e.Cause().Error()
	}
	return errPayload
}

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "webrpc context value " + k.name
}

var (
	// For Client
	HTTPClientRequestHeadersCtxKey = &contextKey{"HTTPClientRequestHeaders"}

	// For Server
	HTTPResponseWriterCtxKey = &contextKey{"HTTPResponseWriter"}

	HTTPRequestCtxKey = &contextKey{"HTTPRequest"}

	ServiceNameCtxKey = &contextKey{"ServiceName"}

	MethodNameCtxKey = &contextKey{"MethodName"}
)
