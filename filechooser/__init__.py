import os
import subprocess


def _get_latest_tag():
    version = "unknown"
    version_file = ""
    if os.environ.get('SNAP_NAME', '') == 'pick-files':
        version_file = os.path.join(os.environ.get('SNAP', ''), 'filechooser')
    version_file = os.path.join(version_file, '.version')
    try:
        with open(version_file, encoding='utf8') as fd:
            version = fd.readline()
        return version
    except FileNotFoundError:
        pass

    git = subprocess.Popen(['git', 'describe', '--tags'],
                           stdout=subprocess.PIPE,
                           universal_newlines=True)
    git.wait()
    if git.returncode == 0:
        raw_version = git.stdout.readlines()[0][1:].rstrip().split('-')
        version = raw_version[0]
        if len(raw_version) > 1:
            version += ".dev" + raw_version[1]
    with open('.version', mode='w', encoding='utf-8') as fd:
        fd.write('{}\n'.format(version))
    return version


__version__ = _get_latest_tag()
