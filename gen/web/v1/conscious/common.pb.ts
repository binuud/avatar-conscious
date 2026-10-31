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

export type Model = {
  uuid?: string
  id?: string
  created?: string
  ownedBy?: string
}

export type HealthRequest = {
  uuid?: string
}

export type HealthResponse = {
  uuid?: string
  status?: HealthResponseHealthStatus
  statusMessage?: string
  errorMessage?: string
  createdAt?: string
}

export type ModelRequest = {
}

export type ModelResponse = {
  models?: Model[]
}

export type ClearMemoryRequest = {
}

export type ClearMemoryResponse = {
}