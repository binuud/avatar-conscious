/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

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

export enum StepStepType {
  STEP_TYPE_ENUM_DEFAULT = "STEP_TYPE_ENUM_DEFAULT",
  STEP_TYPE_ENUM_AI = "STEP_TYPE_ENUM_AI",
  STEP_TYPE_ENUM_UI = "STEP_TYPE_ENUM_UI",
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
  aiAllowedTools?: string[]
}

export type UiCollection = {
  ui?: QuatumUIQuantumUI.QuantumUI[]
}

export type Step = {
  name?: string
  description?: string
  prompt?: string
  stepType?: StepStepType
}

export type Steps = {
  steps?: Step[]
}

export type Avatar = {
  meta?: AvatarMeta
  permissions?: Permission
  systemPrompt?: string
  steps?: Step[]
  uiCollection?: UiCollection
  content?: string
}