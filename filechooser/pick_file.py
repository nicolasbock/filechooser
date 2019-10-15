import mimetypes
import os.path

try:
    from typing import List
except ImportError:
    pass

from filechooser.db import set_timestamp
from filechooser.logger import logger

image_file_extensions = [
    "image/jpeg",
    "image/gif",
    "image/png",
]


def get_image_files(paths):
    # type: (List[str]) -> List[str]
    """Get a list of supported image files under a list of paths.

    Input paths is a list of paths under which to search recursively
    for support image files. The output is a list of files (including
    paths).

    """

    result = []  # type: List[str]
    for path in paths:
        logger.debug("Testing path {}".format(path))
        if not os.path.exists(path):
            raise Exception("The path {} does not exist".format(path))
        if os.path.isdir(path):
            new_paths = [os.path.join(path, p) for p in os.listdir(
                path) if p not in ['.', '..', path]]
            result += get_image_files(new_paths)
        else:
            filetype = mimetypes.guess_type("file://" + path)
            if filetype[0] in image_file_extensions:
                result.append(path)
    return result


def set_image_timestamp(filename):
    # type: (str) -> None
    """Set the timestamp of a file, i.e. when it was chosen.

    The timestamp for "filename" is set to the current time.
    """
    set_timestamp(filename)
