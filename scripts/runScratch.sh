
# export KEY=""
# export LLM_HOST="http://host.docker.internal:11434"
# export LLM_MODEL="gemma4:latest"
# export MQTT_HOST="tcp://host.docker.internal:1883"
export FILE_STORE_PATH=${HOME}/Code/avatar/avatar-avatarconscious/dataVolume/
export SERVABLE_DIRECTORY=${HOME}/Code/avatar/avatar-avatarconscious/dataVolume/gen/



export LLM_HOST="http://localhost:11234"
export MQTT_HOST="tcp://localhost:1883"

export LLM_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"
export IMAGE_MODEL="Runpod/FLUX.2-klein-4B-mflux-4bit"
export VISION_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"

export FILE_STORE_PATH=$PWD/dataVolume
export SERVABLE_DIRECTORY=$PWD/dataVolume/gen

go run cmd/scratch/toolCalling/main.go

#reflex go run cmd/scratch/skill/skillContextRead/main.go
# reflex $(cat .reflex) -- go run cmd/scratch/skill/skillContextRead/main.go