package utils

import (
	"authserver/database"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
)

var SecretKey = os.Getenv("JWT_SECRET")

func GetSession(c *fiber.Ctx) (map[string]string, error) {

	jwt := c.Get("X-JWT")
	if jwt == "" {
		return nil, errors.New("no JWT found in the X-JWT header")
	}

	fmt.Println("X-JWT is: ", jwt)
	claims, err := GetClaimsFromCookie(jwt, SecretKey)
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
