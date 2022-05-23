---
title: "GitHub"
description: "Examples and instructions for setting up deps in GitHub"
---

# GitHub

Deps is designed to run directly inside of GitHub Actions, using a workflow very similar to what you would use for tests.

You'll need two things:

1. A `DEPS_TOKEN` from [3.dependencies.io](https://3.dependencies.io) (set this as a [repo](https://docs.github.com/en/actions/security-guides/encrypted-secrets#creating-encrypted-secrets-for-a-repository) or [org](https://docs.github.com/en/actions/security-guides/encrypted-secrets#creating-encrypted-secrets-for-an-organization) "secret").
1. A `DEPS_GITHUB_TOKEN` that is either a [personal access token](#personal-access-token) or a [GitHub App token](#github-app-token).

```yaml
# .github/workflows/deps.yml
name: deps

on:
  schedule:
  - cron: 0 0 * * Mon  # Weekly
  # cron: 0 0 * * *  # Daily
  # cron: 0 0 1 * *  # Monthly
  workflow_dispatch: {}

jobs:
  deps:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2

    # Configure your languages and package managers
    - uses: actions/setup-python@v2
      with:
        python-version: 3.8
    - run: |
        pip install -U pip pipenv

    # Install and run deps
    - run: curl -sSL https://deps.app/install.sh | bash -s -- -b $HOME/bin
    - run: $HOME/bin/deps ci
      env:
        DEPS_TOKEN: ${{ secrets.DEPS_TOKEN }}
        DEPS_GITHUB_TOKEN: ${{ secrets.DEPS_GITHUB_TOKEN }}
```

## GitHub API access

Deps will need access to the GitHub API to manage pull requests.
The easiest way to do this is with a [personal access token](#personal-access-token),
but a [GitHub App](#github-app-token) is recommended for larger organizations.

We generally recommend *not* using the default `${{ secrets.GITHUB_TOKEN }}`,
as it won't trigger your other workflows to run on the deps commits.

### Personal access token

1. Log in with the account you want deps to use (this will be the author of deps pull requests)
1. Give it access to the repo you're setting up
1. [Generate a new token](https://github.com/settings/tokens) with the `repo` scope
1. Add the new token as a repo or organization `DEPS_GITHUB_TOKEN` "secret"
1. Set the `DEPS_GITHUB_TOKEN` environment variable in your CI as `${{ secrets.DEPS_GITHUB_TOKEN }}`

```yaml
# .github/workflows/deps.yml
    - run: $HOME/bin/deps ci
      env:
        DEPS_GITHUB_TOKEN: ${{ secrets.DEPS_GITHUB_TOKEN }}
        ...
```

### GitHub App token

> We may publish more detailed instructions for this in the future.
[Contact us if you have questions or need help.](https://www.dependencies.io/contact/)

1. [Create an internal GitHub app in your organization](https://developer.github.com/apps/building-github-apps/creating-a-github-app/)
1. Give it access to the repo you're setting up
1. Set the required environment variables in your "secrets" and then in your workflow
    - `DEPS_GITHUB_APP_KEY` - base64 encoded private key
    - `DEPS_GITHUB_APP_ID` - your app ID
    - `DEPS_GITHUB_APP_INSTALLATION_ID` - the ID for the installation in your org

Here's an example workflow that uses another GitHub action for generating the app token:

```yaml
# .github/workflows/deps.yml
    steps:
    - id: generate_token
      uses: tibdex/github-app-token@v1
      with:
        app_id: ${{ secrets.DEPS_GITHUB_APP_ID }}
        private_key: ${{ secrets.DEPS_GITHUB_APP_KEY }}

    - uses: actions/checkout@v2
      with:
        token: ${{ steps.generate_token.outputs.token }}

    - run: curl https://deps.app/install.sh | bash -s -- -b $HOME/bin
    - run: $HOME/bin/deps ci
      env:
        DEPS_TOKEN: ${{ secrets.DEPS_TOKEN }}
        DEPS_GITHUB_TOKEN: ${{ steps.generate_token.outputs.token }}
```

## Pull request settings

When working with a GitHub repo,
there are a few settings you can use to determine what your PRs look like.

```yaml
# deps.yml
version: 3
dependencies:
- type: python
  settings:
    github_labels: ["dependencies"]
    github_base_branch: test
    github_assignees: ["user1"]
    github_milestone: 1
```

If you don't need a `deps.yml` then you can also configure these settings via environment variables.
This is an easy way to put settings directly in your CI config.

Note that they'll need to be in the format of a JSON-encoded string,
with an uppercase name prefixed by `DEPS_SETTING_`.

```console
$ DEPS_SETTING_GITHUB_LABELS='["dependencies"]' deps ci
```

### Automerge

You can use `github_automerge` to enable the [GitHub auto-merge setting](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/automatically-merging-a-pull-request)
(make sure this is enabled on the repo).
If you use `github_automerge` and the pull request does not have any branch requirements (or has already met them),
then deps will immediately merge the PR manually.

The `github_automerge` setting should be set to one of the available merge methods (`merge`, `squash`, or `rebase`).

```yaml
# deps.yml
version: 3
dependencies:
- type: python
  settings:
    github_automerge: squash
```
