#!/bin/bash

set -u -e

: ${PICKFILES:=pick-files}

cat <<EOF > docs/source/supported-options.rst
Supported Options
=================

.. code-block:: console

EOF
${PICKFILES} --help 2>&1| sed --expression 's/^/   /' --expression 's/ \+$//' >> docs/source/supported-options.rst
