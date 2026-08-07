/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as AiModuleAiModuleService from "../aiModule/aiModuleService.pb"
import * as Brain from "./brain.pb"
export class Conscious {
  static In(req: Brain.InRequest, initReq?: fm.InitReq): Promise<Brain.InResponse> {
    return fm.fetchReq<Brain.InRequest, Brain.InResponse>(`/v1/conscious/in`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Health(req: AiModuleAiModuleService.HealthRequest, initReq?: fm.InitReq): Promise<AiModuleAiModuleService.HealthResponse> {
    return fm.fetchReq<AiModuleAiModuleService.HealthRequest, AiModuleAiModuleService.HealthResponse>(`/health?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
}