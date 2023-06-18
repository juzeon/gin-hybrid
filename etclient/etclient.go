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

var ErrNotExist = errors.New("error key not exist")

type Client struct {
	client                        *clientv3.Client
	leaseID                       clientv3.LeaseID
	conf                          Conf
	serviceRegisterEventObservers serviceRegisterEventObserver
}

func NewClient(confV Conf) (*Client, error) {
	client := &Client{conf: confV}
	// create an etcd client; port is not available yet for now
	clientV, err := clientv3.New(clientv3.Config{
		Endpoints:   client.conf.Endpoints,
		DialTimeout: 5 * time.Second,
		Username:    client.conf.User,
		Password:    client.conf.Pass,
	})
	if err != nil {
		return nil, err
	}
	client.client = clientV
	return client, nil
}

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

func (c *Client) RegisterService(confV Conf) error {
	c.conf = confV
	err := c.registerServiceOnce()
	if err != nil {
		return err
	}
	go c.keepAliveThread()
	return nil
}
func (c *Client) UpdateConf(confV Conf) error {
	c.conf = confV
	err := c.updateListDirectory()
	return err
}
func (c *Client) AddServiceRegisterEventListener(f func()) {
	c.serviceRegisterEventObservers.AddListener(f)
}
func (c *Client) registerServiceOnce() error {
	resp, err := c.client.Lease.Grant(context.Background(), leaseTTL)
	if err != nil {
		return err
	}
	c.leaseID = resp.ID
	err = c.updateListDirectory()
	if err != nil {
		return err
	}
	log.Println("registered service successfully with lease id: " + strconv.Itoa(int(c.leaseID)))
	go c.serviceRegisterEventObservers.NotifyAll()
	return nil
}

// updateListDirectory register current service into the etcd directory
func (c *Client) updateListDirectory() error {
	value := c.conf.IP + ":" + strconv.Itoa(c.conf.Port)
	err := c.PutRawKey("list/"+c.conf.Name+"/"+strconv.Itoa(int(c.leaseID)), value, clientv3.WithLease(c.leaseID))
	if err != nil {
		return err
	}
	log.Println("updated list directory: " + value)
	return nil
}
func (c *Client) keepAliveThread() {
	for {
		_, err := c.client.Lease.KeepAliveOnce(context.Background(), c.leaseID)
		if err != nil {
			log.Println("keep lease alive failed: " + err.Error())
			err0 := c.registerServiceOnce()
			if err0 != nil {
				log.Println("register service failed: " + err.Error())
				continue
			}
		}
		time.Sleep(leaseTTL / 2 * time.Second)
	}
}
