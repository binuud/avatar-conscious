
# export KEY=""
# export LLM_HOST="http://host.docker.internal:11234"
# export MQTT_HOST="tcp://host.docker.internal:1883"

export LLM_HOST="http://localhost:11234"
export MQTT_HOST="tcp://localhost:1883"

export LLM_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"
export IMAGE_MODEL="Runpod/FLUX.2-klein-4B-mflux-4bit"
export VISION_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"

export FILE_STORE_PATH=$PWD/dataVolume
export SERVABLE_DIRECTORY=$PWD/dataVolume/static

reflex $(cat .reflex) -- go run cmd/conscious/main.go
# go run cmd/conscious/main.go