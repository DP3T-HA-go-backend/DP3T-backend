ETCDCTL_API=3 etcdctl --endpoints=https://10.0.26.10:2379,https://10.0.26.11:2379,https://10.0.26.13:2379 --cacert=/etc/ssl/etcd/ssl/ca.pem --cert=/etc/ssl/etcd/ssl/node-node2.pem --key=/etc/ssl/etcd/ssl/node-node2-key.pem $*
