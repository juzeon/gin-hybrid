package etclient

import (
	"context"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"strconv"
	"sync"
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
	User      string
	Pass      string
	Port      int
}

type serviceRegisterEventObserver struct {
	Listeners []func()
	Mu        sync.Mutex
}

func (s *serviceRegisterEventObserver) NotifyAll() {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for _, listener := range s.Listeners {
		listener()
	}
}
func (s *serviceRegisterEventObserver) AddListener(f func()) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Listeners = append(s.Listeners, f)
}

var conf Conf
var serviceRegisterEventObservers serviceRegisterEventObserver

func Setup(confV Conf) error {
	conf = confV
	// create an etcd client; port is not available yet for now
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.Endpoints,
		DialTimeout: 5 * time.Second,
		Username:    conf.User,
		Password:    conf.Pass,
	})
	if err != nil {
		return err
	}
	client = c
	return nil
}
func RegisterService(confV Conf) error {
	conf = confV
	err := registerServiceOnce()
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
func AddServiceRegisterEventListener(f func()) {
	serviceRegisterEventObservers.AddListener(f)
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
	go serviceRegisterEventObservers.NotifyAll()
	return nil
}

// updateListDirectory register current service into the etcd directory
func updateListDirectory() error {
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
