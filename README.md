# conscious
Conscious means being awake, alert, and aware of one's surroundings, thoughts, or actions. Aham Bramhasmi.

To give a body, that is not affected in the original sin and the secondary sin. The first life form, that need not consume another life form for its survival and upkeeping.

* Adds memory persistence to llm
* Support tooling
* Support image description
* LLM chat completion

This has to be fronted with a UI.   

## Swagger Editor
Once the proto files are created, and the OpenAPI json file is created, add the following to the final openApi json file

```
  "servers": [{
    "url": "http://localhost:9080/api",
    "description": "Testing Server"
  }],
```

#### External Working Projects


## Useful Commands

Check vision model
* ollama version atleast 0.40
* mlx version atleast 0.32.3
```
 ollama run gemma4:latest  ~/Code/avatar/conscious/dataVolume/sampleImages/elephantWithBall.jpg "what is in this image? Answer in 10 words"
```

check llm chat completion
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

Check via proxy
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

get models
```
curl http://localhost:9080/llmHandler/v1/models
```