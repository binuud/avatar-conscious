---
name: camera-describer
uuid: 10f15989-26f7-4264-9add-0fc46cf9049e
image: /static/avatarStore/camera-describer.png
description: This is a simple demo, that uses LLM to describe the scene visible in the web camera.
triggers:
  - webcam
ui:  
  uuid: ''
  name: First Pane
  description: Country Guesser
  title: I can guess your country
  components:
  - id: webcamview
    class: ''
    label: WebCamera
    value: ''
    placeHolder: ''
    help: ''
    style: ''
    width: 100
    rows: 0
    cols: 0
    type: WEBCAM
  - id: ''
    class: ''
    label: ''
    value: Avatar will try to guess in detail, what it is seeing from the webcamera.
    placeHolder: ''
    help: ''
    style: ''
    width: 100
    rows: 0
    cols: 0
    type: PLAINTEXT  
---

# Camera Describer

## Step 1: Launch an Avatar UI webcamera interface.

Use the following prompt, to launch ui with webcamera
```prompt
Launch an Avatar UI interface, and get an image from the webcamera, that will be used by LLM to describe the image.
```

## Step 2: Describe the image using a vision model

Use the following prompt to describe the image
```prompt
Describe the attached image from webcamera in detail
```
