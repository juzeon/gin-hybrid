package conf

import (
	"gin-hybrid/etclient"
	"github.com/BurntSushi/toml"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

type Init struct {
	Endpoints []string `toml:"endpoints"`
	Namespace string   `toml:"namespace"`
	Name      string   `toml:"name"`
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

var InitConf Init
var ParentConf Parent

func LoadConfig(name string, target any) error {
	_, err := toml.DecodeFile("cmd/"+name+"/config.toml", &InitConf)
	if err != nil {
		return err
	}
	err = etclient.Setup(InitConf.Endpoints, InitConf.Namespace)
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
