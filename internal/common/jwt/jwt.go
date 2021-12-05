package jwt

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/wyy-go/go-web-template/internal/common/errors"
	"time"
)

type Account struct {
	Uid        int64
	DeviceName string
	Platform   string
}

type UserClaims struct {
	jwt.StandardClaims
	Acc *Account
}

func Encode(privateKey *rsa.PrivateKey, userInfo *Account, expiresIn int) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second).Unix(),
		},
		Acc: userInfo,
	}
	if expiresIn == 0 {
		claims.ExpiresAt = 0
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func Decode(publicKey *rsa.PublicKey, tokenString string) (*Account, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		},
	)
	if err == nil {
		if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
			// 验证成功，返回信息
			return claims.Acc, nil
		}
	}

	ve, ok := err.(*jwt.ValidationError)
	if !ok || ve.Errors != jwt.ValidationErrorExpired {
		return nil, errors.ErrInvalidToken
	} else {
		return nil, errors.ErrTokenExpired
	}
	// 验证失败
	//return nil, err
}
