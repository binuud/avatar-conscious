# export LLM_HOST="http://127.0.0.1:12434"
# export LLM_MODEL="gemma4"
# export MQTT_HOST="tcp://localhost:1883"
# export FILE_STORE_PATH=${HOME}/Code/ai/conscious/dataVolume/
# export SERVABLE_DIRECTORY=${HOME}/Code/ai/aiOutData/dataVolume/generatedAiFiles/

export LLM_HOST="http://host.docker.internal:11434"
export LLM_MODEL="gemma4:latest"
export IMAGE_QUERY_MODEL="gemma4:latest"
export MQTT_HOST="tcp://host.docker.internal:1883"
export FILE_STORE_PATH=/app/dataVolume/
export SERVABLE_DIRECTORY=${HOME}/Code/ai/aiOutData/dataVolume/generatedAiFiles/

go test ./pkg/... -v

# go test ./pkg/... -v -run "TestConsciousServer_LLMImageProcess"