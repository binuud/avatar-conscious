
# export KEY=""
export LLM_HOST="http://host.docker.internal:11434"
export LLM_MODEL="gemma4:latest"
export MQTT_HOST="tcp://host.docker.internal:1883"
export FILE_STORE_PATH=${HOME}/Code/avatar/avatar-avatarconscious/dataVolume/
export SERVABLE_DIRECTORY=${HOME}/Code/avatar/avatar-avatarconscious/dataVolume/gen/

reflex $(cat .reflex) -- go run cmd/scratch/main.go