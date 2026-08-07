/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/
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