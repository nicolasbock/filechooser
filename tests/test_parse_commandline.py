import unittest

try:
    from unittest import mock
except ImportError:
    import mock

import sys

from filechooser import parse_commandline


class TestParseCommandline(unittest.TestCase):

    def test_empty_dirs(self):
        with mock.patch.object(sys, "argv", ["/script"]):
            with self.assertRaises(SystemExit):
                parse_commandline.parse_commandline()

    def test_dirs(self):
        testdirs = ["/images-1", "/images-2"]
        with mock.patch.object(sys, "argv", ["/scriptpath"] + testdirs):
            options = parse_commandline.parse_commandline()
        self.assertTrue(isinstance(options.DIR, list))
        self.assertTrue(len(options.DIR) == 2)
        self.assertEqual(options.DIR, testdirs)
