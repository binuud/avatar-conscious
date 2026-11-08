
# export KEY=""
# export LLM_ENDPOINT="http://host.docker.internal:11434"
# export LLM_MODEL="gemma4:latest"
# export MQTT_ENDPOINT="tcp://host.docker.internal:1883"
export FILE_STORE_PATH=${HOME}/Code/avatar/avatar-avatarconscious/dataVolume/
export SERVABLE_DIRECTORY=${HOME}/Code/avatar/avatar-avatarconscious/dataVolume/gen/

export LLM_ENDPOINT="http://localhost:11234/v1"
export MQTT_ENDPOINT="tcp://localhost:1883"

export LLM_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"
export IMAGE_MODEL="Runpod/FLUX.2-klein-4B-mflux-4bit"
export VISION_MODEL="lmstudio-community/Qwen3.5-9B-MLX-4bit"

export FILE_STORE_PATH=$PWD/dataVolume
export SERVABLE_DIRECTORY=$PWD/dataVolume/static


## run avatar election officer
# go run cmd/scratch/avatar/avatarTest/main.go

## read an avatar file, and have llm generate instructions from it
# go run cmd/scratch/avatar/avatarReadContext/main.go


## To check tool calling
go run cmd/scratch/toolCalling/main.go

#reflex go run cmd/scratch/skill/skillContextRead/main.go
# reflex $(cat .reflex) -- go run cmd/scratch/skill/skillContextRead/main.go