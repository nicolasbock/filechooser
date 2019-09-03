import unittest
from filechooser import pick_file


class TestPick(unittest.TestCase):

    def setup(self):
        pass

    def tearDown(self):
        pass

    def test_pick_file(self):
        self.assertEqual(pick_file.pick_file(), "picture_a.jpg")
        assert(True)
