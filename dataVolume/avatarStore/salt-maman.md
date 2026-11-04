---
name: salt-maman
uuid: 995865f9-6849-4219-935e-a5b770b7a7bf
image: /static/avatarStore/salt-maman.png
description: This is a reference avatar to showcase persistant memory over avatar calls.
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
    value: Avatar will try to guess the country of origin of the name and the gender of the name.
    placeHolder: ''
    help: ''
    style: ''
    width: 100
    rows: 0
    cols: 0
    type: PLAINTEXT  
---

# Salt Maman

## Step 1: 

```Description
Launch an Avatar UI interface.
```

Use the following prompt, to launch ui
```prompt
Launch an Avatar UI interface. And get the firstname and lastname of user.
```

## Step 2:

```Description
Try guessing the place of origin or country of the user
```

Use the following prompt, to get information about the user name
```prompt
The name of the user is [firstName] [lastName], try your best to guess the country of the name and the gender of the name.
```
