---
name: magic-camera
uuid: 3d660f8b-f73b-42d2-84fc-d1cfc5daf571
image: /static/avatarStore/magic-camera.png
description: Sample avatar, showcasing web cam input, with prompt amplification affecting certain inputs to agents.
triggers:
  - webcam
---
# Magic Camera

## Step 1: 

```Description
Launch an Avatar UI web-camera interface
```

```prompt
Launch an Avatar UI interface with interface id "web-camera"
```

## Step 2: 

```Description
Launch an Avatar UI get-prompt interface
```

```prompt
Launch a new Avatar UI interface with interface id "get-prompt"
```


## Step 3: 

```Description
Generate an image with the input prompt and given image as reference
```

Use the following prompt to describe the image
```prompt
Edit the given image, using the image of the person from webcamera as reference, the image filename is given in the context. 
Do not modify the face of the person.
Keep the facial features, skin tone of the person exactly as they are in the original image.
Use the prompt user-prompt given in the context Step No: 2 "user-prompt" object entry
```
