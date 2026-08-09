
# export KEY=""
export LLM_ENDPOINT="http://127.0.0.1:12434"
export LLM_MODEL="gemma4"
export MQTT_HOST="tcp://localhost:1883"
export FILE_STORE_PATH=${HOME}/Code/ai/conscious/dataVolume/
export SERVABLE_DIRECTORY=${HOME}/Code/ai/aiOutData/dataVolume/generatedAiFiles/
reflex $(cat .reflex) -- go run cmd/conscious/main.go