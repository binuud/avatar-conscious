---
name: tourist-camera
uuid: cf13858e-4e40-4f90-a792-ba5b09c63651
image: /static/avatarStore/tourist-camera.png
description: This is a reference demo, to showcase webcamera image piple line to agent. It takes a photo of a person using webcamera, and generates a tourist photo of random place.
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
    value: Avatar will try to generate a tourist photo of you.
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
Launch an Avatar UI interface, and get an image from the webcamera, that will be used by LLM to create a tourist photo
```

## Step 2: Generate an image using the image from web camera as reference

Use the following prompt to describe the image
```prompt
Create an image, using the image of the person from webcamera as reference, the image filename is given in the context, the person should be standing facing the camera, in front of a famouse tourist place.
```
