
# export KEY=""
export LLM_ENDPOINT="http://127.0.0.1:11434"
reflex $(cat .reflex) -- go run cmd/conscious/main.go