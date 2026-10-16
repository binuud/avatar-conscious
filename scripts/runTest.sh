export LLM_HOST="http://127.0.0.1:11234"
export MQTT_HOST="tcp://localhost:1883"

# export LLM_HOST="http://host.docker.internal:11234"
# export MQTT_HOST="tcp://host.docker.internal:1883"

# export LLM_MODEL="gemma4"

# export FILE_STORE_PATH=${HOME}/Code/ai/conscious/dataVolume/
# export SERVABLE_DIRECTORY=${HOME}/Code/ai/aiOutData/dataVolume/generatedAiFiles/

export LLM_MODEL="mlx-community/gemma-4-12b-it-4bit"
# export LLM_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"
export IMAGE_QUERY_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"

export FILE_STORE_PATH=${PWD}/dataVolume/
export SERVABLE_DIRECTORY=${PWD}/dataVolume/gen


# go test ./pkg/... -v

# go test ./pkg/... -v -run "Test_RouterViaTool"
go test -v --run Test_LLM_TOOL ./pkg/servers/consciousServer
# go test ./pkg/... -v -run "TestConsciousServer_LLMImageProcess"