## Run Etcd

    docker pull quay.io/coreos/etcd
    
    docker run --rm --net=host quay.io/coreos/etcd

## Populate etcd
    
    docker pull debian
    
    etcdctl set /_coreos.com/fleet/machines/010/object  \
        '{"ID":"010","PublicIP":"127.0.0.1","Metadata":{"role":"services"},"Version":"0.11.7"}'
        
    etcdctl set /_coreos.com/fleet/machines/003/object  \
        '{"ID":"003","PublicIP":"172.17.0.254","Metadata":{"role":"nohere"},"Version":"0.11.7"}'

    etcdctl mkdir /_coreos.com/fleet/machines/004
    
    for i in {1..9}
    do etcdctl set /_coreos.com/fleet/machines/00${i}/object \
        "{\"ID\": \"00${i}\", \"PublicIP\": \"172.17.0.${i}\",\"Metadata\":{\"role\":\"docker\"},\"Version\":\"0.11.7\"}"
        docker run --rm -v $(pwd)/inventory:/inventory debian /inventory &
    done

                    
                    
## Golang

     dirname $(readlink $(which go))

     go get -u github.com/jteeuwen/go-bindata/...
     go get -u github.com/elazarl/go-bindata-assetfs/...
     
     
## Quick deploy


cat << EOF > inventory.service
[Service]
ExecStartPre=/usr/bin/curl -Lk ${BUCKET}/inventory/inventory -o /usr/bin/inventory
ExecStartPre=/bin/chmod +x /usr/bin/inventory
ExecStart=/usr/bin/inventory

[X-Fleet]
Global=true
EOF
fleetctl destroy inventory.service 
fleetctl start inventory.service



etcdctl set /traefik/backends/inventory/servers/server0/weight '1'
etcdctl set /traefik/backends/inventory/servers/server0/url 'http://127.0.0.1:5000'

etcdctl set /traefik/frontends/inventory/backend 'inventory'
etcdctl set /traefik/frontends/inventory/routes/inventory/rule 'PathPrefix:/'
