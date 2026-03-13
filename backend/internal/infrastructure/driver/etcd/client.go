package etcd

import clientv3 "go.etcd.io/etcd/client/v3"

func NewClient(cfg *clientv3.Config) *clientv3.Client {
	client, err := clientv3.New(*cfg)
	if err != nil {
		panic(err)
	}
	return client
}