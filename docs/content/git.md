# Git

This component allows you to track remote repositories (on GitHub or elsewhere) and do a find-and-replace in your repo when new tags are pushed.

This is especially useful for dependencies that don't use a package manager.

## Example `deps.yml`

```yaml
version: 3
dependencies:
- type: git
  settings:
    remotes:
      https://github.com/kubernetes/minikube.git:
        replace_in_files:
        - filename: dev/install.go
          # pattern is a regex that must have 1 capture group
          pattern: minikube version (\S+)
          # you can optionally disable semver parsing on the tags
          # (which means latest version will be the last tag)
          semver: false

      https://github.com/dropseed/deps-git.git:
        replace_in_files:
        - filename: file.txt
          pattern: deps-git (\S+)
          # use a semver range to limit updates
          # https://github.com/blang/semver#ranges
          range: "< 1.0.0"

      https://github.com/getsentry/sentry-javascript.git:
        replace_in_files:
        - filename: file.txt
          pattern: raven==(\S+)
          # only use tags with this prefix (and remove the prefix so we just get the version number)
          tag_prefix: raven-js@
          # include semver pre-releases
          prereleases: true
```

## Support

Any questions or issues with this specific component should be discussed in [GitHub issues](https://github.com/dropseed/deps-git/issues).

If there is private information which needs to be shared then please use the private support channels in [dependencies.io](https://www.dependencies.io/contact/).
