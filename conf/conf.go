package conf

import (
	"gin-hybrid/etclient"
	"github.com/BurntSushi/toml"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"reflect"
)

type Init struct {
	Etcd Etcd   `toml:"etcd"`
	Name string `toml:"name"`
	IP   string `toml:"ip"`
}
type Etcd struct {
	Endpoints []string `toml:"endpoints"`
	Namespace string   `toml:"namespace"`
	User      string   `toml:"user"`
	Pass      string   `toml:"pass"`
}

type Parent struct {
	DB ParentDB `toml:"db"`
}
type ParentDB struct {
	Driver string `toml:"driver"`
	Host   string `toml:"host"`
	Port   int    `toml:"port"`
	User   string `toml:"user"`
	Pass   string `toml:"pass"`
	DB     string `toml:"db"`
}
type Common struct {
	Port int `toml:"port"`
}
type ServiceConfig[T any] struct {
	Name       string
	InitConf   Init
	ParentConf Parent
	SelfConf   T
	Etclient   *etclient.Client
}

func NewServiceConfig[T any](name string) (*ServiceConfig[T], error) {
	srvConf := &ServiceConfig[T]{Name: name}
	etclientIns, err := loadConfig(srvConf)
	if err != nil {
		return nil, err
	}
	srvConf.Etclient = etclientIns
	return srvConf, nil
}

func loadConfig[T any](config *ServiceConfig[T]) (*etclient.Client, error) {
	// load local config
	_, err := toml.DecodeFile("cmd/"+config.Name+"/config.toml", &config.InitConf)
	if err != nil {
		return nil, err
	}
	// initialize etclient using local config
	etclientConf := etclient.Conf{
		Endpoints: config.InitConf.Etcd.Endpoints,
		Namespace: config.InitConf.Etcd.Namespace,
		Name:      config.InitConf.Name,
		IP:        config.InitConf.IP,
		User:      config.InitConf.Etcd.User,
		Pass:      config.InitConf.Etcd.Pass,
		Port:      0, // not available for now
	}
	etclientIns, err := etclient.NewClient(etclientConf)
	if err != nil {
		return nil, err
	}
	parentV, err := etclientIns.GetRawKey("parent_config")
	if err != nil && err != etclient.ErrNotExist {
		return nil, err
	}
	err = toml.Unmarshal([]byte(parentV), &config.ParentConf)
	if err != nil {
		return nil, err
	}
	// initialize config for current service
	configV, err := etclientIns.GetRawKey(config.Name + "/config")
	if err != nil && err != etclient.ErrNotExist {
		return nil, err
	}
	err = toml.Unmarshal([]byte(configV), &config.SelfConf)
	if err != nil {
		return nil, err
	}
	commonV := reflect.ValueOf(&config.SelfConf).Elem().FieldByName("Common").Interface().(Common)
	etclientConf.Port = commonV.Port
	err = etclientIns.RegisterService(etclientConf)
	if err != nil {
		return nil, err
	}
	go watchConfigThread(config)
	return etclientIns, nil
}

// watchConfigThread watches config changes and update
func watchConfigThread[T any](config *ServiceConfig[T]) {
	for {
		watchChan := config.Etclient.WatchKey(config.InitConf.Name + "/config")
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				if event.Type == clientv3.EventTypeDelete {
					continue
				}
				configV := event.Kv.Value
				err := toml.Unmarshal(configV, &config.SelfConf)
				if err != nil {
					log.Println("failed to unmarshal new config: " + err.Error())
					continue
				}
				log.Printf("updated new config: %#v\n", &config.SelfConf)
			}
		}
	}
}
