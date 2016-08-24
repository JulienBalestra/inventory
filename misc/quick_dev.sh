#!/usr/bin/env bash

make re

if [ ! -f inventory ]
then
    echo "inventory not here"
    exit 1
fi

docker run --name etcd --rm --net=host quay.io/coreos/etcd &
until nc -zvv 127.0.0.1 2379 ; do sleep 1 ; done

for i in {2..9}
do
    etcdctl set /_coreos.com/fleet/machines/00${i}/object \
    "{\"ID\": \"00${i}\", \"PublicIP\": \"172.17.0.${i}\",\"Metadata\":{\"role\":\"docker\"},\"Version\":\"0.11.7\"}"
    docker run -e ETCD_ADDRESS="http://172.17.0.1:2379" --name inventory${i} --rm -v $(pwd)/inventory:/inventory debian /inventory > /dev/null &
done

function clean_docker {
    set -x
    docker kill etcd
    for i in {2..9}
    do
        docker kill inventory${i}
    done
}

trap clean_docker 2

HTTP_SERVE=fs ./inventory

wait