import os
import hashlib
import tempfile
import json


def write_json_to_temp_file(data):
    """Writes JSON data to a temporary file and returns the path to it"""
    fp = tempfile.NamedTemporaryFile(delete=False)
    fp.write(json.dumps(data).encode('utf-8'))
    fp.close()
    return fp.name


def mock_manifest_update(manifest_path, dependency_name, version_to_update_to):
    with open(manifest_path, 'r') as f:
        manifest_data = json.load(f)

    manifest_data[dependency_name] = version_to_update_to

    with open(manifest_path, 'w+') as f:
        f.write(json.dumps(manifest_data, indent=4))

    return manifest_data


def mock_lockfile_update(path):
    """
    This is a mock update. In place of this, you might simply shell out
    to a command like `yarn upgrade`.
    """
    updated_lockfile_contents = {
        'package1': '1.2.0'
    }
    with open(path, 'w+') as f:
        f.write(json.dumps(updated_lockfile_contents, indent=4))
    return updated_lockfile_contents


def get_lockfile_fingerprint(path):
    return hashlib.md5(open(path, 'r').read().encode('utf-8')).hexdigest()


def print_settings_example():
    """
    You can use settings to get additional information from the user via their
    dependencies.io configuration file. Settings will be automatically injected as
    env variables with the "SETTING_" prefix.

    All settings will be passed as strings. More complex types will be json
    encoded. You should always provide defaults, if possible.
    """
    SETTING_EXAMPLE_LIST = json.loads(os.getenv('SETTING_EXAMPLE_LIST', '[]'))
    SETTING_EXAMPLE_STRING = os.getenv('SETTING_EXAMPLE_STRING', 'default')

    print('List setting values: {}'.format(SETTING_EXAMPLE_LIST))
    print('String setting value: {}'.format(SETTING_EXAMPLE_STRING))
