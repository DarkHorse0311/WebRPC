// webrpc-sse-chat v1.0.0 45d7ec19bc1e608515372b36a2c528f9451ef36e
// --
// Code generated by webrpc-gen with golang generator. DO NOT EDIT.
//
// webrpc-gen -schema=proto/chat.ridl -target=golang -pkg=proto -server -client -out=proto/chat.gen.go
package proto

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
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
	return "45d7ec19bc1e608515372b36a2c528f9451ef36e"
}

//
// Common types
//

type Message struct {
	Id        uint64    `json:"id"`
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

var WebRPCServices = map[string][]string{
	"Chat": {
		"SendMessage",
		"SubscribeMessages",
	},
}

//
// Server types
//

type Chat interface {
	SendMessage(ctx context.Context, username string, text string) error
	SubscribeMessages(ctx context.Context, username string, stream SubscribeMessagesStreamWriter) error
}

type SubscribeMessagesStreamWriter interface {
	Write(message *Message) error
}

type subscribeMessagesStreamWriter struct {
	streamWriter
}

func (w *subscribeMessagesStreamWriter) Write(message *Message) error {
	out := struct {
		Ret0 *Message `json:"message"`
	}{
		Ret0: message,
	}

	return w.streamWriter.write(out)
}

type streamWriter struct {
	mu sync.Mutex // Guards concurrent writes to w.
	w  http.ResponseWriter
	f  http.Flusher
	e  *json.Encoder

	sendError func(w http.ResponseWriter, r *http.Request, rpcErr WebRPCError)
}

const StreamKeepAliveInterval = 10 * time.Second

func (w *streamWriter) keepAlive(ctx context.Context) {
	for {
		select {
		case <-time.After(StreamKeepAliveInterval):
			err := w.ping()
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (w *streamWriter) ping() error {
	defer w.f.Flush()

	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.w.Write([]byte("\n"))
	return err
}

func (w *streamWriter) write(respPayload interface{}) error {
	defer w.f.Flush()

	w.mu.Lock()
	defer w.mu.Unlock()

	return w.e.Encode(respPayload)
}

//
// Client types
//

type ChatClient interface {
	SendMessage(ctx context.Context, username string, text string) error
	SubscribeMessages(ctx context.Context, username string) (SubscribeMessagesStreamReader, error)
}

type SubscribeMessagesStreamReader interface {
	Read() (message *Message, err error)
}

//
// Server
//

type WebRPCServer interface {
	http.Handler
}

type chatServer struct {
	Chat
	OnError func(r *http.Request, rpcErr *WebRPCError)
}

func NewChatServer(svc Chat) *chatServer {
	return &chatServer{
		Chat: svc,
	}
}

func (s *chatServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// In case of a panic, serve a HTTP 500 error and then panic.
		if rr := recover(); rr != nil {
			s.sendErrorJSON(w, r, ErrWebrpcServerPanic.WithCause(fmt.Errorf("%v", rr)))
			panic(rr)
		}
	}()

	ctx := r.Context()
	ctx = context.WithValue(ctx, HTTPResponseWriterCtxKey, w)
	ctx = context.WithValue(ctx, HTTPRequestCtxKey, r)
	ctx = context.WithValue(ctx, ServiceNameCtxKey, "Chat")

	var handler func(ctx context.Context, w http.ResponseWriter, r *http.Request)
	switch r.URL.Path {
	case "/rpc/Chat/SendMessage":
		handler = s.serveSendMessageJSON
	case "/rpc/Chat/SubscribeMessages":
		handler = s.serveSubscribeMessagesJSONStream
	default:
		err := ErrWebrpcBadRoute.WithCause(fmt.Errorf("no handler for path %q", r.URL.Path))
		s.sendErrorJSON(w, r, err)
		return
	}

	if r.Method != "POST" {
		w.Header().Add("Allow", "POST") // RFC 9110.
		err := ErrWebrpcBadMethod.WithCause(fmt.Errorf("unsupported method %q (only POST is allowed)", r.Method))
		s.sendErrorJSON(w, r, err)
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
		s.sendErrorJSON(w, r, err)
	}
}

func (s *chatServer) serveSendMessageJSON(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = context.WithValue(ctx, MethodNameCtxKey, "SendMessage")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCause(fmt.Errorf("failed to read request data: %w", err)))
		return
	}
	defer r.Body.Close()

	reqPayload := struct {
		Arg0 string `json:"username"`
		Arg1 string `json:"text"`
	}{}
	if err := json.Unmarshal(reqBody, &reqPayload); err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCause(fmt.Errorf("failed to unmarshal request data: %w", err)))
		return
	}

	// Call service method implementation.
	err = s.Chat.SendMessage(ctx, reqPayload.Arg0, reqPayload.Arg1)
	if err != nil {
		rpcErr, ok := err.(WebRPCError)
		if !ok {
			rpcErr = ErrWebrpcEndpoint.WithCause(err)
		}
		s.sendErrorJSON(w, r, rpcErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *chatServer) serveSubscribeMessagesJSONStream(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx = context.WithValue(ctx, MethodNameCtxKey, "SubscribeMessages")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCause(fmt.Errorf("failed to read request data: %w", err)))
		return
	}
	defer r.Body.Close()

	reqPayload := struct {
		Arg0 string `json:"username"`
	}{}
	if err := json.Unmarshal(reqBody, &reqPayload); err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcBadRequest.WithCause(fmt.Errorf("failed to unmarshal request data: %w", err)))
		return
	}

	f, ok := w.(http.Flusher)
	if !ok {
		s.sendErrorJSON(w, r, ErrWebrpcInternalError.WithCause(fmt.Errorf("server http.ResponseWriter doesn't support .Flush() method")))
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	streamWriter := &subscribeMessagesStreamWriter{streamWriter{w: w, f: f, e: json.NewEncoder(w), sendError: s.sendErrorJSON}}
	if err := streamWriter.ping(); err != nil {
		s.sendErrorJSON(w, r, ErrWebrpcStreamLost.WithCause(fmt.Errorf("failed to establish SSE stream: %w", err)))
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	go streamWriter.keepAlive(ctx)

	// Call service method implementation.
	if err := s.Chat.SubscribeMessages(ctx, reqPayload.Arg0, streamWriter); err != nil {
		rpcErr, ok := err.(WebRPCError)
		if !ok {
			rpcErr = ErrWebrpcEndpoint.WithCause(err)
		}
		streamWriter.sendError(w, r, rpcErr)
		return
	}
}

func (s *chatServer) sendErrorJSON(w http.ResponseWriter, r *http.Request, rpcErr WebRPCError) {
	if s.OnError != nil {
		s.OnError(r, &rpcErr)
	}

	if w.Header().Get("Content-Type") == "application/x-ndjson" {
		out := struct {
			WebRPCError WebRPCError `json:"webrpcError"`
		}{WebRPCError: rpcErr}
		json.NewEncoder(w).Encode(out)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rpcErr.HTTPStatus)

	respBody, _ := json.Marshal(rpcErr)
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
// Client
//

const ChatPathPrefix = "/rpc/Chat/"

type chatClient struct {
	client HTTPClient
	urls   [2]string
}

func NewChatClient(addr string, client HTTPClient) ChatClient {
	prefix := urlBase(addr) + ChatPathPrefix
	urls := [2]string{
		prefix + "SendMessage",
		prefix + "SubscribeMessages",
	}
	return &chatClient{
		client: client,
		urls:   urls,
	}
}

func (c *chatClient) SendMessage(ctx context.Context, username string, text string) error {
	in := struct {
		Arg0 string `json:"username"`
		Arg1 string `json:"text"`
	}{username, text}

	resp, err := doHTTPRequest(ctx, c.client, c.urls[0], in, nil)
	if resp != nil {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = ErrWebrpcRequestFailed.WithCause(fmt.Errorf("failed to close response body: %w", cerr))
		}
	}

	return err
}

func (c *chatClient) SubscribeMessages(ctx context.Context, username string) (SubscribeMessagesStreamReader, error) {
	in := struct {
		Arg0 string `json:"username"`
	}{username}

	resp, err := doHTTPRequest(ctx, c.client, c.urls[1], in, nil)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return nil, err
	}

	buf := bufio.NewReader(resp.Body)
	return &subscribeMessagesStreamReader{streamReader{ctx: ctx, c: resp.Body, r: buf}}, nil
}

type subscribeMessagesStreamReader struct {
	streamReader
}

func (r *subscribeMessagesStreamReader) Read() (*Message, error) {
	out := struct {
		Ret0        *Message     `json:"message"`
		WebRPCError *WebRPCError `json:"webrpcError"`
	}{}

	err := r.streamReader.read(&out)
	if err != nil {
		return out.Ret0, err
	}

	if out.WebRPCError != nil {
		return out.Ret0, out.WebRPCError
	}

	return out.Ret0, nil
}

type streamReader struct {
	ctx context.Context
	c   io.Closer
	r   *bufio.Reader
}

func (r *streamReader) read(v interface{}) error {
	for {
		select {
		case <-r.ctx.Done():
			r.c.Close()
			return ErrWebrpcClientDisconnected.WithCause(r.ctx.Err())
		default:
		}

		line, err := r.r.ReadBytes('\n')
		if err != nil {
			return r.handleReadError(err)
		}

		// Eat newlines (keep-alive pings).
		if len(line) == 1 && line[0] == '\n' {
			continue
		}

		if err := json.Unmarshal(line, &v); err != nil {
			return r.handleReadError(err)
		}
		return nil
	}
}

func (r *streamReader) handleReadError(err error) error {
	defer r.c.Close()
	if errors.Is(err, io.EOF) {
		return ErrWebrpcStreamFinished.WithCause(err)
	}
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return ErrWebrpcStreamLost.WithCause(err)
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return ErrWebrpcClientDisconnected.WithCause(err)
	}
	return ErrWebrpcBadResponse.WithCause(fmt.Errorf("reading stream: %w", err))
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
	req, err := http.NewRequestWithContext(ctx, "POST", url, reqBody)
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

// doHTTPRequest is common code to make a request to the remote service.
func doHTTPRequest(ctx context.Context, client HTTPClient, url string, in, out interface{}) (*http.Response, error) {
	reqBody, err := json.Marshal(in)
	if err != nil {
		return nil, ErrWebrpcRequestFailed.WithCause(fmt.Errorf("failed to marshal JSON body: %w", err))
	}
	if err = ctx.Err(); err != nil {
		return nil, ErrWebrpcRequestFailed.WithCause(fmt.Errorf("aborted because context was done: %w", err))
	}

	req, err := newRequest(ctx, url, bytes.NewBuffer(reqBody), "application/json")
	if err != nil {
		return nil, ErrWebrpcRequestFailed.WithCause(fmt.Errorf("could not build request: %w", err))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, ErrWebrpcRequestFailed.WithCause(err)
	}

	if resp.StatusCode != 200 {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, ErrWebrpcBadResponse.WithCause(fmt.Errorf("failed to read server error response body: %w", err))
		}

		var rpcErr WebRPCError
		if err := json.Unmarshal(respBody, &rpcErr); err != nil {
			return nil, ErrWebrpcBadResponse.WithCause(fmt.Errorf("failed to unmarshal server error: %w", err))
		}
		if rpcErr.Cause != "" {
			rpcErr.cause = errors.New(rpcErr.Cause)
		}
		return nil, rpcErr
	}

	if out != nil {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, ErrWebrpcBadResponse.WithCause(fmt.Errorf("failed to read response body: %w", err))
		}

		err = json.Unmarshal(respBody, &out)
		if err != nil {
			return nil, ErrWebrpcBadResponse.WithCause(fmt.Errorf("failed to unmarshal JSON response body: %w", err))
		}
	}

	return resp, nil
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

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "webrpc context value " + k.name
}

var (
	HTTPClientRequestHeadersCtxKey = &contextKey{"HTTPClientRequestHeaders"}
	HTTPResponseWriterCtxKey       = &contextKey{"HTTPResponseWriter"}

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
	if target == nil {
		return false
	}
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

func (e WebRPCError) WithCausef(format string, args ...interface{}) WebRPCError {
	cause := fmt.Errorf(format, args...)
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
	ErrWebrpcEndpoint           = WebRPCError{Code: 0, Name: "WebrpcEndpoint", Message: "endpoint error", HTTPStatus: 400}
	ErrWebrpcRequestFailed      = WebRPCError{Code: -1, Name: "WebrpcRequestFailed", Message: "request failed", HTTPStatus: 400}
	ErrWebrpcBadRoute           = WebRPCError{Code: -2, Name: "WebrpcBadRoute", Message: "bad route", HTTPStatus: 404}
	ErrWebrpcBadMethod          = WebRPCError{Code: -3, Name: "WebrpcBadMethod", Message: "bad method", HTTPStatus: 405}
	ErrWebrpcBadRequest         = WebRPCError{Code: -4, Name: "WebrpcBadRequest", Message: "bad request", HTTPStatus: 400}
	ErrWebrpcBadResponse        = WebRPCError{Code: -5, Name: "WebrpcBadResponse", Message: "bad response", HTTPStatus: 500}
	ErrWebrpcServerPanic        = WebRPCError{Code: -6, Name: "WebrpcServerPanic", Message: "server panic", HTTPStatus: 500}
	ErrWebrpcInternalError      = WebRPCError{Code: -7, Name: "WebrpcInternalError", Message: "internal error", HTTPStatus: 500}
	ErrWebrpcClientDisconnected = WebRPCError{Code: -8, Name: "WebrpcClientDisconnected", Message: "client disconnected", HTTPStatus: 400}
	ErrWebrpcStreamLost         = WebRPCError{Code: -9, Name: "WebrpcStreamLost", Message: "stream lost", HTTPStatus: 400}
	ErrWebrpcStreamFinished     = WebRPCError{Code: -10, Name: "WebrpcStreamFinished", Message: "stream finished", HTTPStatus: 200}
)

// Schema errors
var (
	ErrEmptyUsername = WebRPCError{Code: 100, Name: "EmptyUsername", Message: "Username must be provided.", HTTPStatus: 400}
)
