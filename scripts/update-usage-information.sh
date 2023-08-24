#!/bin/bash

set -u -e

cat <<EOF > docs/source/supported-options.rst
Supported Options
=================

.. code-block:: console

EOF
pick-files --help 2>&1| sed --expression 's/^/   /' >> docs/source/supported-options.rst
