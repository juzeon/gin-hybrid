package etclient

import (
	"context"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"strconv"
	"time"
)

const leaseTTL = 10

var client *clientv3.Client
var leaseID clientv3.LeaseID
var ErrNotExist = errors.New("error key not exist")
var namespace string

func Setup(endpoints []string, namespaceV string) error {
	namespace = namespaceV
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	client = c
	err = registerServiceOnce()
	if err != nil {
		return err
	}
	go keepAliveThread()
	return nil
}
func registerServiceOnce() error {
	resp, err := client.Lease.Grant(context.Background(), leaseTTL)
	if err != nil {
		return err
	}
	leaseID = resp.ID
	log.Println("register service successfully with lease id: " + strconv.Itoa(int(leaseID)))
	return nil
}
func keepAliveThread() {
	for {
		_, err := client.Lease.KeepAliveOnce(context.Background(), leaseID)
		if err != nil {
			log.Println("keep lease alive failed: " + err.Error())
			err0 := registerServiceOnce()
			if err0 != nil {
				log.Println("register service failed: " + err.Error())
				continue
			}
		}
		time.Sleep(leaseTTL / 2 * time.Second)
	}
}
func withKeyNamespace(key string) string {
	return "/" + namespace + "/" + key
}
