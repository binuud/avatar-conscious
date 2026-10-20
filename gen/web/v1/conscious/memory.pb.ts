/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as Common from "./common.pb"
import * as ContextMessage from "./contextMessage.pb"
export type SingleContext = {
  in?: ContextMessage.ContextMesssage
  out?: ContextMessage.ContextMesssage
}

export type Memory = {
  uuid?: string
  name?: string
  description?: string
  model?: string
  systemPrompt?: string
  seed?: string
  temperature?: number
  contexts?: SingleContext[]
  createdAt?: Common.Timestamp
}