/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as GoogleProtobufTimestamp from "../../google/protobuf/timestamp.pb"

export enum HealthResponseHealthStatus {
  HealthIgnore = "HealthIgnore",
  Active = "Active",
  Waiting = "Waiting",
  Down = "Down",
}

export type Request = {
  uuid?: string
  name?: string
  description?: string
  prompt?: string
  bodyText?: string
  bodyBytes?: Uint8Array
  vision?: Uint8Array
  audio?: Uint8Array
  text?: Uint8Array
  data?: string
}

export type Response = {
  uuid?: string
  name?: string
  description?: string
  response?: string
  prompt?: string
  bodyText?: string
  bodyBytes?: Uint8Array
  vision?: Uint8Array
  audio?: Uint8Array
  text?: Uint8Array
  data?: string
}

export type HealthRequest = {
  uuid?: string
}

export type HealthResponse = {
  createdAt?: GoogleProtobufTimestamp.Timestamp
  uuid?: string
  status?: HealthResponseHealthStatus
  statusMessage?: string
  errorMessage?: string
}