---
name: name-detective
uuid: aef82d62-a5d2-4f27-8502-e7dea7735ed3
image: /static/avatarStore/name-detective.png
description: This is a sample avatar agent file for reference. It prompts the user to enter a first and last name. Then tries to guess the country of origin and gender.
triggers:
  - election
  - sir
---

# Country Guesser

## Step 1: 

```Description
Launch an Avatar UI interface.
```

Use the following prompt, to launch ui
```prompt
Launch an Avatar UI interface with id get-name. And get the firstname and lastname of user.
```

## Step 2:

```Description
Try guessing the place of origin or country of the user
```

Use the following prompt, to get information about the user name
```prompt
The name of the user is [firstName] [lastName], try your best to guess the country of the name and the gender of the name.
```
