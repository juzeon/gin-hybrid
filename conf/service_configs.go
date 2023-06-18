package conf

type Gateway struct {
	Common
}

var GatewayServiceConfig *ServiceConfig[Gateway]

type User struct {
	Common
}

var UserServiceConfig *ServiceConfig[User]

type Article struct {
	Common
}

var ArticleServiceConfig *ServiceConfig[Article]
