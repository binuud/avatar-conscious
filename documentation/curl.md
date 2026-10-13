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

### Stream request test
```
curl -N -X POST http://localhost:9080/api/v1/conscious/inStream \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ollama" \
  -d '{
    "model": "gemma4:latest",
    "prompt": "how are you",
    "stream": true
  }'
```

### Stream system metrics 
```
curl -N -X POST http://localhost:9080/api/v1/consciousMetrics/getSystemMetrics \
  -H "Content-Type: application/json" 
```
