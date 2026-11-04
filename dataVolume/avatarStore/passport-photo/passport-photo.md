---
name: passport-photo
uuid: 066c61f2-7caa-480e-8c55-948404469a3a
image: /static/avatarStore/passport-photo.png
description: This is a reference demo, to showcase webcamera, external utility pipeline.
triggers:
  - webcam
---
# Passport Photo Maker

## Step 1: 

```Description
Launch an Avatar UI webcamera interface.
```

Use the following prompt, to launch ui with webcamera
```prompt
Launch an Avatar UI interface with id web-camera, and get an image from the webcamera, that will be used by LLM to create a passport style photo
```

## Step 2: 

```Description
Generate an passport style image using the image from web camera as reference
```

Use the following prompt to describe the image
```prompt
Edit the given image, remove the background, change it to blue, make a passport photo of the person in the image, wearing suit, do not change skin tone, do not smoothen skin.
```



