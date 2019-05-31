import json
from subprocess import run

from utils import mock_lockfile_update, mock_manifest_update, get_lockfile_fingerprint, write_json_to_temp_file
from collect import collect_lockfile_dependencies


def act():
    # An actor will always be given a set of "input" data, so that it knows what
    # exactly it is supposed to update. That JSON data will be stored in a file
    # at /dependencies/input_data.json for you to load.
    with open('/dependencies/input_data.json', 'r') as f:
        data = json.load(f)

    # Start a new branch for this update
    run(['deps', 'component', 'branch'], check=True)

    for lockfile_path, lockfile_data in data.get('lockfiles', {}).items():
        # If "lockfiles" are present then it means that there are updates to
        # those lockfiles that you can make. The most basic way to handle this
        # is to use whatever "update" command is provided by the package
        # manager, and then commit and push the entire update. You can try to be
        # more granular than that if you want, but performing the entire "update"
        # at once is an easier place to start.

        # 1) Do the lockfile update
        #    Since lockfile can change frequently, you'll want to "collect" the
        #    exact update that you end up making, in case it changed slightly from
        #    the original update that it was asked to make.
        updated_lockfile_data = mock_lockfile_update(lockfile_path)
        lockfile_data['updated']['dependencies'] = collect_lockfile_dependencies(updated_lockfile_data)
        lockfile_data['updated']['fingerprint'] = get_lockfile_fingerprint(lockfile_path)

        # 2) Add and commit the changes
        run(['deps', 'component', 'commit', '-m', 'Update ' + lockfile_path, lockfile_path], check=True)

    for manifest_path, manifest_data in data.get('manifests', {}).items():
        for dependency_name, updated_dependency_data in manifest_data['updated']['dependencies'].items():
            installed = manifest_data['current']['dependencies'][dependency_name]['constraint']
            version_to_update_to = updated_dependency_data['constraint']
            mock_manifest_update(manifest_path, dependency_name, version_to_update_to)

            run(['deps', 'component', 'commit', '-m', 'Update {} from {} to {}'.format(dependency_name, installed, version_to_update_to), manifest_path], check=True)

    # Shell out to `pullrequest` to make the actual pull request.
    #    It will automatically use the existing env variables and JSON schema
    #    to submit a pull request, or simulate one a test mode.
    run(['deps', 'component', 'pullrequest', write_json_to_temp_file(data)], check=True)
