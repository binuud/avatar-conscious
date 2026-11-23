/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as GoogleProtobufField_mask from "../../google/protobuf/field_mask.pb"
import * as GoogleProtobufTimestamp from "../../google/protobuf/timestamp.pb"

export enum ApiProtocol {
  API_PROTOCOL_UNSPECIFIED = "API_PROTOCOL_UNSPECIFIED",
  API_PROTOCOL_OPENAI_COMPATIBLE = "API_PROTOCOL_OPENAI_COMPATIBLE",
  API_PROTOCOL_ANTHROPIC = "API_PROTOCOL_ANTHROPIC",
  API_PROTOCOL_GOOGLE_GEMINI = "API_PROTOCOL_GOOGLE_GEMINI",
  API_PROTOCOL_AWS_BEDROCK = "API_PROTOCOL_AWS_BEDROCK",
  API_PROTOCOL_CUSTOM = "API_PROTOCOL_CUSTOM",
}

export enum CredentialStatus {
  CREDENTIAL_STATUS_UNSPECIFIED = "CREDENTIAL_STATUS_UNSPECIFIED",
  CREDENTIAL_STATUS_ACTIVE = "CREDENTIAL_STATUS_ACTIVE",
  CREDENTIAL_STATUS_EXPIRED = "CREDENTIAL_STATUS_EXPIRED",
  CREDENTIAL_STATUS_REVOKED = "CREDENTIAL_STATUS_REVOKED",
  CREDENTIAL_STATUS_RATE_LIMITED = "CREDENTIAL_STATUS_RATE_LIMITED",
}

export type LlmProvider = {
  uuid?: string
  displayName?: string
  baseEndpoint?: string
  apiKeyGenerationUrl?: string
  protocol?: ApiProtocol
  authHeaderName?: string
  authHeaderPrefix?: string
  requiresApiKey?: boolean
  credentials?: ApiCredential[]
  customHeaders?: {[key: string]: string}
  createdAt?: GoogleProtobufTimestamp.Timestamp
  updatedAt?: GoogleProtobufTimestamp.Timestamp
}

export type ApiCredential = {
  id?: string
  label?: string
  apiKeyOrSecretRef?: string
  maskedKey?: string
  envVarName?: string
  status?: CredentialStatus
  createdAt?: GoogleProtobufTimestamp.Timestamp
  expiresAt?: GoogleProtobufTimestamp.Timestamp
}

export type ProviderRegistrySnapshot = {
  providers?: LlmProvider[]
  exportedAt?: GoogleProtobufTimestamp.Timestamp
}

export type CreateProviderRequest = {
  provider?: LlmProvider
}

export type GetProviderRequest = {
  id?: string
}

export type ListProvidersRequest = {
  pageSize?: number
  pageToken?: string
}

export type ListProvidersResponse = {
  providers?: LlmProvider[]
  nextPageToken?: string
}

export type UpdateProviderRequest = {
  provider?: LlmProvider
  updateMask?: GoogleProtobufField_mask.FieldMask
}

export type DeleteProviderRequest = {
  id?: string
}

export type DeleteProviderResponse = {
  deleted?: boolean
}

export class ProviderRegistryService {
  static CreateProvider(req: CreateProviderRequest, initReq?: fm.InitReq): Promise<LlmProvider> {
    return fm.fetchReq<CreateProviderRequest, LlmProvider>(`/llm.ProviderRegistryService/CreateProvider`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static GetProvider(req: GetProviderRequest, initReq?: fm.InitReq): Promise<LlmProvider> {
    return fm.fetchReq<GetProviderRequest, LlmProvider>(`/llm.ProviderRegistryService/GetProvider`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static ListProviders(req: ListProvidersRequest, initReq?: fm.InitReq): Promise<ListProvidersResponse> {
    return fm.fetchReq<ListProvidersRequest, ListProvidersResponse>(`/llm.ProviderRegistryService/ListProviders`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static UpdateProvider(req: UpdateProviderRequest, initReq?: fm.InitReq): Promise<LlmProvider> {
    return fm.fetchReq<UpdateProviderRequest, LlmProvider>(`/llm.ProviderRegistryService/UpdateProvider`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
  static DeleteProvider(req: DeleteProviderRequest, initReq?: fm.InitReq): Promise<DeleteProviderResponse> {
    return fm.fetchReq<DeleteProviderRequest, DeleteProviderResponse>(`/llm.ProviderRegistryService/DeleteProvider`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}