/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as ContextMessage from "../conscious/contextMessage.pb"

export enum Permission {
  PERMISSION_IGNORE = "PERMISSION_IGNORE",
  WEBCAMERA = "WEBCAMERA",
  ROBOT = "ROBOT",
  ROVER = "ROVER",
  ARM = "ARM",
  HOMEAUTOMATION = "HOMEAUTOMATION",
}

export enum AvatarActionType {
  AvatarActionTypeIgnore = "AvatarActionTypeIgnore",
  CommandLine = "CommandLine",
  WebSearch = "WebSearch",
  WebApi = "WebApi",
  Robotic = "Robotic",
  IO = "IO",
  AiQuery = "AiQuery",
  AiImageDescribe = "AiImageDescribe",
  AiImageGeneration = "AiImageGeneration",
}

export enum AvatarActionStatus {
  AvatarActionStatusIgnore = "AvatarActionStatusIgnore",
  Success = "Success",
  Active = "Active",
  Rerun = "Rerun",
  Failed = "Failed",
  Fatal = "Fatal",
}

export type ToolParams = {
  name?: string
  type?: string
  description?: string
  required?: boolean
}

export type Tool = {
  name?: string
  description?: string
  type?: string
  parameters?: {[key: string]: ToolParams}
}

export type AvatarAction = {
  name?: string
  description?: string
  action?: string
  disabled?: boolean
  type?: AvatarActionType
  Status?: AvatarActionStatus
  completedLogic?: string
  successLogic?: string
  failLogic?: string
  model?: string
  prompt?: string
  systemPrompt?: string
  In?: ContextMessage.ContextMesssage
  Out?: ContextMessage.ContextMesssage
}

export type AvatarLoop = {
  name?: string
  description?: string
  model?: string
  prompt?: string
  systemPrompt?: string
  runActions?: AvatarAction[]
  loopAction?: AvatarAction
}

export type AvatarMeta = {
  uuid?: string
  name?: string
  description?: string
  triggers?: string[]
}

export type Avatar = {
  meta?: AvatarMeta
  permissions?: Permission
  systemPrompt?: string
  loops?: AvatarLoop[]
  data?: string
}