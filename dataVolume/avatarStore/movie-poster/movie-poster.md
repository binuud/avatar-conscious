---
name: movie-poster
uuid: e7e578aa-16d1-4095-97a0-8a4ad8a57a3f
image: /static/avatarStore/movie-poster.png
description: Sample Avatar, showcasing prompt amplification, with file injection at later steps
triggers:
  - webcam
---
# Movie-Poster

## Step 1: 

```Description
Launch an Avatar UI web-camera interface
```

```prompt
Launch an Avatar UI interface, the ui interface id is : web-camera
```

## Step 2: 

```Description
Launch an Avatar UI select-movie interface
```

```prompt
Launch an Avatar UI interface, the ui interface id is : select-movie
```


## Step 3: 

```Description
Generate an image with the selected movie name or movie plot and given image as reference
```

Use the following prompt to describe the image
```prompt
Edit the given image, using the image of the person from webcamera as reference, the image filename is given in the context. 
Keep the facial features of the person exactly as they are in the original image.
Change clothes as per the movie character. Design a movie poster with the person as the main character.
Use the movie name or movie plot given in select-movie in Step No.2.
Do not make mistakes in movie name, spell it correctly as given.
```
