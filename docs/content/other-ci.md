# Other CI providers

Don't see what you need?
[Let us know so we can get it added!](https://www.dependencies.io/contact/)

The basic steps for using `deps ci` are generally the same, regardless of provider.

## 1. Set the environment variables

- `DEPS_TOKEN` from [3.dependencies.io](https://3.dependencies.io)
- Credentials for your git host (ex. `GITHUB_TOKEN`)

## 2. Add a scheduled job to run deps

- Use the install script to get the latest version (`curl https://deps.app/install.sh | bash -s -- -b $HOME/bin`)
- Run the `deps ci` command

## 3. Add a `deps.yml` (optional)

Most dependencies and languages we support will be detected automatically,
so you may not even need a `deps.yml`.

But if you need to tweak the settings or point deps to custom dependency locations,
you'll want to [add `deps.yml` to your repo](/config/).
