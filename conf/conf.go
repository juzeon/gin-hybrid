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

var InitConf Init
var ParentConf Parent

func LoadConfig(name string, target any) error {
	_, err := toml.DecodeFile("cmd/"+name+"/config.toml", &InitConf)
	if err != nil {
		return err
	}
	etclientConf := etclient.Conf{
		Endpoints: InitConf.Etcd.Endpoints,
		Namespace: InitConf.Etcd.Namespace,
		Name:      InitConf.Name,
		IP:        InitConf.IP,
		User:      InitConf.Etcd.User,
		Pass:      InitConf.Etcd.Pass,
		Port:      0,
	}
	err = etclient.Setup(etclientConf)
	if err != nil {
		return err
	}
	parentV, err := etclient.GetRawKey("parent_config")
	if err != nil && err != etclient.ErrNotExist {
		return err
	}
	err = toml.Unmarshal([]byte(parentV), &ParentConf)
	if err != nil {
		return err
	}
	configV, err := etclient.GetRawKey(name + "/config")
	if err != nil && err != etclient.ErrNotExist {
		return err
	}
	err = toml.Unmarshal([]byte(configV), target)
	if err != nil {
		return err
	}
	commonV := reflect.ValueOf(target).Elem().FieldByName("Common").Interface().(Common)
	etclientConf.Port = commonV.Port
	err = etclient.UpdateConf(etclientConf)
	if err != nil {
		return err
	}
	go watchConfigThread(target)
	return nil
}
func watchConfigThread(target any) {
	watchChan := etclient.WatchKey(InitConf.Name + "/config")
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			if event.Type == clientv3.EventTypeDelete {
				continue
			}
			configV := event.Kv.Value
			err := toml.Unmarshal(configV, target)
			if err != nil {
				log.Println("failed to unmarshal new config: " + err.Error())
				continue
			}
			log.Printf("updated new config: %#v", target)
		}
	}
}
