#!/bin/sh
go build ./cmd/pullrequest && env $(cat .env.$1 | xargs) ./pullrequest --branch="redux-3.7.2-11.1.0" --title=test --body="Testing it out"
