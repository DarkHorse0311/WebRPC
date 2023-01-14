// example  v0.0.1 efa279ca1879bb4cf8f8189317d349929c55ba87
// --
// Code generated by webrpc-gen@v0.10.x-dev with javascript generator. DO NOT EDIT.
//
// webrpc-gen -schema=example.webrpc.json -target=javascript -client -out=./client/client.gen.mjs

// WebRPC description and code-gen version
export const WebRPCVersion = "v1"

// Schema version of your RIDL schema
export const WebRPCSchemaVersion = " v0.0.1"

// Schema hash generated from your RIDL schema
export const WebRPCSchemaHash = "efa279ca1879bb4cf8f8189317d349929c55ba87"

//
// Types
//

export var Kind;
(function (Kind) {
  Kind["USER"] = "USER"
  Kind["ADMIN"] = "ADMIN"
})(Kind || (Kind = {}))

export class Empty {
  constructor(_data) {
    this._data = {}
    if (_data) {
      
    }
  }
  
  toJSON() {
    return this._data
  }
}

export class GetUserRequest {
  constructor(_data) {
    this._data = {}
    if (_data) {
      this._data['userID'] = _data['userID']
      
    }
  }
  get userID() {
    return this._data['userID']
  }
  set userID(value) {
    this._data['userID'] = value
  }
  
  toJSON() {
    return this._data
  }
}

export class User {
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

export class RandomStuff {
  constructor(_data) {
    this._data = {}
    if (_data) {
      this._data['meta'] = _data['meta']
      this._data['metaNestedExample'] = _data['metaNestedExample']
      this._data['namesList'] = _data['namesList']
      this._data['numsList'] = _data['numsList']
      this._data['doubleArray'] = _data['doubleArray']
      this._data['listOfMaps'] = _data['listOfMaps']
      this._data['listOfUsers'] = _data['listOfUsers']
      this._data['mapOfUsers'] = _data['mapOfUsers']
      this._data['user'] = _data['user']
      
    }
  }
  get meta() {
    return this._data['meta']
  }
  set meta(value) {
    this._data['meta'] = value
  }
  get metaNestedExample() {
    return this._data['metaNestedExample']
  }
  set metaNestedExample(value) {
    this._data['metaNestedExample'] = value
  }
  get namesList() {
    return this._data['namesList']
  }
  set namesList(value) {
    this._data['namesList'] = value
  }
  get numsList() {
    return this._data['numsList']
  }
  set numsList(value) {
    this._data['numsList'] = value
  }
  get doubleArray() {
    return this._data['doubleArray']
  }
  set doubleArray(value) {
    this._data['doubleArray'] = value
  }
  get listOfMaps() {
    return this._data['listOfMaps']
  }
  set listOfMaps(value) {
    this._data['listOfMaps'] = value
  }
  get listOfUsers() {
    return this._data['listOfUsers']
  }
  set listOfUsers(value) {
    this._data['listOfUsers'] = value
  }
  get mapOfUsers() {
    return this._data['mapOfUsers']
  }
  set mapOfUsers(value) {
    this._data['mapOfUsers'] = value
  }
  get user() {
    return this._data['user']
  }
  set user(value) {
    this._data['user'] = value
  }
  
  toJSON() {
    return this._data
  }
}

  
//
// Client
//

export class ExampleService {
  constructor(hostname, fetch) {
    this.path = '/rpc/ExampleService/'
    this.hostname = hostname
    this.fetch = (input, init) => fetch(input, init)
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
