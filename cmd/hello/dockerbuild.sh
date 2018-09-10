#!/bin/bash

# setup
set -ex
scriptdir="$(dirname "$0")"
cd $scriptdir

# main
GOOS=linux go build -o hello
docker build -t hello .

# cleanup
rm hello
set +ex
