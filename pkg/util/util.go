package util

import "golang.org/x/crypto/bcrypt"

type SetupConf struct {
	JWTSecret string
}

func Setup(conf SetupConf) {
	jwtSecret = conf.JWTSecret
}
func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}
func ValidatePassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
