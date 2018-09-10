#!/bin/bash

# setup
set -ex
scriptdir="$(dirname "$0")"

# prerequisites
rootdir=$scriptdir/../../
./$rootdir/cmd/hello/dockerbuild.sh

# main
set -ex
cd $scriptdir
docker-compose up

# cleanup
set +ex
