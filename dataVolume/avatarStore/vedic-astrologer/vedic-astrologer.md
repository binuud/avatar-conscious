---
name: vedic-astrologer
uuid: b44daf15-a204-436e-b693-6e310321bba8
image: /static/avatarStore/vedic-astrologer.png
description: Sample Avatar, with a simple multi field form. Uses MCP Server to show case white listed tool usecase.
triggers:
  - election
  - sir
mcp_servers:
  - astro-server
---

# Vedic Astrologer

## Step 1: 

```Description
Launch an Avatar UI interface.
```

```prompt
Launch an Avatar UI interface with id get-birth-details. To get the date of birth and time of birth
```

## Step 2:

```Description
Get the astrological chart of the user with the date and time of birth details.
```

Use the following prompt, to get astrological chart of the user
```prompt
The date of birth of the user is [dob] and time of birth is [tob], get the astrological chart for the user.
```

## Step 3:

```Description
Predict important life events and give insights about the users life, with the astrological information.
```

Use the following prompt, to get information about the user name
```prompt
Give insights about the user, with the astrological details available to you.
```