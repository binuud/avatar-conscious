/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

export enum HealthResponseHealthStatus {
  HealthIgnore = "HealthIgnore",
  Active = "Active",
  Waiting = "Waiting",
  Down = "Down",
}

export type Timestamp = {
  seconds?: string
  nanos?: number
}

export type InRequest = {
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

export type InResponse = {
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
  createdAt?: Timestamp
  uuid?: string
  status?: HealthResponseHealthStatus
  statusMessage?: string
  errorMessage?: string
}