package rest

import (
	"errors"
	"fmt"
	"gin-hybrid/etclient"
	"gin-hybrid/pkg/app"
	"github.com/go-resty/resty/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Client struct {
	etclientInst *etclient.Client
	httpClient   *resty.Client
	services     []*Service
}

func NewClient(etclientInst *etclient.Client) *Client {
	httpClient := resty.New().SetTimeout(10 * time.Second)
	return &Client{etclientInst: etclientInst, httpClient: httpClient}
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
		httpClient:   c.httpClient,
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
	httpClient   *resty.Client
}

func (s *Service) Call(method string, path string, data any) (any, error) {
	method = strings.ToUpper(method)
	endpoint, err := s.GetEndpointRandomly()
	if err != nil {
		return "", err
	}
	req := s.httpClient.R().SetResult(&app.Result{})
	req.Method = method
	req.URL = "http://" + endpoint + "/api/" + path
	dataValue := reflect.ValueOf(data).Elem()
	values := map[string]string{}
	switch dataValue.Kind() {
	case reflect.Map:
		valuesTmp := dataValue.Interface().(map[string]any)
		for k, v := range valuesTmp {
			values[k] = fmt.Sprintf("%v", v)
		}
	case reflect.Struct:
		values = s.convertStructToMap(data)
	}
	if method == "GET" || method == "HEAD" {
		req.SetQueryParams(values)
	} else {
		req.SetFormData(values)
	}
	resp, err := req.Send()
	if err != nil {
		return "", errors.New("failed to call remote api: " + err.Error())
	}
	result := resp.Result().(*app.Result)
	if result.Code > 399 {
		return "", errors.New("failed to call remote api: " + result.Msg)
	}
	return result.Data, nil
}
func (s *Service) convertStructToMap(data any) map[string]string {
	// Create an empty map to store the result
	result := make(map[string]string)
	// Get the value and type of the struct
	v := reflect.ValueOf(data).Elem()
	t := v.Type()
	// Loop over the fields of the struct
	for i := 0; i < t.NumField(); i++ {
		// Get the field value and type
		fv := v.Field(i)
		ft := t.Field(i)
		// Get the form tag of the field
		tag := ft.Tag.Get("form")
		// If the tag is empty, use the field name as the key
		if tag == "" {
			tag = ft.Name
		}
		// Convert the field value to an interface{}
		fvi := fmt.Sprintf("%v", fv.Interface())
		// Store the key-value pair in the result map
		result[tag] = fvi
	}
	return result
}
func (s *Service) GetEndpointRandomly() (string, error) {
	if len(s.Endpoints) == 0 {
		return "", errors.New("no available service")
	}
	i := rand.Intn(len(s.Endpoints))
	for _, v := range s.Endpoints {
		if i == 0 {
			return v, nil
		}
		i--
	}
	return "", nil
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
