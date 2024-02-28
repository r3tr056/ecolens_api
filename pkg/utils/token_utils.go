package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	Access  string
	Refresh string
}

func GenerateNewTokens(id uint) (*Tokens, error) {
	// Generate JWT Access Token
	accessToken, err := generateNewAccessToken(id)
	if err != nil {
		return nil, err
	}

	// Generate JWT Refresh Token
	refreshToken, err := generateNewRefreshToken()
	if err != nil {
		return nil, err
	}

	return &Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func generateNewAccessToken(id uint) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET_KEY"))

	minutesCount, _ := strconv.Atoi(os.Getenv("JWT_SECRET_KEY_EXPIRE_MIN_COUNT"))

	// create new claims
	claims := jwt.MapClaims{}

	// set public claims
	claims["id"] = id
	claims["issuer"] = id
	claims["expires"] = time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate token
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func generateNewRefreshToken() (string, error) {
	hash := sha256.New()

	refresh := os.Getenv("JWT_REFRESH_KEY") + time.Now().String()

	_, err := hash.Write([]byte(refresh))
	if err != nil {
		return "", err
	}

	hoursCount, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"))

	expireTime := fmt.Sprint(time.Now().Add(time.Hour * time.Duration(hoursCount)).Unix())
	t := hex.EncodeToString(hash.Sum(nil)) + "." + expireTime
	return t, nil
}
