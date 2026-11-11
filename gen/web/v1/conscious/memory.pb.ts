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

export type Memory = {
  uuid?: string
  name?: string
  description?: string
  model?: string
  systemPrompt?: string
  seed?: string
  temperature?: number
  contexts?: SingleContext[]
  createdAt?: CommonCommon.Timestamp
}