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
Launch an Avatar UI interface with id "web-camera"
```
## Step 2: 

```Description
Launch an Avatar UI select-place interface
```

```prompt
Launch an Avatar UI interface, the ui interface id is : "select-place"
```

## Step 3:

```Description
Generate a tourist photograph using the image from web camera as reference, and tourist place as entered by the user
```

Use the following prompt to describe the image
```prompt
Edit the given image, using the image of the person from webcamera as reference, the image filename is given in the context. 
Do not change the face of the person. 
Do not alter the gender or face of the person.
Keep the facial features of the person exactly as they are in the original image. 
Use the tourist place as mentioned in select-place in Step No.2.
Style the clothing depending on the tourist place. 
```
