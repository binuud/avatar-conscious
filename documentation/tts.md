## TTS


For TTS, till I find a placeholder
```
 docker run -p 8880:8880 ghcr.io/remsky/kokoro-fastapi-cpu
```

```
 docker run -it  -p 8880:8880 ghcr.io/remsky/kokoro-fastapi-cpu bash
```

The above opens a portal on http://localhost:8880/web/

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

Currently, tts is rung from ai-modules project, python. It listens to MQTT and saves the output as wav file and sends it to mqtt broker

#### External Non Working Projects

* https://github.com/PierrunoYT/Kokoro-TTS-Local/tree/master Kokorro Docker Image (Not working on mac)