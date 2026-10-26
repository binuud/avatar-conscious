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

export enum StepStatusAvatarStepStatus {
  AvatarActionStatusIgnore = "AvatarActionStatusIgnore",
  Success = "Success",
  Active = "Active",
  Rerun = "Rerun",
  Failed = "Failed",
  Fatal = "Fatal",
  NotStarted = "NotStarted",
}

export enum AvatarControlAvatarStatus {
  AvatarActionStatusIgnore = "AvatarActionStatusIgnore",
  Success = "Success",
  Active = "Active",
  Rerun = "Rerun",
  Failed = "Failed",
  Fatal = "Fatal",
  Stopped = "Stopped",
}

export enum AvatarControlAvatarAction {
  AvatarActionIgnore = "AvatarActionIgnore",
  Run = "Run",
  StartDebug = "StartDebug",
  Stop = "Stop",
  RunStep = "RunStep",
  StopStep = "StopStep",
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
  status?: StepStatusAvatarStepStatus
  createdAt?: Common.Timestamp
  completedAt?: Common.Timestamp
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
  status?: AvatarControlAvatarStatus
  action?: AvatarControlAvatarAction
  currentStep?: number
  stepStatuses?: StepStatus[]
  createdAt?: Common.Timestamp
  completedAt?: Common.Timestamp
}