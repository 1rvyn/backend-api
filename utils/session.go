package utils

import (
	"authserver/database"
	"errors"
	"github.com/gofiber/fiber/v2"
	"os"
)

var SecretKey = os.Getenv("JWT_SECRET")

func GetSession(c *fiber.Ctx) (map[string]string, error) {
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return nil, errors.New("no JWT cookie found")
	}

	claims, err := GetClaimsFromCookie(cookie, SecretKey)
	if err != nil {
		return nil, err
	}

	session, err := database.Redis.GetHMap(claims.Id)
	if err != nil {
		return nil, err
	}

	return session, nil
}
