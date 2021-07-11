import os.path
import subprocess


def _get_latest_tag():
    version = "unknown"
    try:
        with open('.version') as fd:
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
    with open('.version', 'w') as fd:
        fd.write('{}\n'.format(version))
    return version


__version__ = _get_latest_tag()
