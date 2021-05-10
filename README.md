# Coverage Writer

Serverless app (golang!) for writing and retrieving test results from a Firestore DB in GCP.

The main idea is to store data via org/repo:branch with language, and eventually to be able to pull
this data via a frontend to track coverage over time.

Why not codecov? Because why pay for something when you can make a simplier, more configurable
version yourself.

## Development

You have options:

1. After every change, deploy to GCP, and test it that way (pitch drop experiment slow)
2. Run an emulated firebase locally, boot a go server locally, curl against that (much faster)

### Deploying development to GCP

```sh
sls deploy --stage dev
```

### Run locally

Dependencies:

- firebase emulator<br>
  `curl -sL https://firebase.tools | bash`
  `firebase login`
- Patience

You need a couple of terminals:
- One to run firebase emulator: `firebase emulators:start --import seed`
- One to run the go server: `go run cmd/main.go`
- One to curl the things: `curl http://localhost:8000/<function name>`

Note: Want to see some debugs from firebase? `tail -f firebase-debug.log`

## Routes

GET "/"

Query Params:
- org (required)
- repo
- branch
- language
- from (required)
- to (required)

POST "/"
Body: CoverageData (required)
