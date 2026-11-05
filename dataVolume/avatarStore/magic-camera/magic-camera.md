---
name: magic-camera
uuid: 3d660f8b-f73b-42d2-84fc-d1cfc5daf571
image: /static/avatarStore/magic-camera.png
description: This is a reference demo, to showcase multi form input
triggers:
  - webcam
---
# Magic Camera

## Step 1: 

```Description
Launch an Avatar UI web-camera interface
```

```prompt
Launch an Avatar UI interface, the ui interface id is : web-camera
```

## Step 2: 

```Description
Launch an Avatar UI get-prompt interface
```

```prompt
Launch an Avatar UI interface, the ui interface id is : get-prompt
```


## Step 3: 

```Description
Generate an image with the input prompt and given image as reference
```

Use the following prompt to describe the image
```prompt
Edit the given image, using the image of the person from webcamera as reference, the image filename is given in the context. Use the prompt user-prompt given in the context Step No: 2.
```
