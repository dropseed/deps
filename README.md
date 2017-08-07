# pullrequests

A simple Go application that can send pull requests to the different Git hosts.
Can be easily installed inside Docker containers of any kind, so that they don't
have to implement this functionality themselves.

```sh
go build && env $(cat .env | xargs) ./pullrequest --branch="redux-3.7.2-11.1.0" --title=test --body="Testing it out"
```

## Environment variables

### Always required

- `GIT_HOST` - "github" or "gitlab"
- `GIT_BRANCH` - the default branch on the repo (usually "master")

### GitHub

##### Required

- `GITHUB_REPO_FULL_NAME`
- `GITHUB_API_TOKEN`

##### Optional

- `SETTING_GITHUB_LABELS` - JSON encoded list of strings
- `SETTING_GITHUB_ASSIGNEES` - JSON encoded list of strings
- `SETTING_GITHUB_MILESTONE` - int of milestone number
