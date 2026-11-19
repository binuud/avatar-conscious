/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as ConsciousMemory from "./memory.pb"
export class MemoryService {
  static List(req: ConsciousMemory.ListMemoryRequest, initReq?: fm.InitReq): Promise<ConsciousMemory.ListMemoryResponse> {
    return fm.fetchReq<ConsciousMemory.ListMemoryRequest, ConsciousMemory.ListMemoryResponse>(`/v1/memory/list?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static GetMemoryMeta(req: ConsciousMemory.GetMemoryMetaRequest, initReq?: fm.InitReq): Promise<ConsciousMemory.GetMemoryMetaResponse> {
    return fm.fetchReq<ConsciousMemory.GetMemoryMetaRequest, ConsciousMemory.GetMemoryMetaResponse>(`/v1/memory/getMemoryMeta?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static Get(req: ConsciousMemory.GetMemoryRequest, initReq?: fm.InitReq): Promise<ConsciousMemory.GetMemoryResponse> {
    return fm.fetchReq<ConsciousMemory.GetMemoryRequest, ConsciousMemory.GetMemoryResponse>(`/v1/memory/get?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static Add(req: ConsciousMemory.AddMemoryRequest, initReq?: fm.InitReq): Promise<ConsciousMemory.AddMemoryResponse> {
    return fm.fetchReq<ConsciousMemory.AddMemoryRequest, ConsciousMemory.AddMemoryResponse>(`/v1/memory/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Update(req: ConsciousMemory.UpdateMemoryRequest, initReq?: fm.InitReq): Promise<ConsciousMemory.UpdateMemoryResponse> {
    return fm.fetchReq<ConsciousMemory.UpdateMemoryRequest, ConsciousMemory.UpdateMemoryResponse>(`/v1/memory/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Delete(req: ConsciousMemory.DeleteMemoryRequest, initReq?: fm.InitReq): Promise<ConsciousMemory.DeleteMemoryResponse> {
    return fm.fetchReq<ConsciousMemory.DeleteMemoryRequest, ConsciousMemory.DeleteMemoryResponse>(`/v1/memory/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ClearMemory(req: ConsciousMemory.ClearMemoryRequest, initReq?: fm.InitReq): Promise<ConsciousMemory.ClearMemoryResponse> {
    return fm.fetchReq<ConsciousMemory.ClearMemoryRequest, ConsciousMemory.ClearMemoryResponse>(`/v1/memory/clear`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}