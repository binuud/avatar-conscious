# conscious
Conscious means being awake, alert, and aware of one's surroundings, thoughts, or actions. Manomaya Kosha.

To give a body, that is not affected in the original sin and the secondary sin. The first life form, that need not consume another life form for its survival and upkeeping.

* Adds memory persistence to llm
* Support tooling
* Support image description
* LLM chat completion

This has to be fronted with a UI.   

## Usage

Please refer to [Avatar Project](https://github.com/binuud/avatar), for usage, installation details.

Metrics to work, (only on macos)

## Directory Structure

* cmd - the main go file
* documentation - all documentation md files
* pkg/servers - all grpc server implementations go here
* proto/v1 - all protofiles go here
* scripts - support scripts, like launchers etc
* vendor - go vendor files, do not make changes here

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

```
go install github.com/cespare/reflex@latest
```


## Remember

* Avoid using containers during development, since metrics, needs macos access
* Running in docker container, might not show relevant metrics

## Credits


