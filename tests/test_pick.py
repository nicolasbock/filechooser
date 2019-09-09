import os
import pathlib
import shutil
import tempfile
import unittest

from filechooser import pick_file

images = [
    "collection_a/a/a_1.jpg",
    "collection_a/a/a_2.gif",
    "collection_a/b/b_1.jpg",
    "collection_a/b/b_2.txt",
    "collection_b/c/c_1.jpg",
]


class TestPick(unittest.TestCase):

    def setUp(self):
        # Create a filesystem structure in a temporary location. This
        # should eventually be replaced with something in memory, i.e.
        # using `pyfakefs`.
        self.fs_base = tempfile.mkdtemp()
        for image in images:
            path = os.path.dirname(image)
            try:
                os.makedirs(os.path.join(self.fs_base, path), exist_ok=True)
            except FileExistsError as e:
                print("Can not create path {}".format(
                    os.path.join(self.fs_base, path)))
                raise e
            pathlib.Path(os.path.join(self.fs_base, image)).touch()

    def tearDown(self):
        shutil.rmtree(self.fs_base)

    def test_get_image_files(self):
        # Sort the lists so we can compare them directly.
        image_files = sorted(pick_file.get_image_files(
            [os.path.join(self.fs_base, "collection_a"),
             os.path.join(self.fs_base, "collection_b")]))
        reference_image_files = sorted(
            [os.path.join(self.fs_base, p) for p in images if p.split(".")[-1] != "txt"])
        self.assertEqual(image_files, reference_image_files)

    def test_get_image_files_not_exist(self):
        with self.assertRaisesRegex(Exception, "does not exist"):
            pick_file.get_image_files(["does_not_exist"])
