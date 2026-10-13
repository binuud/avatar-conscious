/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
export type SystemMetrics = {
  uuid?: string
  cpuUsage?: number
  totalRamBytes?: string
  usedRamBytes?: string
  ramUsage?: number
  gpuUsage?: number
  totalDiskBytes?: string
  usedDiskBytes?: string
  createdAt?: string
}

export type GetSystemMetricsRequest = {
  intervalSeconds?: number
}

export class Metrics {
  static GetSystemMetrics(req: GetSystemMetricsRequest, entityNotifier?: fm.NotifyStreamEntityArrival<SystemMetrics>, initReq?: fm.InitReq): Promise<void> {
    return fm.fetchStreamingRequest<GetSystemMetricsRequest, SystemMetrics>(`/v1/consciousMetrics/getSystemMetrics`, entityNotifier, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)})
  }
}