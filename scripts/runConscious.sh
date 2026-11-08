
# export KEY=""
# export LLM_ENDPOINT="http://host.docker.internal:11234"
# export LLM_ENDPOINT="tcp://host.docker.internal:1883"

export LLM_ENDPOINT="http://localhost:11234"
export MQTT_ENDPOINT="tcp://localhost:1883"

# export LLM_MODEL="mlx-community/gemma-3-12b-it-4bit"
# export LLM_MODEL="mlx-community/gemma-4-12b-it-4bit"
export LLM_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"

export IMAGE_MODEL="Runpod/FLUX.2-klein-4B-mflux-4bit"
export VISION_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"

export FILE_STORE_PATH=$PWD/dataVolume
export SERVABLE_DIRECTORY=$PWD/dataVolume/static

reflex $(cat .reflex) -- go run cmd/conscious/main.go 
# go run cmd/conscious/main.go