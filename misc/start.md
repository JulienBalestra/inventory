## Populate etcd
    
    etcdctl set /_coreos.com/fleet/machines/001/object  \
        '{"ID":"001","PublicIP":"127.0.0.1","Metadata":{"role":"services"},"Version":"0.11.7"}'
            
    etcdctl set /_coreos.com/fleet/machines/002/object  \
        '{"ID":"002","PublicIP":"localhost","Metadata":{"role":"services"},"Version":"0.11.7"}'
                    
                    
                    
## Golang

     dirname $(readlink $(which go))
     
     
## Quick deploy


    cat << EOF > inventory.service
    [Service]
    ExecStartPre=/usr/bin/curl -Lk inventory/inventory -o /usr/bin/inventory
    ExecStartPre=/bin/chmod +x /usr/bin/inventory
    ExecStart=/usr/bin/inventory
    
    [X-Fleet]
    Global=true
    EOF