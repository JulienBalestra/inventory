## Populate etcd
    
    etcdctl set /_coreos.com/fleet/machines/001/object  \
        '{"ID":"001","PublicIP":"127.0.0.1","Metadata":{"role":"services"},"Version":"0.11.7"}'
            
    etcdctl set /_coreos.com/fleet/machines/002/object  \
        '{"ID":"002","PublicIP":"localhost","Metadata":{"role":"services"},"Version":"0.11.7"}'
                    
                    
                    
## Golang

     dirname $(readlink $(which go))
     