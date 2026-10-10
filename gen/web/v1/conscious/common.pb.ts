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

export type HealthRequest = {
  uuid?: string
}

export type HealthResponse = {
  uuid?: string
  status?: HealthResponseHealthStatus
  statusMessage?: string
  errorMessage?: string
  createdAt?: Timestamp
}