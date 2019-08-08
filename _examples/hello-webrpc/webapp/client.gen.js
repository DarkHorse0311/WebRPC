// hello-webrpc v1.0.0 c929128d878e94653f3a856f80c4671008e22a45
// --
// This file has been generated by https://github.com/webrpc/webrpc using gen/javascript
// Do not edit by hand. Update your webrpc schema and re-generate.

// WebRPC description and code-gen version
export const WebRPCVersion = "v1"

// Schema version of your RIDL schema
export const WebRPCSchemaVersion = "v1.0.0"

// Schema hash generated from your RIDL schema
export const WebRPCSchemaHash = "c929128d878e94653f3a856f80c4671008e22a45"


//
// Types
//

var Kind;
(function (Kind) {
  Kind["USER"] = "USER"
  Kind["ADMIN"] = "ADMIN"
})(Kind || (Kind = {}))

class Empty {
  constructor(_data) {
    this._data = {}
    if (_data) {
      
    }
  }
  
  toJSON() {
    return this._data
  }
}

class User {
  constructor(_data) {
    this._data = {}
    if (_data) {
      this._data['id'] = _data['id']
      this._data['USERNAME'] = _data['USERNAME']
      this._data['created_at'] = _data['created_at']
      
    }
  }
  get id() {
    return this._data['id']
  }
  set id(value) {
    this._data['id'] = value
  }
  get USERNAME() {
    return this._data['USERNAME']
  }
  set USERNAME(value) {
    this._data['USERNAME'] = value
  }
  get created_at() {
    return this._data['created_at']
  }
  set created_at(value) {
    this._data['created_at'] = value
  }
  
  toJSON() {
    return this._data
  }
}

  
//
// Client
//

class ExampleService {
  constructor(hostname, fetch) {
    this.path = '/rpc/ExampleService/'
    this.hostname = hostname
    this.fetch = fetch
  }

  url(name) {
    return this.hostname + this.path + name
  }
  
  ping = (headers) => {
    return this.fetch(
      this.url('Ping'),
      createHTTPRequest({}, headers)
    ).then((res) => {
      return buildResponse(res).then(_data => {
        return {
          status: (_data.status)
        }
      })
    })
  }
  
  getUser = (args, headers) => {
    return this.fetch(
      this.url('GetUser'),
      createHTTPRequest(args, headers)
    ).then((res) => {
      return buildResponse(res).then(_data => {
        return {
          user: new User(_data.user)
        }
      })
    })
  }
  
}

  
const createHTTPRequest = (body = {}, headers = {}) => {
  return {
    method: 'POST',
    headers: { ...headers, 'Content-Type': 'application/json' },
    body: JSON.stringify(body || {})
  }
}

const buildResponse = (res) => {
  return res.text().then(text => {
    let data
    try {
      data = JSON.parse(text)
    } catch(err) {
      throw { code: 'unknown', msg: `expecting JSON, got: ${text}`, status: res.status }
    }
    if (!res.ok) {
      throw data // webrpc error response
    }
    return data
  })
}
