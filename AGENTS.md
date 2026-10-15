# Project Overview
This repository is a GoLang repository. This is a AI Harness, that will mimic a conscious.

## Build & Test Commands
- **Generate proto buffers:** `make protos`

## Code Conventions
- Keep proto files under `proto/v1/conscious/`.

## Restricted Files
- Do not modify files in `dataVolume` without explicit user sign-off.

## Directory Structure

* cmd - the main go file
* documentation - all documentation md files
* pkg/servers - all grpc server implementations go here
* proto/v1 - all protofiles go here
* scripts - support scripts, like launchers etc
* vendor - go vendor files, do not make changes here