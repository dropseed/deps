# Using deps in CI

Most CI systems have a cron or scheduled job function,
which is how we recommend using `deps ci`.
This way you can decide how often you want updates (daily, weekly, monthly, etc.) and keep slow dependency updates out of your normal commit pipeline.

The basic steps for setting up `deps ci` are:

- Copy your `DEPS_TOKEN` from [3.dependencies.io](https://3.dependencies.io)
- Add a scheduled job that installs and runs `deps ci`
- Set any other required environment variables and secrets

Read on for specific instructions for your repository host and CI 👇

### Git hosts

- [GitHub](/github/)

### CI providers

- [CircleCI](/circleci/)
- [TravisCI](/travisci/)
