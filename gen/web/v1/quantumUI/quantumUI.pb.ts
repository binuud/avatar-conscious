/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

export enum uiComponentUiType {
  UiTypeIgnore = "UiTypeIgnore",
  TEXTAREA = "TEXTAREA",
  INPUT = "INPUT",
  CHECKBOX = "CHECKBOX",
  PLAINTEXT = "PLAINTEXT",
  BUTTON = "BUTTON",
  WEBCAM = "WEBCAM",
}

export enum ConfigType {
  DATA_IGNORE = "DATA_IGNORE",
  STRING = "STRING",
  NUMBER = "NUMBER",
  FLOAT = "FLOAT",
  BOOLEAN = "BOOLEAN",
}

export type uiComponent = {
  id?: string
  class?: string
  label?: string
  value?: string
  placeHolder?: string
  help?: string
  style?: string
  width?: number
  rows?: number
  cols?: number
  type?: uiComponentUiType
}

export type quantumUI = {
  uuid?: string
  name?: string
  description?: string
  title?: string
  components?: uiComponent[]
}

export type Config = {
  name?: string
  description?: string
  value?: string
  group?: string
  type?: ConfigType
}

export type ConfigList = {
  configs?: Config[]
}

export type Settings = {
  entry?: {[key: string]: Config}
}