## TTS


For TTS, till I find a placeholder
```
 docker run -p 8880:8880 ghcr.io/remsky/kokoro-fastapi-cpu
```



```
curl -X POST http://localhost:8880/v1/audio/speech \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer not-needed" \
     -d '{
       "model": "kokoro",
       "input": "The small pup gnawed a hole in the sock. The fish twisted and turned on the bent hook.",
       "voice": "af_bella",
       "response_format": "mp3",
       "speed": 1.0
     }' \
     -o outdir/speech.mp3
```     



#### External Non Working Projects

* https://github.com/PierrunoYT/Kokoro-TTS-Local/tree/master Kokorro Docker Image (Not working on mac)