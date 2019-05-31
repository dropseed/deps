from os import path
import json
import sys
from subprocess import run

from utils import get_lockfile_fingerprint, mock_lockfile_update, write_json_to_temp_file


def collect():
    # A collector will be given the path in the repo that it should collect
    # from. A lot of times, it will simply be the root of the repo. The path can
    # be a to a file or directory, and it's your job to handle the kind of input
    # that you expect. In more predictable dependency managers, you will want to
    # expect a directory (which may have multiple files relating to the
    # dependencies) and you can infer which files to look at because the naming
    # is predictable.
    user_given_path_in_repo = sys.argv[1]

    manifest_path = path.join(user_given_path_in_repo, 'example_manifest.json')
    with open(manifest_path, 'r') as f:
        manifest = json.load(f)

    lockfile_path = path.join(user_given_path_in_repo, 'example_lockfile.json')
    with open(lockfile_path, 'r') as f:
        lockfile = json.load(f)

    # MANIFESTS -----------------------------------------------------------------
    # When collecting manifests, your job is to
    # A) figure out which versions of each dependency are installed, and
    # B) find any potential versions that the user could update to
    #    - in most systems the user can use a "constraint" in their manifest,
    #      which will automatically update them (usually in combination with
    #      their lockfile) to versions within that range, so when collecting
    #      versions for a manifest you should only collect versions
    #      *outside* of their constraint, which would require a change to their
    #      manifest to use
    collected_dependencies = collect_manifest_dependencies(manifest, lockfile)

    manifest_output = {
        "manifests": {
            "example_manifest.json": {
                "current": {
                    "dependencies": collected_dependencies
                }
            }
        }
    }
    # Send the data back to dependencies.io
    run(['deps', 'component', 'collect', write_json_to_temp_file(manifest_output)], check=True)

    # LOCKFILES -----------------------------------------------------------------
    #
    # 1) Collect the current contents of the lockfile
    current_lockfile_dependencies = collect_lockfile_dependencies(lockfile)

    # The lockfile "fingerprint" is a unique string representing the contents
    # of the lockfile at the given time. Some package managers will provide a
    # "content-hash" or similar for this purpose, which you can reuse here.
    # Otherwise you can just get an md5 hash of the file yourself.
    current_lockfile_fingerprint = get_lockfile_fingerprint(lockfile_path)

    lockfile_output = {
        'lockfiles': {
            'example_lockfile.json': {
                'current': {
                    'fingerprint': current_lockfile_fingerprint,
                    'dependencies': current_lockfile_dependencies,
                }
            }
        }
    }

    # 2) Update the lockfile
    # Now that you know what the user's lockfile looked like initially, you can
    # check to see if the lockfile is outdated in any way (i.e. dependencies
    # have been defined with ranges, and new versions are available in those
    # ranges).
    updated_lockfile = mock_lockfile_update(lockfile_path)

    # 3) Collect the updated contents of the lockfile
    # Now we can recollect the updated lockfile
    updated_lockfile_dependencies = collect_lockfile_dependencies(updated_lockfile)
    updated_lockfile_fingerprint = get_lockfile_fingerprint(lockfile_path)

    # If the lockfile has changed, then we can finally we can output everything
    # that we know about the update!
    if current_lockfile_fingerprint != updated_lockfile_fingerprint:
        lockfile_output['lockfiles']['example_lockfile.json']['updated'] = {
            'fingerprint': updated_lockfile_fingerprint,
            'dependencies': updated_lockfile_dependencies,
        }

    # 4) Output the results
    run(['deps', 'component', 'collect', write_json_to_temp_file(lockfile_output)], check=True)


def collect_manifest_dependencies(manifest_data, lockfile_data):
    """Convert the manifest format to the dependencies schema"""
    output = {}

    for dependencyName, dependencyConstraint in manifest_data.items():
        output[dependencyName] = {
            # identifies where this dependency is installed from
            'source': 'example-package-manager',
            # the constraint that the user is using (i.e. "> 1.0.0")
            'constraint': dependencyConstraint,
            # all available versions above and outside of their constraint
            # - usually you would need to use the package manager lib or API
            #   to get this information (we just fake it here)
            'available': [
                {'name': '2.0.0'},
            ],
        }

    return output


def collect_lockfile_dependencies(lockfile_data):
    """Convert the lockfile format to the dependencies schema"""
    output = {}

    for dependencyName, installedVersion in lockfile_data.items():
        output[dependencyName] = {
            'source': 'example-package-manager',
            'installed': {'name': installedVersion},
        }

    return output
