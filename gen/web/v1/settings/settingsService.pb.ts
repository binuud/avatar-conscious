/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as SettingsSettings from "./settings.pb"
export type SettingsGetRequest = {
  name?: string
}

export type SettingsGetResponse = {
  entry?: SettingsSettings.Config
}

export type SettingsListRequest = {
  name?: string
}

export type SettingsListResponse = {
  groups?: {[key: string]: SettingsSettings.ConfigList}
}

export type SettingsSetRequest = {
  name?: string
  entry?: SettingsSettings.Config
}

export type SettingsSetResponse = {
  entry?: SettingsSettings.Config
}

export class Settings {
  static Get(req: SettingsGetRequest, initReq?: fm.InitReq): Promise<SettingsGetResponse> {
    return fm.fetchReq<SettingsGetRequest, SettingsGetResponse>(`/v1/Settings/get?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static List(req: SettingsListRequest, initReq?: fm.InitReq): Promise<SettingsListResponse> {
    return fm.fetchReq<SettingsListRequest, SettingsListResponse>(`/v1/Settings/list?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static Set(req: SettingsSetRequest, initReq?: fm.InitReq): Promise<SettingsSetResponse> {
    return fm.fetchReq<SettingsSetRequest, SettingsSetResponse>(`/v1/Settings/set`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}