/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

export enum Permission {
  PERMISSION_IGNORE = "PERMISSION_IGNORE",
  WEBCAMERA = "WEBCAMERA",
  ROBOT = "ROBOT",
  ROVER = "ROVER",
  ARM = "ARM",
  HOMEAUTOMATION = "HOMEAUTOMATION",
}

export type Avatar = {
  uuid?: string
  name?: string
  description?: string
  permissions?: Permission
  systemPrompt?: string
}