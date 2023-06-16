package etclient

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func GetRawKey(key string, opts ...clientv3.OpOption) (string, error) {
	resp, err := client.Get(context.Background(), withKeyNamespace(key), opts...)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", ErrNotExist
	}
	return string(resp.Kvs[0].Value), nil
}
func PutRawKey(key string, value string, opts ...clientv3.OpOption) error {
	_, err := client.Put(context.Background(), withKeyNamespace(key), value, opts...)
	return err
}
func WatchKey(key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	return client.Watch(context.Background(), withKeyNamespace(key), opts...)
}
