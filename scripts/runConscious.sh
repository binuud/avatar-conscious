
# export KEY=""
# export LLM_HOST="http://host.docker.internal:11434"
# export MQTT_HOST="tcp://host.docker.internal:1883"

export LLM_HOST="http://localhost:11234"
export MQTT_HOST="tcp://localhost:1883"

export LLM_MODEL="gemma4:latest"
export FILE_STORE_PATH=$HOME/Code/avatar/conscious/dataVolume
export SERVABLE_DIRECTORY=$HOME/Code/avatar/conscious/dataVolume/gen

reflex $(cat .reflex) -- go run cmd/conscious/main.go