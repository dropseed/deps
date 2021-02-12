# Bitbucket

To give deps write-access to your repo and pull requests, you'll create an *app password*.
You can use your personal account to do this, or a "bot" account that your team has.

## App password

1. Log in with the account you want deps to use (this will be the author of deps pull requests)
1. Give it access to the repo you're setting up
1. Generate a new app password with the repositories and pull requests `write` scopes
    [![Bitbucket app password settings for deps](/assets/img/screenshots/bitbucket-app-password.png)](/assets/img/screenshots/bitbucket-app-password.png)
1. Set the required environment variables in your CI
    - `DEPS_BITBUCKET_USERNAME` to the user who owns the app password
    - `DEPS_BITBUCKET_PASSWORD` to the app password from above

## Pull request settings

When working with a Bitbucket repo,
there are a few settings you can use to determine what your pull requests look like.

```yaml
version: 3
dependencies:
- type: python
  settings:
    bitbucket_destination: "dev"  # branch name
    bitbucket_close_source_branch: true
    bitbucket_reviewers:
    - uuid: "{504c3b62-8120-4f0c-a7bc-87800b9d6f70}"
```

If you don't need a `deps.yml` then you can also configure these settings via environment variables.
This is an easy way to put settings directly in your CI config.

Note that they'll need to be in the format of a JSON-encoded string,
with an uppercase name prefixed by `DEPS_SETTING_`.

```console
$ DEPS_SETTING_BITBUCKET_CLOSE_SOURCE_BRANCH='true' deps ci
```
