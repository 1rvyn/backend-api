package main

import (
	"authserver/database"
	"authserver/models"
	"authserver/utils"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

var SALT = os.Getenv("SALT")
var SecretKey = os.Getenv("JWT_SECRET")
var RedisSessionKey = os.Getenv("REDIS_SECRET")

func main() {
	database.ConnectDb()
	database.ConnectRedis()

	app := fiber.New()

	// simplified CORS - it is a subdomain, so it should be fine
	app.Use(func(c *fiber.Ctx) error {
		c.Response().Header.Set("Access-Control-Allow-Origin", "https://irvyn.xyz")
		c.Response().Header.Set("Access-Control-Allow-Credentials", "true")
		c.Response().Header.Set("Access-Control-Allow-Headers", "Set-Cookie, Cookie , Content-Type")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(200)
		}
		return c.Next()
	})

	setupRoutes(app)

	//err := app.Listen(":8080")
	//if err != nil {
	//	panic(err)
	//}

	//go func() { // for testing using pprof - import _ "net/http/pprof"
	//	http.ListenAndServe("localhost:6060", nil)
	//}()

	// Your existing fiber web app code
	err := app.Listen(":8080")
	if err != nil {
		return
	}

}

func setupRoutes(app *fiber.App) {
	app.Post("/register", Register)
	app.Post("/login", Login)
	app.Post("/user", getUser)
	app.Post("/logout", Logout)
	app.Get("/", Status)

	app.Post("/session", getUserfromSession)
	app.Post("/code", Code)
	app.Post("/account", Account) // return users account from their cookie
	// app.Post("/api/test1", test1)
}

func Account(c *fiber.Ctx) error {
	// get the user from the cookie
	// return their account
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return c.SendStatus(401)
	}
	// get the user from the cookie

	// look for the user in redis
	// if it is there, return the user

	session, err := database.Redis.GetHMap(cookie)
	if err != nil {
		return err
	}
	// they must have a redis session to get a result

	//TODO: re-do this in future

	fmt.Println("session we got from redis: ", session)

	// get the user from the database
	var user models.Account
	if err := database.Database.Db.Where("email = ?", session["email"]).First(&user).Error; err != nil {
		return c.SendStatus(401)
	} else {
		return c.JSON(user)
	}
}

func Code(c *fiber.Ctx) error {

	var data map[string][]byte

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	fmt.Println("\nthe data is : ", data)

	cookie := c.Cookies("jwt")

	fmt.Println("\nthe cookie is :", cookie)

	// validate the cookie
	claims, err := utils.GetClaimsFromCookie(cookie, SecretKey)
	if err != nil {
		// handle error
		return err
	} // validate the cookie

	fmt.Println("\n the claims at submission are : ", claims)

	issuer := claims.Issuer
	fmt.Println("\nthe issuer (userID) is : ", issuer)

	// compare the issuer from the claims to the user from the stored session in redis
	session, err := database.Redis.GetHMap(cookie)
	if err != nil {
		return err
	}

	if issuer != session["userID"] {
		// we should log this as a potential security breach
		// - someone is trying to submit code with a cookie that is not theirs

		// save info on what occured to the errors table

		database.Database.Db.Create(&models.Error{
			Message:    "Cookies didnt match",
			CreatedAt:  time.Now(),
			User:       session["userID"],
			Submission: data["codeitem"],
			IP:         c.Get("X-Forwarded-For"),
			Session:    session,
		})
		fmt.Println("\n Potential Malicious Activity Detected (Saved to errors table) - Cookies didnt match")
		return c.SendStatus(401)
	}

	// else the cookie is valid, save the submission to the database & mark it as pending

	//TODO: Set up getting the language from the frontend
	submission := models.Submission{
		Code:     data["codeitem"],
		UserID:   session["userID"],
		Language: "Python3",
		IP:       c.Get("X-Forwarded-For"),
	}

	database.Database.Db.Create(&submission)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "code was submitted",
	})
}

func Status(c *fiber.Ctx) error {
	// create a bunch of accounts

	return c.SendString("Hello world 👋!")
}

func Login(c *fiber.Ctx) error {
	var loginData map[string]string

	if err := c.BodyParser(&loginData); err != nil {
		return err
	}

	// check the email exists
	var user models.Account

	if err := database.Database.Db.Where("email = ?", loginData["email"]).First(&user).Error; err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	} else {
		// check the password
		hashedPw := make(chan []byte)

		go func() {
			hashedPw <- utils.HashPassword(loginData["password"], []byte(SALT))
			close(hashedPw)
		}()

		// Wait for the hashed password to be returned
		encpw := <-hashedPw

		// compare the passwords using bytes (faster than casting to string)
		if !bytes.Equal(user.EncryptedPassword, encpw) {
			fmt.Println("the passwords do not match")
			return c.Status(401).JSON(fiber.Map{
				"message": "incorrect password",
			})
		} else {
			// create a jwt token and create a session in redis
			claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
				Issuer:    strconv.Itoa(int(user.ID)),
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
			})

			token, err := claims.SignedString([]byte(SecretKey))

			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{
					"message": "could not create cookie",
				})
			}

			cookie := fiber.Cookie{
				Name:     "jwt",
				Value:    token,
				Expires:  time.Now().Add(time.Hour * 24),
				HTTPOnly: true,
				SameSite: "None",
				Secure:   true,
				Path:     "/",
				Domain:   ".irvyn.xyz",
			}
			c.Cookie(&cookie)

			// make new token to represent the session
			sessionToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
				Issuer:    strconv.Itoa(int(user.ID)),
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
			}).SignedString([]byte(RedisSessionKey))

			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{
					"message": "could not create session",
				})
			}

			//fmt.Println("\n The x-forwarded-for header is: ", c.Get("X-Forwarded-For"))

			session := make(map[string]interface{})
			session["sessionID"] = sessionToken
			session["userID"] = user.ID
			session["token"] = token
			session["email"] = user.Email
			session["username"] = user.Name
			session["ip"] = c.Get("X-Forwarded-For")
			session["agent"] = c.Get("User-Agent")
			session["role"] = user.UserRole
			session["created_at"] = time.Now().Format("2006-01-02 15:04:05")

			// save the token in redis with their userID as their key
			err = database.Redis.PutHMap(token, session)
			if err != nil {
				return err
			} else {
				fmt.Println("\nsuccessfully saved session to redis")
			}

			return c.JSON(fiber.Map{
				"message": "successfully logged in",
			})
		}
	}
}

func getUserfromSession(c *fiber.Ctx) error {
	// get the users cookie
	cookie := c.Cookies("jwt")

	fmt.Println("the cookie is :", cookie)

	// search the cookie value in redis to get the session
	session, err := database.Redis.GetHMap(cookie)
	if err != nil {
		return err
	}

	// TODO: we need to validate the token (not a huge security risk but still)

	fmt.Println("\nthe session we got from redis is :", session)

	// get the user's ID from the session
	userID := session["userID"]
	sessionToken := session["sessionID"]

	fmt.Println("\nthe userID is :", userID)
	fmt.Println("\nthe IP is :", session["ip"])

	fmt.Println("\nthe sessionToken is :", sessionToken)

	return c.JSON(fiber.Map{
		"message": "successfully got user from session",
		"userID":  userID,
		"session": session,
	})
}

func getUser(c *fiber.Ctx) error {
	// get the user's ID from the cookie

	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	fmt.Println("user id is :" + c.Get("userID"))

	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	fmt.Println("the token.Raw is :", token.Raw)

	id, err := database.Redis.Get(token.Raw)
	if err != nil {
		fmt.Println("the error is :", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "user not found in Redis",
		})
	}

	// Return the user's ID
	return c.JSON(fiber.Map{
		"userID": id,
	})

}

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	existingUser := &models.Account{}

	if err := database.Database.Db.Where("email = ?", data["email"]).First(existingUser).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	} else {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "email already in use",
		})
	}
	// Create a channel to receive the hashed password
	hashedPw := make(chan []byte)

	go func() {
		hashedPw <- utils.HashPassword(data["password"], []byte(SALT))
		close(hashedPw)
	}()

	// Wait for the hashed password to be returned
	encpw := <-hashedPw
	//fmt.Printf("encpw is : %v \n", encpw)

	if encpw == nil {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Error hashing password",
		})
	}

	user := models.Account{
		Name:              data["name"],
		Email:             data["email"],
		EncryptedPassword: encpw,
	}

	if err := database.Database.Db.Create(&user).Error; err != nil {
		return err
	}
	// return success
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Successfully registered user",
	})
}

func Logout(c *fiber.Ctx) error {
	// delete the cookie
	fmt.Println("\nlogout was called")

	rcookie := c.Cookies("jwt")

	fmt.Println("\nthe 'r' cookie is :", c.Cookies("jwt"))
	// remove the session from redis

	err := database.Redis.DeleteHMap(rcookie)
	if err != nil {
		return err
	} else {
		fmt.Println("\nSuccessfully deleted the session from redis")
	}

	// set the cookie to expire
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour * 24),
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true,
		Path:     "/",
		Domain:   ".irvyn.xyz",
	}
	c.Cookie(&cookie) // this returns a cookie with the date that is expired

	return c.JSON(fiber.Map{
		"message": "successfully logged out",
	})
}
