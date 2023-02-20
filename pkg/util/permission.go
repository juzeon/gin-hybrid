package util

import (
	"gin-hybrid/data/dto"
	"github.com/golang-jwt/jwt"
	"time"
)

const (
	JWTSecret = "fee516aa-8f5b-4248-a80c-bf467d40fd62"
)

func ParseJWT(token string) (*dto.UserClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &dto.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
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
	str, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		panic(err)
	}
	return str
}
