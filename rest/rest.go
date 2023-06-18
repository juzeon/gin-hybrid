package rest

import (
	"gin-hybrid/etclient"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
)

type Client struct {
	etclientInst *etclient.Client
	services     []*Service
}

func NewClient(etclientInst *etclient.Client) *Client {
	return &Client{etclientInst: etclientInst}
}
func (c *Client) GetService(name string) *Service {
	for _, service := range c.services {
		if service.Name == name {
			return service
		}
	}
	return nil
}
func (c *Client) AddService(name string) (*Service, error) {
	service := &Service{
		Name:         name,
		Endpoints:    map[clientv3.LeaseID]string{},
		mu:           sync.Mutex{},
		etclientInst: c.etclientInst,
	}
	err := service.UpdateServiceDirectory()
	if err != nil {
		return nil, err
	}
	go service.updateServiceDirectoryThread()
	c.etclientInst.AddServiceRegisterEventListener(func() {
		err := service.UpdateServiceDirectory()
		if err != nil {
			log.Println("observer failed to update service directory of " + service.Name + ": " + err.Error())
		}
	})
	c.services = append(c.services, service)
	return service, nil
}

type Service struct {
	Name         string
	Endpoints    map[clientv3.LeaseID]string
	mu           sync.Mutex
	etclientInst *etclient.Client
}

func (s *Service) updateServiceDirectoryThread() {
	for {
		watchChan := s.etclientInst.WatchKeysByPrefix("list/" + s.Name)
		for watchResp := range watchChan {
			s.mu.Lock()
			for _, event := range watchResp.Events {
				leaseID := etclient.ConvertStringToLeaseID(etclient.GetKeyLastSegment(string(event.Kv.Key)))
				switch event.Type {
				case clientv3.EventTypePut:
					s.Endpoints[leaseID] = string(event.Kv.Value)
					log.Printf("update %v of service %v\n", leaseID, s.Name)
				case clientv3.EventTypeDelete:
					delete(s.Endpoints, leaseID)
					log.Printf("delete %v of service %v\n", leaseID, s.Name)
				}
			}
			s.mu.Unlock()
		}
	}
}
func (s *Service) UpdateServiceDirectory() error {
	kvArr, err := s.etclientInst.GetRawKeysByPrefix("list/" + s.Name)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Endpoints = map[clientv3.LeaseID]string{}
	for _, kv := range kvArr {
		if string(kv.Value) == "" {
			continue
		}
		leaseID := etclient.ConvertStringToLeaseID(etclient.GetKeyLastSegment(string(kv.Key)))
		s.Endpoints[leaseID] = string(kv.Value)
	}
	log.Printf("updated the directory of service "+s.Name+": %#v\n", s.Endpoints)
	return nil
}
