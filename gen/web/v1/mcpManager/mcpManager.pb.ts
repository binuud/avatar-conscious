/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as CommonCommon from "../common/common.pb"

export enum MCPTransportType {
  MCP_TYPE_IGNORE = "MCP_TYPE_IGNORE",
  MCP_TYPE_STDIO = "MCP_TYPE_STDIO",
  MCP_TYPE_HTTP = "MCP_TYPE_HTTP",
  MCP_TYPE_SSE = "MCP_TYPE_SSE",
}

export enum MCPServerStatus {
  SERVER_STATUS_IGNORE = "SERVER_STATUS_IGNORE",
  SERVER_STATUS_CONNECTING = "SERVER_STATUS_CONNECTING",
  SERVER_STATUS_CONNECTED = "SERVER_STATUS_CONNECTED",
  SERVER_STATUS_DISCONNECTED = "SERVER_STATUS_DISCONNECTED",
  SERVER_STATUS_ERROR = "SERVER_STATUS_ERROR",
}

export type MCPServerConfig = {
  uuid?: string
  name?: string
  mcpType?: MCPTransportType
  command?: string
  args?: string[]
  env?: {[key: string]: string}
  url?: string
  timeout?: CommonCommon.Timestamp
}

export type MCPTool = {
  name?: string
  description?: string
  inputSchema?: string
}

export type MCPServerConfigs = {
  configs?: MCPServerConfig[]
}

export type MCPServerInstance = {
  config?: MCPServerConfig
  status?: MCPServerStatus
  lastChecked?: CommonCommon.Timestamp
  lastError?: string
  tools?: MCPTool[]
  resources?: string[]
  prompts?: string[]
  prmDoc?: string
}

export type CheckRequest = {
  config?: MCPServerConfig
}

export type CheckResponse = {
  instance?: MCPServerInstance
}

export type RegisterRequest = {
  config?: MCPServerConfig
}

export type RegisterResponse = {
  instance?: MCPServerInstance
}

export type GetRequest = {
  uuid?: string
}

export type GetResponse = {
  instance?: MCPServerInstance
}

export type ListRequest = {
}

export type ListResponse = {
  servers?: MCPServerConfig[]
}

export type DeleteRequest = {
  uuid?: string
}

export type DeleteResponse = {
  success?: boolean
}

export class MCPManagerService {
  static Check(req: CheckRequest, initReq?: fm.InitReq): Promise<CheckResponse> {
    return fm.fetchReq<CheckRequest, CheckResponse>(`/v1/mcpManager/check`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Register(req: RegisterRequest, initReq?: fm.InitReq): Promise<RegisterResponse> {
    return fm.fetchReq<RegisterRequest, RegisterResponse>(`/v1/mcpManager/register`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Get(req: GetRequest, initReq?: fm.InitReq): Promise<GetResponse> {
    return fm.fetchReq<GetRequest, GetResponse>(`/v1/mcpManager/get?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static List(req: ListRequest, initReq?: fm.InitReq): Promise<ListResponse> {
    return fm.fetchReq<ListRequest, ListResponse>(`/v1/mcpManager/list?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static Delete(req: DeleteRequest, initReq?: fm.InitReq): Promise<DeleteResponse> {
    return fm.fetchReq<DeleteRequest, DeleteResponse>(`/v1/mcpManager/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}