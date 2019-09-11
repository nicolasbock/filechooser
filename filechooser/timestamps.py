from tinydb import TinyDB, Query

from filechooser.logger import logger

database = "file-timestamps-db.json"
db = None
query = None


def initialize_db():
    """Initializes the database, creating it in case it didn't exist
    before.
    """
    global db, query
    logger.debug("using database file {}".format(database))
    if db is None:
        logger.debug("initializing new database")
        db = TinyDB(database)
        query = Query()


def store_timestamp(filename, timestamp):
    """Stores or updates the timestamp of filename.
    """
    global db, query
    logger.debug(
        "storing/updating timestamp {} {}".format(filename, timestamp))
    db.upsert({'filename': filename, 'timestamp': timestamp},
              query.filename == filename)


def get_timestamp(filename):
    """Get the timestamp of filename. If there is no timestamp in the
    database, return None.
    """
    global db, query
    return db.search(query.filename == filename)


def dump_db():
    """Prints out all records currently stored in the database.
    """
    global db
    return(db.all())
