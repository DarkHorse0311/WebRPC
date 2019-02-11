// Code generated by statik. DO NOT EDIT.

// Package contains static assets.
package embed

var	Asset = "PK\x03\x04\x14\x00\x08\x00\x00\x00\xf4\x90HN\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00client.ts.tmplUT\x05\x00\x01m\xc5]\\{{define \"client\"}}\n\n{{- if .Services}}\n// Client\n{{range .Services}}\nconst {{.Name | constPathPrefix}} = \"/rpc/{{.Name}}/\"\n{{end}}\n\n{{- range .Services}}\nexport class {{.Name}} implements {{.Name | serviceInterfaceName}} {\n  private hostname: string\n  private fetch: Fetch\n  private path = '/rpc/{{.Name}}/'\n\n  constructor(hostname: string, fetch: Fetch) {\n    this.hostname = hostname\n    this.fetch = fetch\n  }\n\n  private url(name: string): string {\n    return this.hostname + this.path + name\n  }\n\n  {{range .Methods}}\n  {{.Name}}({{.Inputs | methodInputs}} = {}): {{.Outputs | methodOutputs}} {\n    return this.fetch(\n      this.url('{{.Name}}'),\n      {{if .Inputs | len}}\n      createHTTPRequest(params, headers)\n      {{else}}\n      createHTTPRequest({}, headers)\n      {{end}}\n    ).then((res) => {\n      if (!res.ok) {\n        return throwHTTPError(res)\n      }\n      {{range $output := .Outputs}}\n      return res.json().then((_data) => {return {{$output | newResponseConcreteType}}(_data)})\n      {{end}}\n    })\n  }\n  {{end}}\n}\n{{end -}}\n\n{{end -}}\n{{end}}\nPK\x07\x08\xd1\xbb]\x0f-\x04\x00\x00-\x04\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\xf4\x90HN\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x00	\x00helpers.ts.tmplUT\x05\x00\x01m\xc5]\\{{define \"helpers\"}}\n\nexport interface WebRPCError extends Error {\n  code: string\n  msg: string\n	status: number\n}\n\nexport const throwHTTPError = (resp: Response) => {\n  return resp.json().then((err: WebRPCError) => { throw err })\n}\n\nexport const createHTTPRequest = (body: object = {}, headers: object = {}): object => {\n  return {\n    method: 'POST',\n    headers: { ...headers, 'Content-Type': 'application/json' },\n    body: JSON.stringify(body || {})\n  }\n}\n\nexport type Fetch = (input: RequestInfo, init?: RequestInit) => Promise<Response>\n{{end}}\nPK\x07\x08\x90d\x1d='\x02\x00\x00'\x02\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00#\x89JN\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x11\x00	\x00proto.gen.ts.tmplUT\x05\x00\x01\xb2Z`\\{{define \"proto\"}}\n/* tslint:disable */\n\n// This file has been generated by https://github.com/webrpc/webrpc using gen/typescript\n// Do not edit by hand. Update your webrpc schema and re-generate.\n\n{{template \"types\" .}}\n{{- if .TargetOpts.Client}}\n  {{template \"client\" .}}\n{{- end}}\n{{- if .TargetOpts.Server}}\n  {{template \"server\" .}}\n{{- end}}\n{{template \"helpers\" .}}\n{{end}}\nPK\x07\x08\xd2\x8d\x1f\xc5~\x01\x00\x00~\x01\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00;tHN\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00server.ts.tmplUT\x05\x00\x01R\x93]\\{{define \"server\"}}\n{{- if .Services}}\n// TODO: Server\n{{end -}}\n{{end}}\nPK\x07\x08\x8a@[\xefI\x00\x00\x00I\x00\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00/\x0fKN\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0d\x00	\x00types.ts.tmplUT\x05\x00\x01\x8a\xd6`\\{{define \"types\"}}\n\n{{- if .Messages -}}\n{{range .Messages -}}\n\n{{if .Type | isEnum -}}\n{{$enumName := .Name}}\nexport enum {{$enumName}} {\n{{- range $i, $field := .Fields}}\n  {{- if $i}},{{end}}\n  {{$field.Name}} = '{{$field.Name}}'\n{{- end}}\n}\n{{end -}}\n\n{{- if .Type | isStruct  }}\nexport interface {{.Name | interfaceName}} {\n  {{range .Fields -}}\n  {{. | exportedJSONField}}{{if .Optional}}?{{end}}: {{.Type | fieldType}}\n  {{end -}}\n\n  toJSON?(): object\n}\n\nexport class {{.Name}} implements {{.Name | interfaceName}} {\n  private _data: {{.Name | interfaceName}}\n  constructor(_data?: {{.Name | interfaceName}}) {\n    this._data = {} as {{.Name | interfaceName}}\n    if (_data) {\n      {{range .Fields -}}\n      this._data['{{. | exportedJSONField}}'] = _data['{{. | exportedJSONField}}']!\n      {{end}}\n    }\n  }\n  {{ range .Fields -}}\n  public get {{. | exportedJSONField}}(): {{.Type | fieldType}} {\n    return this._data['{{. | exportedJSONField }}']!\n  }\n  public set {{. | exportedJSONField}}(value: {{.Type | fieldType}}) {\n    this._data['{{. | exportedJSONField}}'] = value\n  }\n  {{end}}\n  public toJSON(): object {\n    return this._data\n  }\n}\n{{end -}}\n{{end -}}\n{{end -}}\n\n{{if .Services}}\n{{range .Services -}}\nexport interface {{.Name | serviceInterfaceName}} {\n  {{range .Methods -}}\n    {{.Name}}({{.Inputs | methodInputs}}): {{.Outputs | methodOutputs}}\n  {{end}}\n}\n{{- end}}\n{{end -}}\n\n{{end}}\nPK\x07\x08w\xefkl\x87\x05\x00\x00\x87\x05\x00\x00PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xf4\x90HN\xd1\xbb]\x0f-\x04\x00\x00-\x04\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x00\x00\x00\x00client.ts.tmplUT\x05\x00\x01m\xc5]\\PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\xf4\x90HN\x90d\x1d='\x02\x00\x00'\x02\x00\x00\x0f\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81r\x04\x00\x00helpers.ts.tmplUT\x05\x00\x01m\xc5]\\PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00#\x89JN\xd2\x8d\x1f\xc5~\x01\x00\x00~\x01\x00\x00\x11\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\xdf\x06\x00\x00proto.gen.ts.tmplUT\x05\x00\x01\xb2Z`\\PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00;tHN\x8a@[\xefI\x00\x00\x00I\x00\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\xa5\x08\x00\x00server.ts.tmplUT\x05\x00\x01R\x93]\\PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00/\x0fKNw\xefkl\x87\x05\x00\x00\x87\x05\x00\x00\x0d\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x813	\x00\x00types.ts.tmplUT\x05\x00\x01\x8a\xd6`\\PK\x05\x06\x00\x00\x00\x00\x05\x00\x05\x00\\\x01\x00\x00\xfe\x0e\x00\x00\x00\x00"
