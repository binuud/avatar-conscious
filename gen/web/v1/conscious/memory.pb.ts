/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as CommonCommon from "../common/common.pb"
import * as ConsciousContextMessage from "./contextMessage.pb"
export type SingleContext = {
  in?: ConsciousContextMessage.ContextMesssage
  out?: ConsciousContextMessage.ContextMesssage
}

export type MemoryMeta = {
  uuid?: string
  name?: string
  description?: string
  model?: string
  systemPrompt?: string
  seed?: string
  temperature?: number
  createdAt?: CommonCommon.Timestamp
}

export type Memory = {
  meta?: MemoryMeta
  contexts?: SingleContext[]
}

export type ListMemoryRequest = {
  name?: string
}

export type ListMemoryResponse = {
  memories?: MemoryMeta[]
}

export type GetMemoryMetaRequest = {
  uuid?: string
  name?: string
}

export type GetMemoryMetaResponse = {
  meta?: MemoryMeta
}

export type GetMemoryRequest = {
  uuid?: string
  name?: string
}

export type GetMemoryResponse = {
  memory?: Memory
}

export type AddMemoryRequest = {
  meta?: MemoryMeta
}

export type AddMemoryResponse = {
  meta?: MemoryMeta
}

export type UpdateMemoryRequest = {
  meta?: MemoryMeta
}

export type UpdateMemoryResponse = {
  meta?: MemoryMeta
}

export type DeleteMemoryRequest = {
  uuid?: string
  name?: string
}

export type DeleteMemoryResponse = {
  success?: boolean
}

export type ClearMemoryRequest = {
  uuid?: string
  name?: string
}

export type ClearMemoryResponse = {
  success?: boolean
}