/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as AvatarAvatar from "./avatar.pb"
export type AvatarStoreListRequest = {
  limit?: number
  offset?: number
}

export type AvatarStoreListResponse = {
  avatars?: AvatarAvatar.Avatar[]
  totalCount?: number
}

export type AvatarStoreGetRequest = {
  uuid?: string
}

export type AvatarStoreGetResponse = {
  avatar?: AvatarAvatar.Avatar
}

export type AvatarStoreCreateRequest = {
  name?: string
  description?: string
  permissions?: AvatarAvatar.Permission
  systemPrompt?: string
}

export type AvatarStoreCreateResponse = {
  avatar?: AvatarAvatar.Avatar
}

export type AvatarStoreUpdateRequest = {
  uuid?: string
  name?: string
  description?: string
  permissions?: AvatarAvatar.Permission
  systemPrompt?: string
}

export type AvatarStoreUpdateResponse = {
  avatar?: AvatarAvatar.Avatar
}

export type AvatarStoreDeleteRequest = {
  uuid?: string
}

export type AvatarStoreDeleteResponse = {
  message?: string
}

export class AvatarStoreService {
  static List(req: AvatarStoreListRequest, initReq?: fm.InitReq): Promise<AvatarStoreListResponse> {
    return fm.fetchReq<AvatarStoreListRequest, AvatarStoreListResponse>(`/v1/avatarStore/list?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static Get(req: AvatarStoreGetRequest, initReq?: fm.InitReq): Promise<AvatarStoreGetResponse> {
    return fm.fetchReq<AvatarStoreGetRequest, AvatarStoreGetResponse>(`/v1/avatarStore/get?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
  static Create(req: AvatarStoreCreateRequest, initReq?: fm.InitReq): Promise<AvatarStoreCreateResponse> {
    return fm.fetchReq<AvatarStoreCreateRequest, AvatarStoreCreateResponse>(`/v1/avatarStore/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Update(req: AvatarStoreUpdateRequest, initReq?: fm.InitReq): Promise<AvatarStoreUpdateResponse> {
    return fm.fetchReq<AvatarStoreUpdateRequest, AvatarStoreUpdateResponse>(`/v1/avatarStore/update`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static Delete(req: AvatarStoreDeleteRequest, initReq?: fm.InitReq): Promise<AvatarStoreDeleteResponse> {
    return fm.fetchReq<AvatarStoreDeleteRequest, AvatarStoreDeleteResponse>(`/v1/avatarStore/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}