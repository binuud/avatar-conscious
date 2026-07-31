/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as Brain from "./brain.pb"
export class Conscious {
  static In(req: Brain.Request, initReq?: fm.InitReq): Promise<Brain.Response> {
    return fm.fetchReq<Brain.Request, Brain.Response>(`/v1/brainRunner/run`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Health(req: Brain.HealthRequest, initReq?: fm.InitReq): Promise<Brain.HealthResponse> {
    return fm.fetchReq<Brain.HealthRequest, Brain.HealthResponse>(`/health?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
}