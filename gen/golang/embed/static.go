// Code generated by statik. DO NOT EDIT.

// Package contains static assets.
package embed

var	Asset = "PK\x03\x04\x14\x00\x08\x00\x00\x00`\xbfMN\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00client.go.tmplUT\x05\x00\x01D\xafd\\{{define \"client\"}}\n{{if .Services}}\n  // Client\n\n  {{range .Services}}\n  const {{.Name | constPathPrefix}} = \"/rpc/{{.Name}}/\"\n  {{end}}\n\n  {{range .Services}}\n    {{ $serviceName := .Name | clientServiceName}}\n    type {{$serviceName}} struct {\n      client HTTPClient\n      urls   [{{.Methods | countMethods}}]string\n    }\n\n    func {{.Name | newClientServiceName }}(addr string, client HTTPClient) {{.Name}} {\n      prefix := urlBase(addr) + {{.Name | constPathPrefix}}\n      urls := [{{.Methods | countMethods}}]string{\n        {{- range .Methods}}\n        prefix + \"{{.Name}}\",\n        {{- end}}\n      }\n      return &{{$serviceName}}{\n        client: client,\n        urls:   urls,\n      }\n    }\n\n    {{range $i, $method := .Methods}}\n      func (c *{{$serviceName}}) {{.Name}}({{.Inputs | methodInputs}}) ({{.Outputs | methodOutputs }}) {\n        {{- range .Outputs}}\n          out := new({{.Type | fieldType}})\n        {{- end}}\n        {{- if .Inputs | len}}\n        err := doJSONRequest(ctx, c.client, c.urls[{{$i}}], {{index .Inputs 0 | methodInputName}}, out)\n        {{- else}}\n        err := doJSONRequest(ctx, c.client, c.urls[{{$i}}], nil, out)\n        {{- end}}\n        if err != nil {\n          return nil, err\n        }\n        return out, nil\n      }\n    {{end}}\n  {{end}}\n{{end}}\n{{end}}\nPK\x07\x08\xc2\xa3X\xb7\x1d\x05\x00\x00\x1d\x05\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x87\xb5GN\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x00	\x00helpers.go.tmplUT\x05\x00\x01\xbf\xb4\\\\{{define \"helpers\"}}\n//\n// Helpers\n//\n\n// HTTPClient is the interface used by generated clients to send HTTP requests.\n// It is fulfilled by *(net/http).Client, which is sufficient for most users.\n// Users can provide their own implementation for special retry policies.\ntype HTTPClient interface {\n  Do(req *http.Request) (*http.Response, error)\n}\n\ntype WebRPCServer interface {\n  http.Handler\n  WebRPCVersion() string\n  ServiceVersion() string\n}\n\ntype errResponse struct {\n  Status int    `json:\"status\"`\n  Code   string `json:\"code\"`\n  Msg    string `json:\"msg\"`\n  Cause  string `json:\"cause,omitempty\"`\n}\n\nfunc writeJSONError(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {\n  rpcErr, ok := err.(webrpc.Error)\n  if !ok {\n    rpcErr = webrpc.WrapError(webrpc.ErrInternal, err, \"webrpc error\")\n  }\n\n  statusCode := webrpc.HTTPStatusFromErrorCode(rpcErr.Code())\n\n  w.Header().Set(\"Content-Type\", \"application/json\")\n  w.WriteHeader(statusCode)\n\n  errResp := errResponse{\n    Status: statusCode,\n    Code:   string(rpcErr.Code()),\n    Msg:    rpcErr.Error(),\n  }\n  respBody, _ := json.Marshal(errResp)\n  w.Write(respBody)\n}\n\n// urlBase helps ensure that addr specifies a scheme. If it is unparsable\n// as a URL, it returns addr unchanged.\nfunc urlBase(addr string) string {\n  // If the addr specifies a scheme, use it. If not, default to\n  // http. If url.Parse fails on it, return it unchanged.\n  url, err := url.Parse(addr)\n  if err != nil {\n    return addr\n  }\n  if url.Scheme == \"\" {\n    url.Scheme = \"http\"\n  }\n  return url.String()\n}\n\n// newRequest makes an http.Request from a client, adding common headers.\nfunc newRequest(ctx context.Context, url string, reqBody io.Reader, contentType string) (*http.Request, error) {\n  req, err := http.NewRequest(\"POST\", url, reqBody)\n  if err != nil {\n    return nil, err\n  }\n  req.Header.Set(\"Accept\", contentType)\n  req.Header.Set(\"Content-Type\", contentType)\n  return req, nil\n}\n\n// doJSONRequest is common code to make a request to the remote service.\nfunc doJSONRequest(ctx context.Context, client HTTPClient, url string, in, out interface{}) error {\n	reqBody, err := json.Marshal(in)\n	if err != nil {\n		return clientError(\"failed to marshal json request\", err)\n	}\n	if err = ctx.Err(); err != nil {\n		return clientError(\"aborted because context was done\", err)\n	}\n\n	req, err := newRequest(ctx, url, bytes.NewBuffer(reqBody), \"application/json\")\n	if err != nil {\n		return clientError(\"could not build request\", err)\n	}\n	resp, err := client.Do(req)\n	if err != nil {\n		return clientError(\"request failed\", err)\n	}\n\n	defer func() {\n		cerr := resp.Body.Close()\n		if err == nil && cerr != nil {\n			err = clientError(\"failed to close response body\", cerr)\n		}\n	}()\n\n	if err = ctx.Err(); err != nil {\n		return clientError(\"aborted because context was done\", err)\n	}\n\n	if resp.StatusCode != 200 {\n		return errorFromResponse(resp)\n	}\n\n	respBody, err := ioutil.ReadAll(resp.Body)\n	if err != nil {\n		return clientError(\"failed to read response body\", err)\n	}\n\n	err = json.Unmarshal(respBody, &out)\n	if err != nil {\n		return clientError(\"failed to unmarshal json response body\", err)\n	}\n	if err = ctx.Err(); err != nil {\n		return clientError(\"aborted because context was done\", err)\n	}\n\n	return nil\n}\n\n// errorFromResponse builds a webrpc.Error from a non-200 HTTP response.\nfunc errorFromResponse(resp *http.Response) webrpc.Error {\n	respBody, err := ioutil.ReadAll(resp.Body)\n	if err != nil {\n		return clientError(\"failed to read server error response body\", err)\n	}\n\n	var respErr errResponse\n	if err := json.Unmarshal(respBody, &respErr); err != nil {\n		return clientError(\"failed unmarshal error response\", err)\n	}\n\n	errCode := webrpc.ErrorCode(respErr.Code)\n\n	if webrpc.HTTPStatusFromErrorCode(errCode) == 0 {\n		return webrpc.ErrorInternal(\"invalid code returned from server error response: %s\", respErr.Code)\n	}\n\n	return webrpc.Errorf(errCode, respErr.Msg)\n}\n\nfunc clientError(desc string, err error) webrpc.Error {\n	return webrpc.WrapError(webrpc.ErrInternal, err, desc)\n}\n{{end}}\nPK\x07\x08\x1a\xc0\x05B\xd3\x0f\x00\x00\xd3\x0f\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\xb3\x83JN\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x11\x00	\x00proto.gen.go.tmplUT\x05\x00\x01sQ`\\{{define \"proto\"}}\n// This file has been generated by https://github.com/webrpc/webrpc using gen/golang\n// Do not edit by hand. Update your webrpc schema and re-generate.\npackage {{.TargetOpts.PkgName}}\n\nimport (\n  \"bytes\"\n  \"context\"\n  \"encoding/json\"\n  \"io\"\n  \"io/ioutil\"\n  \"net/http\"\n  \"net/url\"\n  \"strings\"\n  \"time\"\n\n  \"github.com/webrpc/webrpc/lib/webrpc-go\"\n)\n\n{{template \"types\" .}}\n\n{{if .TargetOpts.Client}}\n  {{template \"client\" .}}\n{{end}}\n\n{{if .TargetOpts.Server}}\n  {{template \"server\" .}}\n{{end}}\n\n{{template \"helpers\" .}}\n\n{{end}}\nPK\x07\x08\x8c\xe4DY#\x02\x00\x00#\x02\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x99{GN\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00server.go.tmplUT\x05\x00\x01\xb3N\\\\{{define \"server\"}}\n{{if .Services}}\n  // Server\n  {{- range .Services}}\n    {{$name := .Name}}\n    {{$serviceName := .Name | serverServiceName}}\n\n    type {{$serviceName}} struct {\n      {{.Name}}\n    }\n\n    func {{ .Name | newServerServiceName }}(svc {{.Name}}) WebRPCServer {\n      return &{{$serviceName}}{\n        {{.Name}}: svc,\n      }\n    }\n\n\n    func (s *{{$serviceName}}) WebRPCVersion() string {\n      return \"v0.0.1\"\n    }\n\n    func (s *{{$serviceName}}) ServiceVersion() string {\n      return \"v0.1.0\"\n    }\n\n    func (s *{{$serviceName}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {\n      ctx := r.Context()\n      ctx = webrpc.WithResponseWriter(ctx, w)\n      ctx = webrpc.WithServiceName(ctx, \"{{.Name}}\")\n\n      if r.Method != \"POST\" {\n        err := webrpc.Errorf(webrpc.ErrBadRoute, \"unsupported method %q (only POST is allowed)\", r.Method)\n        writeJSONError(ctx, w, r, err)\n        return\n      }\n\n      switch r.URL.Path {\n      {{- range .Methods}}\n      case \"/rpc/{{$name}}/{{.Name}}\":\n        s.{{.Name | serviceMethodName}}(ctx, w, r)\n        return\n      {{- end}}\n      default:\n        err := webrpc.Errorf(webrpc.ErrBadRoute, \"no handler for path %q\", r.URL.Path)\n        writeJSONError(ctx, w, r, err)\n        return\n      }\n    }\n\n    {{range .Methods}}\n      func (s *{{$serviceName}}) {{.Name | serviceMethodName}}(ctx context.Context, w http.ResponseWriter, r *http.Request) {\n        header := r.Header.Get(\"Content-Type\")\n        i := strings.Index(header, \";\")\n        if i == -1 {\n          i = len(header)\n        }\n\n        switch strings.TrimSpace(strings.ToLower(header[:i])) {\n        case \"application/json\":\n          s.{{ .Name | serviceMethodJSONName }}(ctx, w, r)\n        default:\n          err := webrpc.Errorf(webrpc.ErrBadRoute, \"unexpected Content-Type: %q\", r.Header.Get(\"Content-Type\"))\n          writeJSONError(ctx, w, r, err)\n        }\n      }\n\n      func (s *{{$serviceName}}) {{.Name | serviceMethodJSONName}}(ctx context.Context, w http.ResponseWriter, r *http.Request) {\n        var err error\n        ctx = webrpc.WithMethodName(ctx, \"{{.Name}}\")\n\n        {{- if .Inputs|len}}\n\n        {{- range .Inputs}}\n        reqContent := new({{.Type | fieldType}})\n        {{- end}}\n\n        reqBody, err := ioutil.ReadAll(r.Body)\n        if err != nil {\n          err = webrpc.WrapError(webrpc.ErrInternal, err, \"failed to read request data\")\n          writeJSONError(ctx, w, r, err)\n          return\n        }\n        defer r.Body.Close()\n\n        err = json.Unmarshal(reqBody, reqContent)\n        if err != nil {\n          err = webrpc.WrapError(webrpc.ErrInternal, err, \"failed to unmarshal request data\")\n          writeJSONError(ctx, w, r, err)\n          return\n        }\n        {{- end}}\n\n        // Call service method\n        {{- range .Outputs}}\n        var respContent *{{.Type | fieldType}}\n        {{- end}}\n        func() {\n          defer func() {\n            // In case of a panic, serve a 500 error and then panic.\n            if rr := recover(); rr != nil {\n              writeJSONError(ctx, w, r, webrpc.ErrorInternal(\"internal service panic\"))\n              panic(rr)\n            }\n          }()\n          {{- if .Inputs|len}}\n          respContent, err = s.{{$name}}.{{.Name}}(ctx, reqContent)\n          {{- else}}\n          respContent, err = s.{{$name}}.{{.Name}}(ctx)\n          {{- end}}\n        }()\n\n        if err != nil {\n          writeJSONError(ctx, w, r, err)\n          return\n        }\n\n        respBody, err := json.Marshal(respContent)\n        if err != nil {\n          err = webrpc.WrapError(webrpc.ErrInternal, err, \"failed to marshal json response\")\n          writeJSONError(ctx, w, r, err)\n          return\n        }\n\n        w.Header().Set(\"Content-Type\", \"application/json\")\n        w.WriteHeader(http.StatusOK)\n\n        if n, err := w.Write(respBody); err != nil {\n          _ = n\n          _ = err\n          // TODO: failing silently for now..\n          // msg := fmt.Sprintf(\"failed to write response, %d of %d bytes written: %s\", n, len(respBytes), err.Error())\n        }\n      }\n    {{end}}\n  {{- end}}\n{{end}}\n{{end}}\nPK\x07\x08\x93[\xb4e\x1f\x10\x00\x00\x1f\x10\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\xfdpNN\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0d\x00	\x00types.go.tmplUT\x05\x00\x01?ve\\{{define \"types\"}}\n\n{{if .Messages}}\n  {{range .Messages}}\n    {{if .Type | isEnum}}\n      {{$enumName := .Name}}\n      {{$enumType := .EnumType}}\n      type {{$enumName}} {{$enumType}}\n\n      const (\n        {{- range .Fields}}\n          {{$enumName}}_{{.Name}} {{$enumName}} = {{.Value}}\n        {{- end}}\n      )\n\n      var {{$enumName}}_name = map[{{$enumType}}]string {\n        {{- range .Fields}}\n          {{.Value}}: \"{{.Name}}\",\n        {{- end}}\n      }\n\n      var {{$enumName}}_value = map[string]{{$enumType}} {\n        {{- range .Fields}}\n          \"{{.Name}}\": {{.Value}},\n        {{- end}}\n      }\n\n      func (x {{$enumName}}) String() string {\n        return {{$enumName}}_name[{{$enumType}}(x)]\n      }\n\n      func (x {{$enumName}}) MarshalJSON() ([]byte, error) {\n        buf := bytes.NewBufferString(`\"`)\n        buf.WriteString({{$enumName}}_name[{{$enumType}}(x)])\n        buf.WriteString(`\"`)\n        return buf.Bytes(), nil\n      }\n\n      func (x *{{$enumName}}) UnmarshalJSON(b []byte) error {\n        var j string\n        err := json.Unmarshal(b, &j)\n        if err != nil {\n          return err\n        }\n        *x = {{$enumName}}({{$enumName}}_value[j])\n        return nil\n      }\n    {{end}}\n    {{if .Type | isStruct  }}\n      type {{.Name}} struct {\n        {{- range .Fields}}\n          {{.Name | exportedField}} {{. | fieldOptional}}{{.Type | fieldType}} {{.Meta | fieldTags}}\n        {{- end}}\n      }\n    {{end}}\n  {{end}}\n{{end}}\n{{if .Services}}\n  {{range .Services}}\n    type {{.Name}} interface {\n      {{- range .Methods}}\n        {{.Name}}({{.Inputs | methodInterfaceInputs}}) ({{.Outputs | methodOutputs}})\n      {{- end}}\n    }\n  {{end}}\n  var Services = map[string][]string{\n    {{- range .Services}}\n      \"{{.Name}}\": {\n        {{- range .Methods}}\n          \"{{.Name}}\",\n        {{- end}}\n      },\n    {{- end}}\n  }\n{{end}}\n\n{{end}}\nPK\x07\x08\xf7\xf6\xfd.Y\x07\x00\x00Y\x07\x00\x00PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00`\xbfMN\xc2\xa3X\xb7\x1d\x05\x00\x00\x1d\x05\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x00\x00\x00\x00client.go.tmplUT\x05\x00\x01D\xafd\\PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x87\xb5GN\x1a\xc0\x05B\xd3\x0f\x00\x00\xd3\x0f\x00\x00\x0f\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81b\x05\x00\x00helpers.go.tmplUT\x05\x00\x01\xbf\xb4\\\\PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xb3\x83JN\x8c\xe4DY#\x02\x00\x00#\x02\x00\x00\x11\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81{\x15\x00\x00proto.gen.go.tmplUT\x05\x00\x01sQ`\\PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x99{GN\x93[\xb4e\x1f\x10\x00\x00\x1f\x10\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\xe6\x17\x00\x00server.go.tmplUT\x05\x00\x01\xb3N\\\\PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xfdpNN\xf7\xf6\xfd.Y\x07\x00\x00Y\x07\x00\x00\x0d\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81J(\x00\x00types.go.tmplUT\x05\x00\x01?ve\\PK\x05\x06\x00\x00\x00\x00\x05\x00\x05\x00\\\x01\x00\x00\xe7/\x00\x00\x00\x00"
