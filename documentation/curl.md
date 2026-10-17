# Curl command for most of the api's


### check llm chat completion
```
curl http://localhost:11434/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d '{
    "model": "gemma4:latest",
    "messages": [
      {
        "role": "system",
        "content": "You are a helpful assistant."
      },
      {
        "role": "user",
        "content": "Hello!"
      }
    ]
  }'
```

### Check via proxy
```
curl http://localhost:9080/llmHandler/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ollama" \
  -d '{
    "model": "gemma4:latest",
    "messages": [
      {
        "role": "system",
        "content": "You are a helpful assistant."
      },
      {
        "role": "user",
        "content": "Hello!"
      }
    ]
  }'

```

### Get models
```
curl http://localhost:9080/llmHandler/v1/models
```

```
curl http://localhost:9080/api/v1/conscious/models
```

### In request test
```
curl -N -X POST http://localhost:9080/api/v1/conscious/in \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ollama" \
  -d '{
    "model": "mlx-community/gemma-4-12b-it-4bit",
    "prompt": "how are you",
    "stream": true
  }'
```

### Stream request test
```
curl -N -X POST http://localhost:9080/api/v1/conscious/inStream \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ollama" \
  -d '{
    "model": "mlx-community/gemma-4-12b-it-4bit",
    "prompt": "how are you",
    "stream": true
  }'
```

### Generate request with inbuilt router

Normal chat streaming
```
curl -N -X POST http://localhost:9080/api/v1/conscious/generate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ollama" \
  -d '{
    "model": "mlx-community/gemma-4-12b-it-4bit",
    "prompt": "how are you",
    "stream": true
  }'
```

Image generation
```
curl -N -X POST http://localhost:9080/api/v1/conscious/generate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ollama" \
  -d '{
    "uuid": "bddd24ba-0cd9-4be2-a027-bad056cd31f8"
    "model": "mlx-community/gemma-4-12b-it-4bit",
    "prompt": "Create an image of a cat jumping over moon",
    "stream": true
  }'
```

### Stream system metrics 
```
curl -N -X POST http://localhost:9080/api/v1/consciousMetrics/getSystemMetrics \
  -H "Content-Type: application/json" 
```

```
## model that supports tool calling
export LLM_MODEL="mlx-community/gemma-4-12b-it-4bit" 
```