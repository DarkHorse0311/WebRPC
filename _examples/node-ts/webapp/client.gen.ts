/* tslint:disable */
// node-ts v1.0.0 4d2858fa129683e5775e9b863ceceb740e7e09b1
// --
// This file has been generated by https://github.com/webrpc/webrpc using gen/typescript
// Do not edit by hand. Update your webrpc schema and re-generate.

// WebRPC description and code-gen version
export const WebRPCVersion = "v1"

// Schema version of your RIDL schema
export const WebRPCSchemaVersion = "v1.0.0"

// Schema hash generated from your RIDL schema
export const WebRPCSchemaHash = "4d2858fa129683e5775e9b863ceceb740e7e09b1"


//
// Types
//
export enum Kind {
  USER = 'USER',
  ADMIN = 'ADMIN'
}

export interface User {
  id: number
  USERNAME: string
  role: Kind
  meta: {[key: string]: any}
  
  createdAt?: string
}

export interface Page {
  num: number
}

export interface ExampleService {
  ping(headers?: object): Promise<PingReturn>
  getUser(args: GetUserArgs, headers?: object): Promise<GetUserReturn>
}

export interface PingArgs {
}

export interface PingReturn {
  status: boolean  
}
export interface GetUserArgs {
  userID: number
}

export interface GetUserReturn {
  code: number
  user: User  
}


  
//
// Client
//
export class ExampleService implements ExampleService {
  private hostname: string
  private fetch: Fetch
  private path = '/rpc/ExampleService/'

  constructor(hostname: string, fetch: Fetch) {
    this.hostname = hostname
    this.fetch = fetch
  }

  private url(name: string): string {
    return this.hostname + this.path + name
  }
  
  ping = (headers?: object): Promise<PingReturn> => {
    return this.fetch(
      this.url('Ping'),
      createHTTPRequest({}, headers)
      ).then((res) => {
      return buildResponse(res).then(_data => {
        return {
          status: <boolean>(_data.status)
        }
      })
    })
  }
  
  getUser = (args: GetUserArgs, headers?: object): Promise<GetUserReturn> => {
    return this.fetch(
      this.url('GetUser'),
      createHTTPRequest(args, headers)).then((res) => {
      return buildResponse(res).then(_data => {
        return {
          code: <number>(_data.code), 
          user: <User>(_data.user)
        }
      })
    })
  }
  
}

  
export interface WebRPCError extends Error {
  code: string
  msg: string
	status: number
}

const createHTTPRequest = (body: object = {}, headers: object = {}): object => {
  return {
    method: 'POST',
    headers: { ...headers, 'Content-Type': 'application/json' },
    body: JSON.stringify(body || {})
  }
}

const buildResponse = (res: Response): Promise<any> => {
  return res.text().then(text => {
    let data
    try {
      data = JSON.parse(text)
    } catch(err) {
      throw { code: 'unknown', msg: `expecting JSON, got: ${text}`, status: res.status } as WebRPCError
    }
    if (!res.ok) {
      throw data // webrpc error response
    }
    return data
  })
}

export type Fetch = (input: RequestInfo, init?: RequestInit) => Promise<Response>
