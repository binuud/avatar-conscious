/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

export enum ConfigType {
  DATA_IGNORE = "DATA_IGNORE",
  STRING = "STRING",
  NUMBER = "NUMBER",
  FLOAT = "FLOAT",
  BOOLEAN = "BOOLEAN",
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

export type Configs = {
  entry?: {[key: string]: Config}
}