# Testing a dependencies.io component

The deps binary comes with a testing framework for components. You can install
it on your machine and in CI to ensure that your component gives the proper
output based on specified input.

## Installing deps on your machine

```console
$ curl https://deps.app/install.sh | bash -s 2.5.0-beta.1
```

Or, you can download binaries manually from the [releases on
GitHub](https://github.com/dropseed/deps/releases).

## Writing tests

Test cases are written in a YAML file named `dependencies_test.yml`. These will
be discovered automatically. The YAML format is as follows...

### Collector tests

```yaml
cases:
- name: basic
  type: collector
  dockerfile: Dockerfile
  # path from the repo root to a directory that will be mounted as the test repo
  # (it will be turned into a git repo automatically)
  repo_contents: tests/collector/basic/repo
  output_data_path: tests/collector/basic/expected_output_data.json
  # user_config is only required if you need a path other than ".",
  # or need to test settings
  user_config:
    path: /
    settings:
      example_string: foo
      example_list:
        - foo
        - bar
```

### Actor tests

```yaml
cases:
- name: basic
  type: actor
  dockerfile: Dockerfile
  repo_contents: tests/actor/basic/repo
  # input_data_path is required when using "type: actor"
  input_data_path: tests/actor/basic/input_data.json
  output_data_path: tests/actor/basic/expected_output_data.json
  user_config:
    path: /
    settings:
      example_string: foo
      example_list:
        - foo
        - bar
```

Your file structure might look something like this:
```console
example-component/tests/
├── actor
│   ├── basic
│   │   ├── expected_output_data.json
│   │   ├── input_data.json
│   │   └── repo
│   │       ├── example_lockfile.json
│   │       └── example_manifest.json
│   └── dependencies_tests.yml
└── collector
    ├── basic
    │   ├── expected_output_data.json
    │   └── repo
    │       ├── example_lockfile.json
    │       └── example_manifest.json
    └── dependencies_tests.yml
```

## Running tests

From the root of your repo:
```console
$ deps dev test
```

You can read about all of the test options with `deps dev test --help`, but one
of the more commonly used options is to automatically update your
`output_data_path` files. This makes it easy to save the output that your
component generates (especially when dealing with lockfiles). Be sure to
actually review the changes though and ensure it is still generating the output
you expect.

If your component deals with lockfiles that tend to update frequently (like
`yarn.lock`), you may want to use `deps dev test --loose-output-data-comparison`
while running in CI. This does a less accurate comparison of the output, but
makes it a little easier to keep your CI pipeline running and will still catch
major issues.

A `.travis.yml` might look like this:
```yaml

```
