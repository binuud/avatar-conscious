
# export KEY=""
export LLM_ENDPOINT="http://127.0.0.1:11434"
export MQTT_HOST="tcp://localhost:1883"
export SERVABLE_DIRECTORY=${HOME}/Code/ai/aiOutData/dataVolume/generatedAiFiles/
reflex $(cat .reflex) -- go run cmd/conscious/main.go