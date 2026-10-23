---
name: sample
uuid: aef82d62-a5d2-4f27-8502-e7dea7735ed3
description: This is a sample avatar file for reference. It prompts the user to enter a first and last name.
triggers:
  - election
  - sir
ui:  
  uuid: ''
  name: First Pane
  description: Country Guesser
  title: I can guess your country
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

## Step 1: Prompt user to enter his First Name and Last name

## Step 2: Try guessing the place of origin or country of the user

Use the following prompt, to get information about the user name
```prompt
The name of the user is [firstName] [lastName], which country does he belong to
```
