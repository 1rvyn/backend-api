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
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return nil, errors.New("no JWT cookie found")
	}

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
