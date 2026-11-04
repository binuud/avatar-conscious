---
name: camera-describer
uuid: 10f15989-26f7-4264-9add-0fc46cf9049e
image: /static/avatarStore/camera-describer.png
description: This is a simple demo, that uses LLM to describe the scene visible in the web camera.
triggers:
- webcam
---
# Camera Describer

## Step 1: 

```Description
Launch an Avatar UI webcamera interface.
```

```prompt
Launch an Avatar UI interface with id web-camera, and get an image from the webcamera, that will be used by LLM to describe the image.
```

## Step 2: 

```Description
Describe the image using a vision model
```

Use the following prompt to describe the image
```prompt
Describe the attached image from webcamera in detail
```
