package utils

import (
	"authserver/database"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func GetSession(c *fiber.Ctx) (map[string]string, error) {

	jwt := c.Get("X-JWT")
	if jwt == "" {
		return nil, errors.New("no JWT found in the X-JWT header")
	}

	fmt.Println("X-JWT is: ", jwt)

	session, err := database.Redis.GetHMap(jwt)
	if err != nil {
		return nil, err
	}

	fmt.Println("session is: ", session)

	return session, nil
}
