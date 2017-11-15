#!/bin/sh
go build ./cmd/pullrequest && env $(cat .env.$1 | xargs) ./pullrequest --branch="react-16.0.0-11.0.0" --title="Deps tester" --body="Testing it out" --related-pr-title-search "Deps tester"
