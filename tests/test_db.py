import os
import tempfile
import unittest

import filechooser.timestamps as timestamps
from filechooser.logger import logger


class TestDB(unittest.TestCase):

    def setUp(self):
        from logging import DEBUG
        logger.setLevel(DEBUG)
        self.db = tempfile.NamedTemporaryFile(delete=False)
        self.db.close()
        timestamps.database = self.db.name
        timestamps.initialize_db()

    def tearDown(self):
        os.remove(self.db.name)

    def test_store_timestamp(self):
        timestamps.store_timestamp("a.gif", "timestamp1")
        result = timestamps.get_timestamp("a.gif")
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0]['filename'], "a.gif")

    def test_update_timestamp(self):
        timestamps.store_timestamp("b.gif", "timestamp1")
        timestamps.store_timestamp("b.gif", "timestamp2")
        result = timestamps.get_timestamp("b.gif")
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0]['timestamp'], "timestamp2")
