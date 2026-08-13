## WORKING as on 12th Aug 2026. need to check with golang on api implementation
## this shows docker model runner has image url support for vision model

Q="What do you see in this image?"
IMG="$(base64 -i ./dataVolume/sampleImages/elephantWithBall.jpg)"
# echo 'Test '"${Q}"' finished  '"${IMG}"' img'
curl http://localhost:12434/engines/llama.cpp/v1/chat/completions \
-H "Content-Type: application/json" \
-d "{
  \"model\": \"ai/gemma3:4B-Q4_K_M\",
  \"messages\": [
    {
      \"role\": \"user\",
      \"content\": [
        {
          \"type\": \"text\",
          \"text\": \"${Q}\"
        },
        {
          \"type\": \"image_url\",
          \"image_url\": {\"url\": \"data:image/jpeg;base64,${IMG}\"}
        }
      ]
    }
  ]
}"
