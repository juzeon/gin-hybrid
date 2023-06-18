package etclient

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strconv"
	"strings"
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
func GetRawKeysByPrefix(prefix string, opts ...clientv3.OpOption) ([]*mvccpb.KeyValue, error) {
	opts = append(opts, clientv3.WithPrefix())
	resp, err := client.Get(context.Background(), withKeyNamespace(prefix), opts...)
	if err != nil {
		return nil, err
	}
	return resp.Kvs, nil
}
func PutRawKey(key string, value string, opts ...clientv3.OpOption) error {
	_, err := client.Put(context.Background(), withKeyNamespace(key), value, opts...)
	return err
}
func WatchKey(key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	return client.Watch(context.Background(), withKeyNamespace(key), opts...)
}
func WatchKeysByPrefix(prefix string, opts ...clientv3.OpOption) clientv3.WatchChan {
	opts = append(opts, clientv3.WithPrefix())
	return client.Watch(context.Background(), withKeyNamespace(prefix), opts...)
}
func withKeyNamespace(key string) string {
	return "/" + conf.Namespace + "/" + key
}
func GetKeyLastSegment(fullKey string) string {
	arr := strings.Split(fullKey, "/")
	return arr[len(arr)-1]
}
func ConvertStringToLeaseID(str string) clientv3.LeaseID {
	i, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return clientv3.LeaseID(i)
}
