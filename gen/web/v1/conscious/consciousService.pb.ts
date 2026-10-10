/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as Common from "./common.pb"
import * as ContextMessage from "./contextMessage.pb"
export class Conscious {
  static In(req: ContextMessage.contextMesssage, initReq?: fm.InitReq): Promise<ContextMessage.contextMesssage> {
    return fm.fetchReq<ContextMessage.contextMesssage, ContextMessage.contextMesssage>(`/v1/conscious/in`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static InStream(req: ContextMessage.contextMesssage, entityNotifier?: fm.NotifyStreamEntityArrival<ContextMessage.contextMesssage>, initReq?: fm.InitReq): Promise<void> {
    return fm.fetchStreamingRequest<ContextMessage.contextMesssage, ContextMessage.contextMesssage>(`/v1/conscious/inStream`, entityNotifier, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Health(req: Common.HealthRequest, initReq?: fm.InitReq): Promise<Common.HealthResponse> {
    return fm.fetchReq<Common.HealthRequest, Common.HealthResponse>(`/v1/health?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
}