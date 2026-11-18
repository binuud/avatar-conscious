/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as ConsciousMemory from "./memory.pb"
export class MemoryService {
  static List(req: ConsciousMemory.ListMemoryRequest, initReq?: fm.InitReq): Promise<ConsciousMemory.ListMemoryResponse> {
    return fm.fetchReq<ConsciousMemory.ListMemoryRequest, ConsciousMemory.ListMemoryResponse>(`/v1/memory/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Add(req: ConsciousMemory.AddMemoryRequest, initReq?: fm.InitReq): Promise<ConsciousMemory.AddMemoryResponse> {
    return fm.fetchReq<ConsciousMemory.AddMemoryRequest, ConsciousMemory.AddMemoryResponse>(`/v1/memory/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Update(req: ConsciousMemory.UpdateMemoryRequest, initReq?: fm.InitReq): Promise<ConsciousMemory.UpdateMemoryResponse> {
    return fm.fetchReq<ConsciousMemory.UpdateMemoryRequest, ConsciousMemory.UpdateMemoryResponse>(`/v1/memory/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}