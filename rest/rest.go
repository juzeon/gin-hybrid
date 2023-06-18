package rest

import (
	"gin-hybrid/etclient"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
)

func Setup() {

}

type Service struct {
	Name      string
	Endpoints map[clientv3.LeaseID]string
	mu        sync.Mutex
}

func (s *Service) updateServiceDirectoryThread() {
	for {
		watchChan := etclient.WatchKeysByPrefix("list/" + s.Name)
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
	kvArr, err := etclient.GetRawKeysByPrefix("list/" + s.Name)
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
func NewService(name string) (*Service, error) {
	service := &Service{
		Name:      name,
		Endpoints: map[clientv3.LeaseID]string{},
		mu:        sync.Mutex{},
	}
	err := service.UpdateServiceDirectory()
	if err != nil {
		return nil, err
	}
	go service.updateServiceDirectoryThread()
	etclient.AddServiceRegisterEventListener(func() {
		err := service.UpdateServiceDirectory()
		if err != nil {
			log.Println("observer failed to update service directory of " + service.Name + ": " + err.Error())
		}
	})
	return service, nil
}
