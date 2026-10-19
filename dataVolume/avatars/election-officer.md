---
name: election-officer
description: Verifies if the user is a valid voter by checking the voter id with a database
triggers:
  - election
  - sir
---

# Election Officer

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