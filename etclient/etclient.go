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

type Conf struct {
	Endpoints []string
	Namespace string
	Name      string
	IP        string
	Port      int
}

var conf Conf

func Setup(confV Conf) error {
	conf = confV
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.Endpoints,
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
func UpdateConf(confV Conf) error {
	conf = confV
	err := updateListDirectory()
	return err
}
func registerServiceOnce() error {
	resp, err := client.Lease.Grant(context.Background(), leaseTTL)
	if err != nil {
		return err
	}
	leaseID = resp.ID
	err = updateListDirectory()
	if err != nil {
		return err
	}
	log.Println("registered service successfully with lease id: " + strconv.Itoa(int(leaseID)))
	return nil
}
func updateListDirectory() error {
	if conf.Port == 0 {
		return nil
	}
	value := conf.IP + ":" + strconv.Itoa(conf.Port)
	err := PutRawKey("list/"+conf.Name+"/"+strconv.Itoa(int(leaseID)), value, clientv3.WithLease(leaseID))
	if err != nil {
		return err
	}
	log.Println("updated list directory: " + value)
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
	return "/" + conf.Namespace + "/" + key
}
