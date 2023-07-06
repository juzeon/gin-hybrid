package etclient

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (c *Client) GetRawKey(key string, opts ...clientv3.OpOption) (string, error) {
	resp, err := c.client.Get(context.Background(), c.withKeyNamespace(key), opts...)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", ErrNotExist
	}
	return string(resp.Kvs[0].Value), nil
}
func (c *Client) GetRawKeysByPrefix(prefix string, opts ...clientv3.OpOption) ([]*mvccpb.KeyValue, error) {
	opts = append(opts, clientv3.WithPrefix())
	resp, err := c.client.Get(context.Background(), c.withKeyNamespace(prefix), opts...)
	if err != nil {
		return nil, err
	}
	return resp.Kvs, nil
}
func (c *Client) PutRawKey(key string, value string, opts ...clientv3.OpOption) error {
	_, err := c.client.Put(context.Background(), c.withKeyNamespace(key), value, opts...)
	return err
}
func (c *Client) PutRawKeyWithTTL(key string, value string, ttl int, opts ...clientv3.OpOption) error {
	leaseResp, err := c.client.Lease.Grant(context.Background(), int64(ttl))
	if err != nil {
		return err
	}
	opts = append(opts, clientv3.WithLease(leaseResp.ID))
	return c.PutRawKey(key, value, opts...)
}
func (c *Client) WatchKey(key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	return c.client.Watch(context.Background(), c.withKeyNamespace(key), opts...)
}
func (c *Client) WatchKeysByPrefix(prefix string, opts ...clientv3.OpOption) clientv3.WatchChan {
	opts = append(opts, clientv3.WithPrefix())
	return c.client.Watch(context.Background(), c.withKeyNamespace(prefix), opts...)
}
func (c *Client) withKeyNamespace(key string) string {
	return "/" + c.conf.Namespace + "/" + key
}
