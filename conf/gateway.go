package conf

type Gateway struct {
	Port int `toml:"port"`
}

var GatewayConf Gateway

func (g *Gateway) Load() {
	err := LoadConfig("gateway", g)
	if err != nil {
		panic(err)
	}
}
