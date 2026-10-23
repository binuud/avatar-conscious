---
name: sample
uuid: aef82d62-a5d2-4f27-8502-e7dea7735ed3
description: This is a sample avatar file for reference.
triggers:
  - election
  - sir
ui:  
  uuid: ''
  name: First Pane
  description: This is a test pane
  title: Avatar Pane
  components:
  - id: firstName
    class: ''
    label: First Name
    value: ''
    placeHolder: Enter your first name
    help: ''
    style: ''
    width: 50
    rows: 0
    cols: 0
    type: INPUT
  - id: lastName
    class: ''
    label: Last Name
    value: ''
    placeHolder: Enter your last name
    help: ''
    style: ''
    width: 50
    rows: 0
    cols: 0
    type: INPUT
  - id: ''
    class: ''
    label: ''
    value: lorem ipsum dolor sit amet consectetur adipiscing elit optio ducimus temporibus
      ullamco ipsum dolor omnis aut et fugiat fugiat et voluptas et vero est autem minim
      possimus et libero facere imperdiet nihil consequatur eiusmod nostrud voluptas
      vero aliquip possimus duis qui dolorum et cupiditate elit minus quas ea voluptas
      consequatur
    placeHolder: ''
    help: ''
    style: ''
    width: 100
    rows: 0
    cols: 0
    type: PLAINTEXT  
---

# Sample

## Step 1: Prompt user to upload a license image, or show license on webcam

## Step 2: Get information on the user from the uploaded image

Use the following prompt, to get information from the image
```prompt
Describe the image from file sampleLicense.jpg, and get the license number and driver name, in format {'license': [LICENCE NUMBER],'name':'[NAME OF PERSON']}
```

## Step 3: Verify that the uploaded document is valid, by verifying it against an api call

Use the api call to endpoint 
```api
http://localhost:8081/checklicense/[LICENSE NUMBER]
```

Compare the returned user name with the one got from the uploaded image