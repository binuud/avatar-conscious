/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as Common from "./common.pb"

export enum ContextMesssageContextType {
  CONTEXT_TYPE_IGNORE = "CONTEXT_TYPE_IGNORE",
  LLM_STREAM_QUERY = "LLM_STREAM_QUERY",
  LLM_QUERY = "LLM_QUERY",
  LLM_IMAGE_DESCRIBE = "LLM_IMAGE_DESCRIBE",
  IMAGE_GENERATION = "IMAGE_GENERATION",
  TTS = "TTS",
}

export enum ContextMesssageContextRole {
  ROLE_IGNORE = "ROLE_IGNORE",
  USER = "USER",
  ASSISTANT = "ASSISTANT",
  SYSTEM = "SYSTEM",
  AGENT = "AGENT",
}

export enum ContextMesssageContextStatus {
  STATUS_IGNORE = "STATUS_IGNORE",
  SUCCESS = "SUCCESS",
  FAIL = "FAIL",
  ONGOING = "ONGOING",
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

export type ContextMesssage = {
  uuid?: string
  name?: string
  description?: string
  model?: string
  type?: ContextMesssageContextType
  role?: ContextMesssageContextRole
  stream?: boolean
  prompt?: string
  systemPrompt?: string
  seed?: string
  temperature?: number
  data?: string
  bodyText?: string
  schema?: string
  bodyBytes?: Uint8Array
  vision?: Uint8Array
  audio?: Uint8Array
  text?: Uint8Array
  status?: ContextMesssageContextStatus
  statusMessage?: string
  errorMessage?: string
  response?: string
  responseImage?: string
  reasoning?: string
  thinking?: string
  createdAt?: Common.Timestamp
  completedAt?: Common.Timestamp
}