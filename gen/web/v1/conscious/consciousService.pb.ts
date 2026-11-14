/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as ConsciousCommon from "./common.pb"
import * as ConsciousContextMessage from "./contextMessage.pb"
export class Conscious {
  static In(req: ConsciousContextMessage.ContextMesssage, initReq?: fm.InitReq): Promise<ConsciousContextMessage.ContextMesssage> {
    return fm.fetchReq<ConsciousContextMessage.ContextMesssage, ConsciousContextMessage.ContextMesssage>(`/v1/conscious/in`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Generate(req: ConsciousContextMessage.ContextMesssage, entityNotifier?: fm.NotifyStreamEntityArrival<ConsciousContextMessage.ContextMesssage>, initReq?: fm.InitReq): Promise<void> {
    return fm.fetchStreamingRequest<ConsciousContextMessage.ContextMesssage, ConsciousContextMessage.ContextMesssage>(`/v1/conscious/generate`, entityNotifier, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static InStream(req: ConsciousContextMessage.ContextMesssage, entityNotifier?: fm.NotifyStreamEntityArrival<ConsciousContextMessage.ContextMesssage>, initReq?: fm.InitReq): Promise<void> {
    return fm.fetchStreamingRequest<ConsciousContextMessage.ContextMesssage, ConsciousContextMessage.ContextMesssage>(`/v1/conscious/inStream`, entityNotifier, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetModels(req: ConsciousCommon.ModelRequest, initReq?: fm.InitReq): Promise<ConsciousCommon.ModelResponse> {
    return fm.fetchReq<ConsciousCommon.ModelRequest, ConsciousCommon.ModelResponse>(`/v1/conscious/models?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static FileUpload(req: ConsciousContextMessage.UploadRequest, initReq?: fm.InitReq): Promise<ConsciousContextMessage.UploadResponse> {
    return fm.fetchReq<ConsciousContextMessage.UploadRequest, ConsciousContextMessage.UploadResponse>(`/v1/conscious/fileUpload`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ClearMemory(req: ConsciousCommon.ClearMemoryRequest, initReq?: fm.InitReq): Promise<ConsciousCommon.ClearMemoryResponse> {
    return fm.fetchReq<ConsciousCommon.ClearMemoryRequest, ConsciousCommon.ClearMemoryResponse>(`/v1/conscious/clearMemory`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Health(req: ConsciousCommon.HealthRequest, initReq?: fm.InitReq): Promise<ConsciousCommon.HealthResponse> {
    return fm.fetchReq<ConsciousCommon.HealthRequest, ConsciousCommon.HealthResponse>(`/v1/health?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
}