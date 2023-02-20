package dto

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type UserClaims struct {
	jwt.StandardClaims
	UserID    int       `json:"user_id"`
	RoleID    int       `json:"role_id"`
	RoleName  string    `json:"role_name"`
	LoginTime time.Time `json:"login_time"`
}
