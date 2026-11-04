---
name: tourist-camera
uuid: cf13858e-4e40-4f90-a792-ba5b09c63651
image: /static/avatarStore/tourist-camera.png
description: This is a reference demo, to showcase webcamera image piple line to agent. It takes a photo of a person using webcamera, and generates a tourist photo of random place.
triggers:
  - webcam
---

# Tourist Photographer

## Step 1:

```Description
Launch an Avatar UI webcamera interface.
```

Use the following prompt, to launch ui with webcamera
```prompt
Launch an Avatar UI interface with id web-camera, and get an image from the webcamera, that will be used by LLM to create a tourist photo
```

## Step 2:

```Description
Generate an image using the image from web camera as reference
```

Use the following prompt to describe the image
```prompt
Edit the given image, remove the background. Then create a new image, using the image of the person from webcamera as reference, the image filename is given in the context, the person should be standing facing the camera, in front of a place like any of the seven wonders, a beach, big cities, river front. Choose a place at random. Style the clothing depending on the tourist place. Do not alter the persons skin tone or facial features.
```
