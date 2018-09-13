#!/bin/bash

# adapted from https://github.com/cha87de/couchbase-docker-cloudnative

function bucketCreate(){
    couchbase-cli bucket-create -c localhost -u Administrator -p password \
        --bucket=$BUCKET \
        --bucket-type=couchbase \
        --bucket-ramsize=512 \
        --bucket-replica=1 \
        --wait
    if [[ $? != 0 ]]; then
        return 1
    fi
}

function userCreate(){
    createOutput=$(couchbase-cli user-manage -c localhost -u Administrator -p password \
     --set --rbac-username $USER --rbac-password $PASS \
     --roles admin --auth-domain local)
    if [[ $? != 0 ]]; then
        echo $createOutput >&2
        return 1
    fi
}

function indexCreate(){
    cmd='CREATE PRIMARY INDEX ON `'$BUCKET'`'
    createOutput=$(cbq -u $USER -p $PASS --script="$cmd")
    if [[ $? != 0 ]]; then
        echo $createOutput >&2
        return 1
    fi
}

function clusterUp(){
    # wait for service to come up
    until $(curl --output /dev/null --silent --head --fail http://localhost:8091); do
        printf '.'
        sleep 1
    done

    # initialize cluster
    initOutput=$(couchbase-cli cluster-init -c localhost \
            --cluster-username=Administrator \
            --cluster-password=password \
            --cluster-port=8091 \
            --services=data,index,query,fts \
            --cluster-ramsize=1024 \
            --cluster-index-ramsize=256 \
            --cluster-fts-ramsize=256 \
            --index-storage-setting=default)
    if [[ $? != 0 ]]; then
        echo $initOutput >&2
        return 1
    fi
}

function main(){
    echo "Couchbase UI :8091"
    echo "Couchbase logs /opt/couchbase/var/lib/couchbase/logs"
    exec /usr/sbin/runsvdir-start &
    if [[ $? != 0 ]]; then
        echo "Couchbase startup failed. Exiting." >&2
        exit 1
    fi

    clusterUp
    if [[ $? != 0 ]]; then
        echo "Cluster init failed. Exiting." >&2
        exit 1
    fi

    bucketCreate
    if [[ $? != 0 ]]; then
        echo "Bucket create failed. Exiting." >&2
        exit 1
    fi

    userCreate
    if [[ $? != 0 ]]; then
        echo "User create failed. Exiting." >&2
        exit 1
    fi

    indexCreate
    if [[ $? != 0 ]]; then
        echo "Index create failed. Exiting." >&2
        exit 1
    fi
}

set -ex
main
set +ex
