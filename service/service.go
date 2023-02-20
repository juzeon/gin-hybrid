package service

var ExUser *UserService

func Setup() {
	ExUser = NewUserService()
}
