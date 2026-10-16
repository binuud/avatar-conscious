/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as Common from "./common.pb"

export enum contextMesssageContextType {
  CONTEXT_TYPE_IGNORE = "CONTEXT_TYPE_IGNORE",
  LLM_STREAM_QUERY = "LLM_STREAM_QUERY",
  LLM_QUERY = "LLM_QUERY",
  LLM_IMAGE_DESCRIBE = "LLM_IMAGE_DESCRIBE",
  IMAGE_GENERATION = "IMAGE_GENERATION",
  TTS = "TTS",
}

export enum contextMesssageContextRole {
  ROLE_IGNORE = "ROLE_IGNORE",
  USER = "USER",
  ASSISTANT = "ASSISTANT",
  SYSTEM = "SYSTEM",
  AGENT = "AGENT",
}

export enum contextMesssageContextStatus {
  STATUS_IGNORE = "STATUS_IGNORE",
  SUCCESS = "SUCCESS",
  FAIL = "FAIL",
  FATAL = "FATAL",
}

export type UploadRequest = {
  fileName?: string
  fileData?: string
  contentType?: string
}

export type UploadResponse = {
  uploadId?: string
  status?: string
  errorMessage?: string
}

export type contextMesssage = {
  uuid?: string
  name?: string
  description?: string
  model?: string
  type?: contextMesssageContextType
  role?: contextMesssageContextRole
  stream?: boolean
  prompt?: string
  systemPrompt?: string
  data?: string
  bodyText?: string
  bodyBytes?: Uint8Array
  vision?: Uint8Array
  audio?: Uint8Array
  text?: Uint8Array
  status?: contextMesssageContextStatus
  statusMessage?: string
  errorMessage?: string
  response?: string
  responseImage?: string
  reasoning?: string
  thinking?: string
  createdAt?: Common.Timestamp
}