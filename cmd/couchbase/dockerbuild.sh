#!/bin/bash

# setup
set -ex
scriptdir="$(dirname "$0")"
cd $scriptdir

# main
docker run -d --name tmpcouch -v "$(pwd)"/tmpconfig:/opt/couchbase/var couchbase
docker cp config.sh tmpcouch:/config.sh
docker exec -t \
	-e USER=admin \
	-e PASS=password \
	-e BUCKET=testbucket \
        tmpcouch /bin/bash config.sh
docker kill tmpcouch
docker rm tmpcouch
docker build -t couchlocal .

# cleanup
rm -r tmpconfig
set +ex
