// example  v0.0.1 d1747ce08492a2ae743797161ea176c0ff23f049
// --
// Code generated by webrpc-gen@v0.14.0-dev with golang generator. DO NOT EDIT.
//
// webrpc-gen -schema=example.webrpc.json -target=golang -pkg=main -server -out=./server/server.gen.go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// WebRPC description and code-gen version
func WebRPCVersion() string {
	return "v1"
}

// Schema version of your RIDL schema
func WebRPCSchemaVersion() string {
	return " v0.0.1"
}

// Schema hash generated from your RIDL schema
func WebRPCSchemaHash() string {
	return "d1747ce08492a2ae743797161ea176c0ff23f049"
}

//
// Types
//

type Kind uint32

const (
	Kind_USER  Kind = 1
	Kind_ADMIN Kind = 2
)

var Kind_name = map[uint32]string{
	1: "USER",
	2: "ADMIN",
}

var Kind_value = map[string]uint32{
	"USER":  1,
	"ADMIN": 2,
}

func (x Kind) String() string {
	return Kind_name[uint32(x)]
}

func (x Kind) MarshalText() ([]byte, error) {
	return []byte(Kind_name[uint32(x)]), nil
}

func (x *Kind) UnmarshalText(b []byte) error {
	*x = Kind(Kind_value[string(b)])
	return nil
}

type Empty struct {
}

type GetUserRequest struct {
	UserID uint64 `json:"userID"`
}

type User struct {
	ID        uint64     `json:"id" db:"id"`
	Username  string     `json:"USERNAME" db:"username"`
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
}

type RandomStuff struct {
	Meta              map[string]interface{}       `json:"meta"`
	MetaNestedExample map[string]map[string]uint32 `json:"metaNestedExample"`
	NamesList         []string                     `json:"namesList"`
	NumsList          []int64                      `json:"numsList"`
	DoubleArray       [][]string                   `json:"doubleArray"`
	ListOfMaps        []map[string]uint32          `json:"listOfMaps"`
	ListOfUsers       []*User                      `json:"listOfUsers"`
	MapOfUsers        map[string]*User             `json:"mapOfUsers"`
	User              *User                        `json:"user"`
}

type ExampleService interface {
	Ping(ctx context.Context) (bool, error)
	GetUser(ctx context.Context, req *GetUserRequest) (*User, error)
}

var WebRPCServices = map[string][]string{
	"ExampleService": {
		"Ping",
		"GetUser",
	},
}

//
// Server
//

type WebRPCServer interface {
	http.Handler
}

type exampleServiceServer struct {
	ExampleService
}

func NewExampleServiceServer(svc ExampleService) WebRPCServer {
	return &exampleServiceServer{
		ExampleService: svc,
	}
}

func (s *exampleServiceServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// In case of a panic, serve a HTTP 500 error and then panic.
		if rr := recover(); rr != nil {
			RespondWithError(w, ErrWebrpcServerPanic.WithCause(fmt.Errorf("%v", rr)))
			panic(rr)
		}
	}()

	ctx := r.Context()
	ctx = context.WithValue(ctx, HTTPResponseWriterCtxKey, w)
	ctx = context.WithValue(ctx, HTTPRequestCtxKey, r)
	ctx = context.WithValue(ctx, ServiceNameCtxKey, "ExampleService")

	var handler func(ctx context.Context, w http.ResponseWriter, r *http.Request)
	switch r.URL.Path {
	case "/rpc/ExampleService/Ping":
		handler = s.servePingJSON
	case "/rpc/ExampleService/GetUser":
		handler = s.serveGetUserJSON
	default:
		err := ErrWebrpcBadRoute.WithCause(fmt.Errorf("no handler for path %q", r.URL.Path))
		RespondWithError(w, err)
		return
	}

	if r.Method != "POST" {
		w.Header().Add("Allow", "POST") // RFC 9110.
		err := ErrWebrpcBadMethod.WithCause(fmt.Errorf("unsupported method %q (only POST is allowed)", r.Method))
		RespondWithError(w, err)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if i := strings.Index(contentType, ";"); i >= 0 {
		contentType = contentType[:i]
	}
	contentType = strings.TrimSpace(strings.ToLower(contentType))

	switch contentType {
	case "application/json":
		handler(ctx, w, r)
	default:
		err := ErrWebrpcBadRequest.WithCause(fmt.Errorf("unexpected Content-Type: %q", r.Header.Get("Content-Type")))
		RespondWithError(w, err)
	}
}

func (s *exampleServiceServer) servePingJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = context.WithValue(ctx, MethodNameCtxKey, "Ping")

	// Call service method implementation.
	ret0, err := s.ExampleService.Ping(ctx)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	respPayload := struct {
		Ret0 bool `json:"status"`
	}{ret0}
	respBody, err := json.Marshal(respPayload)
	if err != nil {
		RespondWithError(w, ErrWebrpcBadResponse.WithCause(fmt.Errorf("failed to marshal json response: %w", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (s *exampleServiceServer) serveGetUserJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = context.WithValue(ctx, MethodNameCtxKey, "GetUser")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, ErrWebrpcBadRequest.WithCause(fmt.Errorf("failed to read request data: %w", err)))
		return
	}
	defer r.Body.Close()

	reqPayload := struct {
		Arg0 *GetUserRequest `json:"req"`
	}{}
	if err := json.Unmarshal(reqBody, &reqPayload); err != nil {
		RespondWithError(w, ErrWebrpcBadRequest.WithCause(fmt.Errorf("failed to unmarshal request data: %w", err)))
		return
	}

	// Call service method implementation.
	ret0, err := s.ExampleService.GetUser(ctx, reqPayload.Arg0)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	respPayload := struct {
		Ret0 *User `json:"user"`
	}{ret0}
	respBody, err := json.Marshal(respPayload)
	if err != nil {
		RespondWithError(w, ErrWebrpcBadResponse.WithCause(fmt.Errorf("failed to marshal json response: %w", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func RespondWithError(w http.ResponseWriter, err error) {
	rpcErr, ok := err.(WebRPCError)
	if !ok {
		rpcErr = ErrWebrpcEndpoint.WithCause(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rpcErr.HTTPStatus)

	respBody, _ := json.Marshal(rpcErr)
	w.Write(respBody)
}

//
// Helpers
//

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "webrpc context value " + k.name
}

var (
	HTTPResponseWriterCtxKey = &contextKey{"HTTPResponseWriter"}

	HTTPRequestCtxKey = &contextKey{"HTTPRequest"}

	ServiceNameCtxKey = &contextKey{"ServiceName"}

	MethodNameCtxKey = &contextKey{"MethodName"}
)

func ServiceNameFromContext(ctx context.Context) string {
	service, _ := ctx.Value(ServiceNameCtxKey).(string)
	return service
}

func MethodNameFromContext(ctx context.Context) string {
	method, _ := ctx.Value(MethodNameCtxKey).(string)
	return method
}

func RequestFromContext(ctx context.Context) *http.Request {
	r, _ := ctx.Value(HTTPRequestCtxKey).(*http.Request)
	return r
}
func ResponseWriterFromContext(ctx context.Context) http.ResponseWriter {
	w, _ := ctx.Value(HTTPResponseWriterCtxKey).(http.ResponseWriter)
	return w
}

//
// Errors
//

type WebRPCError struct {
	Name       string `json:"error"`
	Code       int    `json:"code"`
	Message    string `json:"msg"`
	Cause      string `json:"cause,omitempty"`
	HTTPStatus int    `json:"status"`
	cause      error
}

var _ error = WebRPCError{}

func (e WebRPCError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s %d: %s: %v", e.Name, e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s %d: %s", e.Name, e.Code, e.Message)
}

func (e WebRPCError) Is(target error) bool {
	if rpcErr, ok := target.(WebRPCError); ok {
		return rpcErr.Code == e.Code
	}
	return errors.Is(e.cause, target)
}

func (e WebRPCError) Unwrap() error {
	return e.cause
}

func (e WebRPCError) WithCause(cause error) WebRPCError {
	err := e
	err.cause = cause
	err.Cause = cause.Error()
	return err
}

// Deprecated: Use .WithCause() method on WebRPCError.
func ErrorWithCause(rpcErr WebRPCError, cause error) WebRPCError {
	return rpcErr.WithCause(cause)
}

// Webrpc errors
var (
	ErrWebrpcEndpoint      = WebRPCError{Code: 0, Name: "WebrpcEndpoint", Message: "endpoint error", HTTPStatus: 400}
	ErrWebrpcRequestFailed = WebRPCError{Code: -1, Name: "WebrpcRequestFailed", Message: "request failed", HTTPStatus: 400}
	ErrWebrpcBadRoute      = WebRPCError{Code: -2, Name: "WebrpcBadRoute", Message: "bad route", HTTPStatus: 404}
	ErrWebrpcBadMethod     = WebRPCError{Code: -3, Name: "WebrpcBadMethod", Message: "bad method", HTTPStatus: 405}
	ErrWebrpcBadRequest    = WebRPCError{Code: -4, Name: "WebrpcBadRequest", Message: "bad request", HTTPStatus: 400}
	ErrWebrpcBadResponse   = WebRPCError{Code: -5, Name: "WebrpcBadResponse", Message: "bad response", HTTPStatus: 500}
	ErrWebrpcServerPanic   = WebRPCError{Code: -6, Name: "WebrpcServerPanic", Message: "server panic", HTTPStatus: 500}
	ErrWebrpcInternalError = WebRPCError{Code: -7, Name: "WebrpcInternalError", Message: "internal error", HTTPStatus: 500}
)

// Schema errors
var (
	ErrUserNotFound = WebRPCError{Code: 1000, Name: "UserNotFound", Message: "user not found", HTTPStatus: 404}
)
