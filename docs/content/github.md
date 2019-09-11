---
title: "GitHub"
description: "Examples and instructions for setting up deps in GitHub"
---

# GitHub

The easiest way to give deps write-access is with a personal access token --
either from your personal account or a bot account that your company already has set up.

You can also create or use an existing internal GitHub App.
The setup is more involved but offers a different way to manage permissions and repo access.

## Personal access token

1. Log in with the account you want deps to use (this will be the author of deps pull requests)
1. Give it access to the repo you're setting up
1. [Generate a new token](https://github.com/settings/tokens) with the `repo` scope
1. Set the `DEPS_GITHUB_TOKEN` environment variable in your CI

## GitHub App

> We may publish more detailed instructions for this in the future.
[Contact us if you have questions or need help.](https://www.dependencies.io/contact/)

1. [Create an internal GitHub app in your organization](https://developer.github.com/apps/building-github-apps/creating-a-github-app/)
1. Give it access to the repo you're setting up
1. Set the required environment variables in your CI
    - `DEPS_GITHUB_APP_KEY` - base64 encoded private key
    - `DEPS_GITHUB_APP_ID` - your app ID
    - `DEPS_GITHUB_APP_INSTALLATION_ID` - the ID for the installation in your org

## Pull request settings

When working with a GitHub repo,
there are a few settings you can use to determine what your PRs look like.

```yaml
version: 3
dependencies:
- type: python
  settings:
    github_labels: ["automerge"]
    github_base_branch: test
    github_assignees: ["user1"]
    github_milestone: 1
```

If you don't need a `deps.yml` then you can also configure these settings via environment variables.
This is an easy way to put settings directly in your CI config.

Note that they'll need to be in the format of a JSON-encoded string,
with an uppercase name prefixed by `DEPS_SETTING_`.

```sh
$ DEPS_SETTING_GITHUB_LABELS='["automerge"]' deps ci
```
