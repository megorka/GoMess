package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("b25d03a4749df33c94f7e34ac81a0a27219e512882cf44bcdd74cd2584b387f57d58a377b03859e398a5c28cc8b4816806bdfdf4addcd463cbc87b3bb21a6ee0c2b61f4bdc986baad395c8fb772d561c8201ccc64e98aa77e92d03a341dc2ad987ce3d751bcd2c1d5dca7b168092f77c9912ef83b964bf2dee9a83a51e760df2")

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateToken(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"user_id": strconv.Itoa(userID),
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//func ParseToken(tokenString string) (int, error) {
//	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		return jwtKey, nil
//	})
//	if err != nil {
//		return 0, err
//	}
//
//	claims, ok := token.Claims.(jwt.MapClaims)
//	if !ok || !token.Valid {
//		return 0, err
//	}
//
//	userID, ok := claims["user_id"].(int)
//	if !ok {
//		return 0, err
//	}
//
//	return userID, nil
//}
