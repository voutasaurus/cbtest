#!/bin/bash

# setup
set -ex
scriptdir="$(dirname "$0")"
cd $scriptdir

# main
docker build -t couchlocal .

# cleanup
set +ex
