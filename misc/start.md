## Populate etcd
    
    etcdctl set /_coreos.com/fleet/machines/001/object  \
        '{"ID":"001","PublicIP":"127.0.0.1","Metadata":{"role":"services"},"Version":"0.11.7"}'
            
    etcdctl set /_coreos.com/fleet/machines/002/object  \
        '{"ID":"002","PublicIP":"localhost","Metadata":{"role":"services"},"Version":"0.11.7"}'
            
    etcdctl set /_coreos.com/fleet/machines/003/object  \
        '{"ID":"003","PublicIP":"172.17.0.254","Metadata":{"role":"nohere"},"Version":"0.11.7"}'
                    
                    
                    
## Golang

     dirname $(readlink $(which go))
     
     
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
