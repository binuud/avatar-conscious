/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as Common from "../conscious/common.pb"
import * as ContextMessage from "../conscious/contextMessage.pb"

export enum AvatarResult {
  RESULT_IGNORE = "RESULT_IGNORE",
  RESULT_SUCCESS = "RESULT_SUCCESS",
  RESULT_ERROR = "RESULT_ERROR",
}

export enum AvatarStatus {
  STATUS_IGNORE = "STATUS_IGNORE",
  STATUS_NOT_STARTED = "STATUS_NOT_STARTED",
  STATUS_INIT = "STATUS_INIT",
  STATUS_ACTIVE = "STATUS_ACTIVE",
  STATUS_STOPPED = "STATUS_STOPPED",
  STATUS_COMPLETED = "STATUS_COMPLETED",
  STATUS_WAITING = "STATUS_WAITING",
  STATUS_ERROR = "STATUS_ERROR",
}

export enum AvatarAction {
  ACTION_IGNORE = "ACTION_IGNORE",
  ACTION_INIT = "ACTION_INIT",
  ACTION_BEGIN = "ACTION_BEGIN",
  ACTION_DEBUG = "ACTION_DEBUG",
  ACTION_WAIT = "ACTION_WAIT",
  ACTION_RESUME = "ACTION_RESUME",
  ACTION_BEGIN_STEP = "ACTION_BEGIN_STEP",
  ACTION_STOP_STEP = "ACTION_STOP_STEP",
  ACTION_COMPLETE = "ACTION_COMPLETE",
  ACTION_ERROR = "ACTION_ERROR",
  ACTION_STOP = "ACTION_STOP",
}

export enum StepStatusResponseType {
  RESPONSE_TYPE_IGNORE = "RESPONSE_TYPE_IGNORE",
  IMAGEDATA = "IMAGEDATA",
  TEXT = "TEXT",
}

export type StepStatus = {
  status?: AvatarStatus
  responseType?: StepStatusResponseType
  reqeust?: string
  response?: string
  result?: AvatarResult
  createdAt?: Common.Timestamp
  completedAt?: Common.Timestamp
}

export type AvatarControl = {
  uuid?: string
  avatarUuid?: string
  name?: string
  status?: AvatarStatus
  action?: AvatarAction
  currentStep?: number
  stepStatuses?: StepStatus[]
  response?: ContextMessage.ContextMesssage
  result?: AvatarResult
  createdAt?: Common.Timestamp
  completedAt?: Common.Timestamp
}