/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

export enum UiComponentUiType {
  UiTypeIgnore = "UiTypeIgnore",
  TEXTAREA = "TEXTAREA",
  INPUT = "INPUT",
  CHECKBOX = "CHECKBOX",
  PLAINTEXT = "PLAINTEXT",
  BUTTON = "BUTTON",
  WEBCAM = "WEBCAM",
}

export type UiComponent = {
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
  type?: UiComponentUiType
}

export type QuantumUI = {
  uuid?: string
  name?: string
  description?: string
  title?: string
  components?: UiComponent[]
}