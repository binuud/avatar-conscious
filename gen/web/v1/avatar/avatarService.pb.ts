/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as AvatarAvatar from "./avatar.pb"
export type GetRequest = {
  uuid?: string
  name?: string
}

export type GetResponse = {
  avatar?: AvatarAvatar.Avatar
}

export type StartRequest = {
  uuid?: string
  name?: string
}

export type StartResponse = {
  avatar?: AvatarAvatar.Avatar
}

export type ControlRequest = {
  control?: AvatarAvatar.AvatarControl
}

export type ControlResponse = {
  control?: AvatarAvatar.AvatarControl
}

export type StopRequest = {
  uuid?: string
  name?: string
}

export type StopResponse = {
  avatar?: AvatarAvatar.Avatar
}

export class AvatarService {
  static Get(req: GetRequest, initReq?: fm.InitReq): Promise<GetResponse> {
    return fm.fetchReq<GetRequest, GetResponse>(`/v1/avatar/get?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static Start(req: StartRequest, initReq?: fm.InitReq): Promise<StartResponse> {
    return fm.fetchReq<StartRequest, StartResponse>(`/v1/avatar/start`, {...initReq, method: "POST"})
  }
  static Stop(req: StopRequest, initReq?: fm.InitReq): Promise<StopResponse> {
    return fm.fetchReq<StopRequest, StopResponse>(`/v1/avatar/stop`, {...initReq, method: "POST"})
  }
}