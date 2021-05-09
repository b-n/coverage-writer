# Coverage Writer

Serverless app (golang!) for writing and retrieving test results from a Firestore DB in GCP.

The main idea is to store data via org/repo:branch with language, and eventually to be able to pull
this data via a frontend to track coverage over time.

Why not codecov? Because why pay for something when you can make a simplier, more configurable
version yourself.

## Development

Deploying to GCP to test each function is painfully slow, and if you happen to be writing and
testing a lot of code, you could go over your limits. Instead, use the included `cmd` to make your
life easier :tm:.

```sh
go run cmd/main.go
```

then
```sh
curl http://localhost:8080/<function name>
```

Note: Any changes to your source will need to restart the server (sigh)


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
