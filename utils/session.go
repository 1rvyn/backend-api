package utils

import (
	"authserver/database"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"strings"
)

var SecretKey = os.Getenv("JWT_SECRET")

func GetSession(c *fiber.Ctx) (map[string]string, error) {

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("no JWT found in the Authorization header")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("invalid Authorization header format")
	}

	cookie := strings.TrimPrefix(authHeader, "Bearer ")

	fmt.Println("cookie is: ", cookie)
	claims, err := GetClaimsFromCookie(cookie, SecretKey)
	if err != nil {
		return nil, err
	}

	fmt.Println("claims are: ", claims)

	session, err := database.Redis.GetHMap(claims.Id)
	if err != nil {
		return nil, err
	}

	fmt.Println("session is: ", session)

	return session, nil
}
