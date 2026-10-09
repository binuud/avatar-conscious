# Maintainer Binu Udayakumar <binu@dronasys.com>

REPO=dronasys


CONSCIOUS_APP_NAME= conscious

DEBUG_IMAGE= debug-golang

BUILD_VER= a1.0.0
DOCKER_HUB_TAG= a1.0.0

GOLANG_IMAGE=golang:1.27.1

TAGGED_NAME = $(REPO)/$(APP_NAME)

MULTI_PLATFORM_DOCKER = --platform=linux/amd64,linux/arm64

# Optional environmental variables

BASE_DIR	:=	$(shell git rev-parse --show-toplevel)
BIN			:=	$(BASE_DIR)/bin

## Needs protoc to be installed

ifndef PROTOC
PROTOC = protoc
endif

ifndef PROTOWEB
PROTOWEB = protoc-gen-grpc-web
endif

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help vendor vendor-update
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

##
## Build, Vendor, Publish
##

main: ## Build main
	go build $(BUILDFLAGS) -o main .

build: ## build backend and frontend as a single image
	docker build --no-cache  --build-arg BUILD_VER=$(BUILD_VER) $(MULTI_PLATFORM_DOCKER) -t $(REPO)/$(CONSCIOUS_APP_NAME) -f deployment/docker/conscious.dockerfile .

build-go-binary: ## build the go binary
	go build -o brainBackend cmd/conscious/main.go

publish-docker: ## Publish docker backend and combined images to docker !! PUBLIC IMAGE !!
	@echo "====================> Pushing image to dockerhub !! PUBLIC IMAGE !!"
## @docker login --username=$(DOCKER_HUB_USER) --password="$(DOCKER_HUB_PASSWD)"
	docker tag $(REPO)/${CONSCIOUS_APP_NAME} $(REPO)/${CONSCIOUS_APP_NAME}:$(DOCKER_HUB_TAG)
	docker push $(REPO)/${CONSCIOUS_APP_NAME}:$(DOCKER_HUB_TAG)

publish-private-microservices: ## Deploy all images to docker hub 
	@echo "====================> Pushing image to dockerhub$(AWS_ECR) ."
	## @docker login --username=$(DOCKER_HUB_USER) --password="$(DOCKER_HUB_PASSWD)"
	aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin $(AWS_ECR)

	docker tag $(REPO)/${CONSCIOUS_APP_NAME} $(AWS_ECR)/${CONSCIOUS_APP_NAME}:$(DOCKER_HUB_TAG)
	docker push $(AWS_ECR)/$(CONSCIOUS_APP_NAME):$(DOCKER_HUB_TAG)

tidy: ## run tidy
	go mod tidy

vendor: ## run vendor
	go mod tidy
	go mod vendor

staticcheck:
	GOFLAGS="" go install honnef.co/go/tools/cmd/staticcheck@v0.2.2
	staticcheck $(PKGS)

errcheck:
	go get -u github.com/kisielk/errcheck
	errcheck -ignoregenerated -verbose -blank $(PKGS)


### the ordering of sub targets is very important
### do not change
protos: protos-conscious ## run all protos


## protos
protos-conscious: ## Build Go and Web Proto
# ****** Conscious AiModule 
	$(PROTOC) -I=./proto/.  \
	--go_out=./gen/go/ \
	--go_opt paths=source_relative \
	--go-grpc_out=./gen/go/ \
	--go-grpc_opt paths=source_relative \
	--grpc-gateway_out=./gen/ts/ \
	--grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	--grpc-gateway-ts_out=./gen/ts/ \
	--grpc-gateway-ts_opt paths=source_relative \
	--grpc-gateway-ts_opt generate_unbound_methods=true \
	--openapiv2_out=./gen/openapi/v1/brain/ \
	proto/v1/brain/brain.proto proto/v1/brain/consciousService.proto \
	proto//google/protobuf/timestamp.proto
	yq eval ./gen/go/v1/brain/openapi.yaml -o=json -P > ./gen/go/v1/brain/openapi.json
# 	--oas_out ./gen/go/v1/brain/ \
##
## Test Code
##

test: ## Run test 
	./scripts/runTest.sh



run: ## combine all services for single user case
	./scripts/runConscious.sh
