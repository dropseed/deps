---
title: "GitHub"
description: "Examples and instructions for setting up deps in GitHub"
---

# GitHub

To use deps on a GitHub repo,
you'll need to make sure it has permission to commit and push from CI.

Most CI provders use read-only access for the normal cloning and testing. The easiest way to give deps write-access is with a personal access token --
either from your personal account or a bot account that your company already has set up.

You can also create or use an existing internal GitHub App.
The setup is more involved but offers a different way to manage permissions and repo access.

> Note: If your CI is already configured with write-access and you have others tools that commit back to the repo,
you can use those same credentials for deps if you want!

## Personal access token

1. Log in with the account you want deps to use (this will be the author of deps pull requests)
1. Give it access to the repo you're setting up
1. [Generate a new token](https://github.com/settings/tokens) with the `repo` scope
1. Set the required environment variables in your CI
    - `GITHUB_TOKEN` or `DEPS_GITHUB_TOKEN`

## GitHub App

> We may publish more detailed instructions for this in the future.
[Contact us if you have questions or need help.](https://www.dependencies.io/contact/)

1. [Create an internal GitHub app in your organization](https://developer.github.com/apps/building-github-apps/creating-a-github-app/)
1. Give it access to the repo you're setting up
1. Set the required environment variables in your CI
    - `DEPS_GITHUB_APP_KEY` - base64 encoded private key
    - `DEPS_GITHUB_APP_ID` - your app ID
    - `DEPS_GITHUB_APP_INSTALLATION_ID` - the ID for the installation in your org
