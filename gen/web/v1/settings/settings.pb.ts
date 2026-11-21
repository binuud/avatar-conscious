/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

export enum ConfigType {
  DATA_IGNORE = "DATA_IGNORE",
  STRING = "STRING",
  MULTI_LINE_STRING = "MULTI_LINE_STRING",
  NUMBER = "NUMBER",
  FLOAT = "FLOAT",
  BOOLEAN = "BOOLEAN",
  MULTI_OPTION = "MULTI_OPTION",
  DATALIST = "DATALIST",
}

export type Config = {
  name?: string
  description?: string
  group?: string
  value?: string
  default?: string
  envValue?: string
  options?: string[]
  weight?: string
  type?: ConfigType
  numLines?: string
  min?: string
  max?: string
}

export type ConfigList = {
  configs?: Config[]
}

export type Configs = {
  name?: string
  entry?: {[key: string]: Config}
}