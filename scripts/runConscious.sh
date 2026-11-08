
# export KEY=""
# export LLM_ENDPOINT="http://host.docker.internal:11234"
# export LLM_ENDPOINT="tcp://host.docker.internal:1883"

export MQTT_ENDPOINT="tcp://localhost:1883"

# for open router
# export LLM_ENDPOINT="https://openrouter.ai/api/v1"
# export LLM_MODEL="qwen/qwen3.5-9b-20260310"
# export IMAGE_MODEL="black-forest-labs/flux-3-image-klein"
# export VISION_MODEL="qwen/qwen3.5-9b-20260310"


# export LLM_MODEL="mlx-community/gemma-3-12b-it-4bit"
# export LLM_MODEL="mlx-community/gemma-4-12b-it-4bit"


export LLM_ENDPOINT="http://localhost:11234/v1"
export LLM_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"
export IMAGE_MODEL="Runpod/FLUX.2-klein-4B-mflux-4bit"
export VISION_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"

export FILE_STORE_PATH=$PWD/dataVolume
export SERVABLE_DIRECTORY=$PWD/dataVolume/static

## get the LLM API KEY
source .env

reflex $(cat .reflex) -- go run cmd/conscious/main.go 
# go run cmd/conscious/main.go