package util

import (
	"gin-hybrid/data/dto"
	"github.com/golang-jwt/jwt"
	"time"
)

var jwtSecret string

func ParseJWT(token string) (*dto.UserClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &dto.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*dto.UserClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func GenerateJWT(userID int, roleID int, roleName string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, dto.UserClaims{
		UserID:         userID,
		RoleID:         roleID,
		RoleName:       roleName,
		LoginTime:      time.Now(),
		StandardClaims: jwt.StandardClaims{},
	})
	str, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		panic(err)
	}
	return str
}
