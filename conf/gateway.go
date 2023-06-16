package conf

type Gateway struct {
	Common
}

var GatewayConf Gateway

func (g *Gateway) Load() {
	err := LoadConfig("gateway", g)
	if err != nil {
		panic(err)
	}
}
