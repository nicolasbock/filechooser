import os
import tempfile
import unittest

import filechooser.db as db
from filechooser.logger import logger


class TestDB(unittest.TestCase):

    def setUp(self):
        from logging import DEBUG
        logger.setLevel(DEBUG)
        self.db = tempfile.NamedTemporaryFile(delete=False)
        self.db.close()
        db.database = self.db.name
        db.initialize_db()

    def tearDown(self):
        os.remove(self.db.name)

    def test_set_timestamp(self):
        db.set_timestamp("a.gif", "timestamp1")
        result = db.get_timestamp("a.gif")
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0]['filename'], "a.gif")

    def test_update_timestamp(self):
        db.set_timestamp("b.gif", "timestamp1")
        db.set_timestamp("b.gif", "timestamp2")
        result = db.get_timestamp("b.gif")
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0]['timestamp'], "timestamp2")
