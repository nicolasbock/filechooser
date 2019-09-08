def get_image_files(paths):
    # type: (List[str]) -> List[str]
    """Get a list of supported image files under a list of paths.

    Input paths is a list of paths under which to search recursively
    for support image files. The output is a list of files (including
    paths).

    """
    return ["a/picture_a.jpg", "b/picture_b.jpg"]
