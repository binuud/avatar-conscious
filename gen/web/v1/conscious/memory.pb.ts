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

export type GetMemoryRequest = {
  uuid?: string
  name?: string
}

export type GetMemoryResponse = {
  meta?: MemoryMeta
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