# Using deps in CI

Most CI systems have a cron or scheduled job function,
which is how we recommend using `deps ci`.
This way you can decide how often you want updates (daily, weekly, monthly, etc.) and keep slow dependency updates out of your normal commit pipeline.
This way you'll also get dependency updates even if you haven't committed to the repo recently.

The basic steps for setting up `deps ci` are:

- Copy your `DEPS_TOKEN` from [3.dependencies.io](https://3.dependencies.io)
- Add a scheduled job that installs and runs `deps ci`
- Set any other required environment variables and secrets

👈 Find the instructions for your repository host and CI in the menu.

## Autoconfigure

By default, `deps ci` automatically detects and configures various settings for common CI providers.
This includes setting `git config` user name and email,
and changing the git remote from ssh to https.

To disable the automatic configuration (if you have other requirements or want to do it yourself) use `deps ci --manual`.

For specifics on what is configured and how, you can [read the code here](https://github.com/dropseed/deps/search?l=Go&q=autoconfigure).

## Filtering by type

In some container-based CI systems,
you'll only have certain languages and requirements installed in certain containers.

You can use the `--type` option to run the appropriate updates based on the container you're in. For example, use `deps ci --type js` in your container with your JavaScript environment and `deps ci --type python` in your Python container.
