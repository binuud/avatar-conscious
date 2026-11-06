---
name: dress-shop
uuid: 23cb7a34-556f-466d-be56-5e98a9996f7f
image: /static/avatarStore/dress-shop.png
description: This is a reference demo, to showcase multi form input
triggers:
  - webcam
---
# Dress Shop

## Step 1: 

```Description
Launch an Avatar UI web-camera interface
```

```prompt
Launch an Avatar UI interface with interface id "web-camera"
```

## Step 2: 

```Description
Launch an Avatar UI select-dress interface
```

```prompt
Launch an Avatar UI interface with interface id "select-dress"
```


## Step 3: 

```Description
Generate an image of the you wearing the selected dress of your choice
```

Use the following prompt to describe the image
```prompt
Edit the given image, using the image of the person from webcamera as reference, the image filename is given in the context. 
Do not modify the face of the person.
Keep the facial features, skin tone of the person exactly as they are in the original image.
Create a new photo with the person wearing this clothes as given in Step 2, "select-dress" objec entry
```
