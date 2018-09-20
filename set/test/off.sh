#!/bin/bash

# setup
set -ex
scriptdir="$(dirname "$0")"
cd $scriptdir

# main
docker-compose down
docker ps -aq | xargs docker rm -f;
docker network prune -f;

# cleanup
set +ex
