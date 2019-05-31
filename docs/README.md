# Building a dependencies.io component

Each dependencies.io component is built as a Docker container. There are only a
few requirements for how that is configured.

- Git should be installed and use the email address "bot@dependencies.io" and name "Dependencies.io Bot". Usually done with `git config --system` in your Dockerfile.
- The repo will be mounted at `/repo`. Your Dockerfile should set `WORKDIR /repo`.
- The Go binary `deps` should be installed.
- Specify the `ENTRYPOINT` to run your code, and expect 1 positional argument (`CMD`) to be a repo-relative path given by the user.
- Using a debian based image is suggested, but not currently required.

Simple example:
```Dockerfile
FROM python:3.7.1

# Add the deps utility
ENV DEPS_VERSION=2.5.0-beta
RUN curl https://www.dependencies.io/install.sh | bash -s -- -b /usr/local/bin $DEPS_VERSION

# Configure git settings for git commit
RUN git config --system user.email "bot@dependencies.io"
RUN git config --system user.name "Dependencies.io Bot"

# Put your code somewhere
RUN mkdir -p /usr/src/app
ADD src/ /usr/src/app/

# Run your code with /repo as the working directory
WORKDIR /repo

ENTRYPOINT ["/usr/src/app/entrypoint.sh"]
```

## Overview

Each component is responsible for parsing the manifests ("package.json",
"requirements.txt", etc.) and lockfiles ("yarn.lock", "composer.lock", etc.)
that they are targeting. That parsed data is output into a pre-determined JSON
schema so that the dependencies data can be passed throughout the system.

The `deps component` commands are frequently used to perform the steps that are
common to every component, like running user-injected commands or creating a
pull request.

One of the best ways to see how components are built is to look at the ones that
already exist! It is usually easier to start simple, and then add in
language-specific features as needed. One of the simpler existing examples is
our custom `git-repos` component
[here](https://github.com/dropseed/git-repos).

### Entrypoint

The entrypoint is usually pretty simple, and is a good place to configure any
component-wide [settings](settings.md).

The `RUN_AS` environment variable will tell you whether to "collect" the
dependencies, or "act" on them.

```python
import os

from collect import collect
from act import act


RUN_AS = os.getenv('RUN_AS')

if os.getenv('SETTING_EXAMPLE') == 'test':
    configure_example_setting()

if RUN_AS == 'collector':
    collect()
elif RUN_AS == 'actor':
    act()
```

### Collect

1. Run `deps component hook before_update`.
1. The dependency path provided by the user will be given to your script as the
argument (ex. `sys.argv[1]`). It can be either a directory or a file, depending
on your implementation. By default it will be `.` which simply points to the
current directory (top of the repo).
1. Manifests and lockfiles should be converted into their respective schemas.
[Take a look at those formats for more info on what needs to be collected.](schema.md) The only
additional data that needs to be found is the "available" versions for the direct dependencies. Again, this can usually be done using the native package manager or registry API in some way.
1. Run `deps component collect <path to JSON file>` to send the results.
1. Done!

Simplified example:

```python
run('deps hook before_update')

path_given = sys.argv[1]

output_data = {
    "manifests": {
        path_given: {
            "current": {
                "dependencies": manifest_dependencies(path_given)
            }
        }
    }
}

if lockfile_exists():
    output_data["lockfiles"] {
        lockfile_path: {
            "current": {
                "dependencies": lockfile_dependencies(lockfile_path),
                "fingerprint": lockfile_md5sum(lockfile_path)
            }
        }
    }

    update_lockfile(lockfile_path)

    if current_lockfile_fingerprint != updated_lockfile_fingerprint:
        output_data["lockfiles"]["updated"] = {
            "dependencies": lockfile_dependencies(lockfile_path),
            "fingerprint": lockfile_md5sum(lockfile_path)
        }

# Send the results to dependencies.io
temp_file_path = write_json_to_temp_file(output_data)
run('deps collect {temp_file_path}')
```

### Act

1. Run `deps component branch` to create a new branch for the update.
1. The JSON for the update will be mounted at `/dependencies/input_data.json`.
1. Using the input_data, you can now perform the updates however you need, and
then use `deps component commit -m "Update xyz to 1.0.0" <paths to commit>` to
commit the changes. Typically, the thing to do at this point is to use the
native tools for the language/ecosystem you're working on to do the actual
update -- mirroring what a user would do on their own machine as closely as
possible. If there are lockfile updates available, those are usually made by
running the lockfile update command, and then committing those changes. To
protect against the lockfile update having changed since it was originally
collected & requested, you can re-collect the new lockfile and update the JSON
with it (see example below).
1. Run `deps component pullrequest <path to JSON file>` to send a pull request
to the hosted repo.

Simplified example:

```python
input_data = json.load(open("/dependencies/input_data.json"))

# Create a new branch for this update
run('deps branch')

for lockfile_path, lockfile_data in data.get('lockfiles', {}).items():
    # If "lockfiles" are present then it means that there are updates to
    # those lockfiles that you can make. The most basic way to handle this
    # is to use whatever "update" command is provided by the package
    # manager, and then commit and push the entire update. You can try to be
    # more granular than that if you want, but performing the entire "update"
    # at once is an easier place to start.

    update_lockfile(lockfile_path)

    # Update lockfile.updated with the data that we have now -- in case
    # it differs from when the original lockfile update was created
    lockfile_data["updated"]["dependencies"] = lockfile_dependencies(lockfile_path)
    lockfile_data["updated"]["fingerprint"] = lockfile_md5sum(lockfile_path)

    # Commit the lockfile update to the repo
    run('deps commit -m "Update {lockfile_path}" {lockfile_path}')

for manifest_path, manifest_data in data.get('manifests', {}).items():
    for dependency_name, updated_dependency_data in manifest_data['updated']['dependencies'].items():
        # We'll update each dependency individually
        installed = manifest_data["current"]["dependencies"][dependency_name]["constraint"]
        new_version = updated_dependency_data["constraint"]

        # Update the dependency in the manifest
        update_dependency(manifest_path, dependency_name, new_version)

        # Commit the manifest change to the repo
        run('deps commit -m "Update {dependency_name} from {installed} to {new_version}" {manifest_path}')

# Send a pull request to the repo
temp_file_path = write_json_to_temp_file(input_data)
run('deps pullrequest {temp_file_path}')
```

## Testing

Once your component is underway, you'll want to add tests to ensure things are
working properly. We have a custom-built testing framework that tests the input
and output JSON data, which in the end is the most important thing. Feel free to
add your own language-specific tests too, but at a minimum you should use the
deps testing framework.

[Read more about testing.](testing.md)

## Schema

A JSON schema is used to validate collected dependencies data. Validation
happens automatically when you use `deps component collect`.

[See the schema.](schema.md)
