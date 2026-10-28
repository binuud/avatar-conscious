/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as Common from "../conscious/common.pb"
import * as QuatumUIQuantumUI from "../quantumUI/quantumUI.pb"

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

export enum AvatarStatus {
  AvatarActionStatusIgnore = "AvatarActionStatusIgnore",
  STATUS_INIT = "STATUS_INIT",
  STATUS_SUCCESS = "STATUS_SUCCESS",
  STATUS_ACTIVE = "STATUS_ACTIVE",
  STATUS_RERUN = "STATUS_RERUN",
  STATUS_ERROR = "STATUS_ERROR",
  STATUS_FAILED = "STATUS_FAILED",
  STATUS_FATAL = "STATUS_FATAL",
  STATUS_STOPPED = "STATUS_STOPPED",
  STATUS_NOT_STARTED = "STATUS_NOT_STARTED",
  STATUS_WAITING = "STATUS_WAITING",
  STATUS_UI_WAITING = "STATUS_UI_WAITING",
  STATUS_IO_WAITING = "STATUS_IO_WAITING",
}

export enum StepStatusResponseType {
  RESPONSE_TYPE_IGNORE = "RESPONSE_TYPE_IGNORE",
  IMAGEDATA = "IMAGEDATA",
  TEXT = "TEXT",
}

export enum AvatarControlAvatarAction {
  AVATAR_ACTION_IGNORE = "AVATAR_ACTION_IGNORE",
  INIT = "INIT",
  RUN = "RUN",
  START_DEBUG = "START_DEBUG",
  STOP = "STOP",
  RUN_STEP = "RUN_STEP",
  STOP_STEP = "STOP_STEP",
  SEND_DATA = "SEND_DATA",
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

export type AvatarMeta = {
  uuid?: string
  name?: string
  description?: string
  image?: string
  triggers?: string[]
  ui?: QuatumUIQuantumUI.QuantumUI
}

export type Step = {
  name?: string
  description?: string
  prompt?: string
}

export type Steps = {
  steps?: Step[]
}

export type StepStatus = {
  status?: AvatarStatus
  createdAt?: Common.Timestamp
  completedAt?: Common.Timestamp
  responseType?: StepStatusResponseType
  reqeust?: string
  response?: string
}

export type Avatar = {
  meta?: AvatarMeta
  permissions?: Permission
  systemPrompt?: string
  steps?: Step[]
  content?: string
}

export type AvatarControl = {
  meta?: AvatarMeta
  status?: AvatarStatus
  action?: AvatarControlAvatarAction
  currentStep?: number
  stepStatuses?: StepStatus[]
  response?: string
  createdAt?: Common.Timestamp
  completedAt?: Common.Timestamp
}