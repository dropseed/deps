# pullrequests

A simple Go application that can send pull requests to the different Git hosts.
Can be easily installed inside Docker containers of any kind, so that they don't
have to implement this functionality themselves.

This was built for being run within the [dependencies.io](https://www.dependencies.io) platform, though it could probably used however. The key though is the environment variables, which match up with what [dependencies.io](https://www.dependencies.io) uses.

## Usage in a dependencies.io actor

Add to your Dockerfile.

```sh
ADD https://github.com/dependencies-io/pullrequest/releases/download/0.1.0/pullrequest_0.1.0_linux_amd64.tar.gz /usr/src/actor/pullrequest
```

Do your `git commit` and `git push`, so that the branch exists on the repo.

Then use `pullrequest`, and be sure to properly escape your strings (ex. `shlex.quote` in [python](https://docs.python.org/3/library/shlex.html#shlex.quote)).

```sh
/usr/src/actor/pullrequest/pullrequest --branch="branch-name" --title="PR title" --body="PR body"
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

## Development

```sh
go build && env $(cat .env | xargs) ./pullrequest --branch="redux-3.7.2-11.1.0" --title=test --body="Testing it out"
```
