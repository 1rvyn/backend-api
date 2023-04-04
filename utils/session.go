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

	cookieHeader := c.Cookies("Cookie")
	if cookieHeader == "" {
		return nil, errors.New("no JWT found in the Cookie header")
	}

	// Extract the jwt cookie from the cookie header
	jwtCookie := ""
	cookies := strings.Split(cookieHeader, ";")
	for _, cookie := range cookies {
		if strings.HasPrefix(strings.TrimSpace(cookie), "jwt=") {
			jwtCookie = strings.TrimSpace(strings.TrimPrefix(cookie, "jwt="))
			break
		}
	}

	fmt.Println("new cookie format is: ", jwtCookie)
	if jwtCookie == "" {
		return nil, errors.New("no JWT cookie found")
	}

	fmt.Println("cookie is: ", jwtCookie)
	claims, err := GetClaimsFromCookie(jwtCookie, SecretKey)
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
