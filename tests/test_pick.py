import unittest
from filechooser import pick_file


class TestPick(unittest.TestCase):

    def setup(self):
        pass

    def tearDown(self):
        pass

    def test_get_image_files(self):
        self.assertEqual(pick_file.get_image_files(
            ["a", "b"]), ["a/picture_a.jpg", "b/picture_b.jpg"])
        assert(True)
