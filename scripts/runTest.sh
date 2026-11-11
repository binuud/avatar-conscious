# export LLM_ENDPOINT="http://localhost:11234/v1"
# export MQTT_ENDPOINT="tcp://localhost:1883"

# export LLM_MODEL="gemma4"

# export FILE_STORE_PATH=${HOME}/Code/ai/conscious/dataVolume/
# export SERVABLE_DIRECTORY=${HOME}/Code/ai/aiOutData/dataVolume/generatedAiFiles/

export LLM_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"
export IMAGE_QUERY_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"

export FILE_STORE_PATH=$PWD/dataVolume
export SERVABLE_DIRECTORY=$PWD/dataVolume/static

# this entry is for testing settings package
export SETTINGS_IGNORE_TESTCASE_ENDPOINT="Sample123"

# go test ./pkg/...
# go test ./pkg/...  -v

# go test -v --run TestConsciousServer_LLMChat ./pkg/servers/consciousServer

## test run for running avatar step by step
# go test -v  --run Test_Avatar_By_Step_Election_Officer ./pkg/servers/avatarServer

## test run for running avatar step by step, control via bidirectional server
#go test -v  --run Test_AvatarServer_Control ./pkg/servers/avatarServer
# go test -v --run  Test_UiLaunchTool ./pkg/tools/testTools/
# go test -v --run  Test_ImageEdit ./pkg/tools/testTools/
# go test -v  ./pkg/servers/avatarStoreServer
go test -v  ./pkg/servers/mcpManagerServer
# go test -v  ./pkg/settings/
# go test ./pkg/... -v -run "Test_RouterViaTool"
# go test -v --run Test_LLM_TOOL ./pkg/servers/consciousServer
# go test ./pkg/... -v -run "TestConsciousServer_LLMImageProcess"