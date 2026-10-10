
# export KEY=""
export LLM_HOST="http://host.docker.internal:11434"
export LLM_MODEL="gemma4:latest"
export MQTT_HOST="tcp://host.docker.internal:1883"
export FILE_STORE_PATH=/app/dataVolume
export SERVABLE_DIRECTORY=/app/dataVolume/gen/

reflex $(cat .reflex) -- go run cmd/conscious/main.go