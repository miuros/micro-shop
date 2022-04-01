package util

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type Claim struct {
	jwt.StandardClaims
	Username string
	RoleName string
	UserUuid string
}

func ReleaseToken(userName, userUuid string, roleName string) (string, error) {
	var claim = Claim{
		Username: userName,
		UserUuid: userUuid,
		RoleName: roleName,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "micro-shop",
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
		},
	}

	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	secretKey := viper.GetString("secretKey")
	token, err := tokenClaim.SignedString([]byte(secretKey))
	if err != nil {
		//zap.L().Info(err.Error())
		return "", err
	}
	return token, nil

}

func ParseToken(token string) (*Claim, error) {
	tokenClaim, err := jwt.ParseWithClaims(token[7:], &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("secretKey")), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaim.Valid {
		claim, ok := tokenClaim.Claims.(*Claim)
		if ok {
			return claim, nil
		}
	}
	return nil, errors.New("fail to parse token")

}
