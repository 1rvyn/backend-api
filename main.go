package main

import (
	"authserver/database"
	"authserver/models"
	"authserver/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

var SALT = os.Getenv("SALT")
var SecretKey = os.Getenv("JWT_SECRET")
var RedisSessionKey = os.Getenv("REDIS_SECRET")

//var Flask = os.Getenv("FLASK_API_ENDPOINT")

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
	app.Get("/logout", Logout)
	app.Get("/", Status)

	app.Get("/submissions", getSubmissions)
	app.Post("/session", getUserfromSession)
	app.Post("/getsession", getSession)
	app.Post("/code", Code)
	app.Post("/account", Account) // return users account from their cookie
	app.Post("/bugreport", BugReport)
	app.Get("/question/:id/:language", Question)
	app.Get("/questionall/:id", Questionall)
	app.Get("/questions", Questions)
	app.Post("/new_question", CreateQuestion)
	app.Post("/results-endpoint", ResultsEndpoint)

	app.Post("/tested", Tested) // endpoint to return the GKE job status
	//app.Get("/mailgun", Mailgun)

	app.Get("/verify", VerifyAccount)

	app.Get("/admin", Admin)
	//app.Post("/vemail", VerifyEmail)
	// app.Post("/api/test1", test1)

	app.Use(pprof.New())

}

func Admin(c *fiber.Ctx) error {
	// Get the cookie from the request
	fmt.Println("Admin handler HIT")
	fmt.Println(c.GetReqHeaders())
	session, err := utils.GetSession(c)

	if err != nil {
		return err
	}
	fmt.Println(session)

	// check if the user is an admin
	if session["role"] == "2" {
		return c.SendStatus(200)
	} else {
		return c.SendStatus(403)
	}

}

func Tested(c *fiber.Ctx) error {
	// Parse the JSON payload
	var payload map[string]interface{}
	if err := json.Unmarshal(c.Body(), &payload); err != nil {
		return err
	}

	fmt.Println("Tested endpoint hit \n", payload)
	// Extract submission_id and results
	//fmt.Println(payload["submission_id"], payload["results"])

	// TODO: Update the submission record in your database with the results

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "results received",
	})
}

func ResultsEndpoint(c *fiber.Ctx) error {
	// Parse the JSON payload
	var payload map[string]interface{}
	if err := json.Unmarshal(c.Body(), &payload); err != nil {
		return err
	}
	// Extract submission_id and results
	fmt.Println(payload["submission_id"], payload["results"])

	// TODO: Update the submission record in your database with the results

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "results received",
	})
}

func getSession(c *fiber.Ctx) error {
	session, err := utils.GetSession(c)
	if err != nil {
		return err
	}
	fmt.Println(session)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

// TODO: Protect this route behind admin perms only

func CreateQuestion(c *fiber.Ctx) error {
	fmt.Println("CreateQuestion handler HIT")

	// check if user has a cookie and is logged in

	cookie := c.Cookies("jwt")
	if cookie == "" {
		return c.SendStatus(401)
	} else {
		// check if the user has a session in redis

		session, err := database.Redis.GetHMap(cookie)
		if err != nil {
			return err
		}

		if session == nil {
			return c.SendStatus(401)
		}
	}

	// get the question and its data
	var questionData models.Question

	if err := c.BodyParser(&questionData); err != nil {
		return err
	}

	fmt.Println(questionData)

	//TODO: Speed up this - marshal and unmarshal is slow

	// Unmarshal the TemplateCode JSON data into a map[string]string
	var templateCodeMap map[string]string
	if err := json.Unmarshal(questionData.TemplateCode, &templateCodeMap); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid template code JSON format.",
		})
	}

	// Validate the presence of the required languages in the TemplateCode field
	if templateCodeMap["python"] == "" || templateCodeMap["javascript"] == "" || templateCodeMap["go"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Template code for Python, JavaScript, and Go are required.",
		})
	}

	database.Database.Db.Create(&questionData)

	// return the question to the user

	return c.JSON(fiber.Map{
		"message": "Question created successfully.",
	})
}

func VerifyAccount(c *fiber.Ctx) error {
	// get the code from the url
	code := c.Query("code")
	email := c.Query("email")

	// get the user from the database
	var user models.Account
	if err := database.Database.Db.Where("email = ?", email).First(&user).Error; err != nil {
		return c.SendStatus(401)
	}

	// check to see if the user is already verified
	if user.Verified {
		return c.Redirect("https://irvyn.xyz/login?message=account+already+verified")
	}

	// check if the code is correct
	if strconv.Itoa(user.EmailCode) == code {
		// update the user
		database.Database.Db.Model(&user).Update("verified", true)
		return c.JSON("verified")
	}

	// update the user
	database.Database.Db.Model(&user).Update("verified", true)

	return c.Redirect("https://irvyn.xyz/login?message=successfully+verified+email")

}

// TODO: Only allow the frontend fiber app to access this route

func Questions(c *fiber.Ctx) error {
	// get all the questions from the database
	var questions []models.Question
	database.Database.Db.Find(&questions)

	// return the questions to the user
	return c.JSON(questions)
}

// TODO: Cache this / make it event driven

func Question(c *fiber.Ctx) error {
	// get the question with the ID from the URL
	id := c.Params("id")
	language := c.Params("language")
	fmt.Println("the language is: ", language)
	fmt.Println("the id is: ", id)

	var question models.Question
	database.Database.Db.Where("id = ?", c.Params("id")).First(&question)

	var templateCodeMap map[string]string
	// parse the individual language from the 'TemplateCode' item
	err := json.Unmarshal(question.TemplateCode, &templateCodeMap)
	if err != nil {
		fmt.Println("Error unmarshaling template code:", err)
		return c.SendStatus(500)
	}

	code, exists := templateCodeMap[language]
	if !exists {
		fmt.Println("Language not found in template code")
		return c.SendStatus(400)
	}

	response := map[string]string{
		language: code,
	}

	return c.JSON(response)
}

func Questionall(c *fiber.Ctx) error {
	// get the question with the ID from the URL
	id := c.Params("id")
	fmt.Println("the id is: ", id)

	var question models.Question
	database.Database.Db.Where("id = ?", c.Params("id")).First(&question)

	return c.JSON(question)
}

func BugReport(c *fiber.Ctx) error {

	// get the user from the cookie

	// if they arent logged in / dont have a valid cookie then we will return a 401
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return c.SendStatus(401)
	}

	// get their session from redis -
	// errors are handled in the function because my redis code is amazing
	session, err := database.Redis.GetHMap(cookie)
	if err != nil {
		return err
	}

	// get the body from the request
	var body map[string]string
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	fmt.Println("the body is: ", body)

	fmt.Println("the bug report is: ", body["bugReport"])

	// get the bug report
	bugReport := body["bugReport"]

	// create a new bug report
	bug := models.Bug{
		Email: session["email"],
		Title: body["title"],
		Body:  bugReport,
	}

	// save it to the database
	database.Database.Db.Create(&bug)

	return c.SendStatus(200)

}

func getSubmissions(c *fiber.Ctx) error {
	// get the user from the cookie
	cookie := c.Cookies("jwt")
	// get their session from redis
	if cookie == "" {
		return c.SendStatus(401)
	}

	session, err := database.Redis.GetHMap(cookie)
	if err != nil {
		return err
	}

	if session["email"] == "" {
		return c.SendStatus(401)
	}

	// get users submission from the submissions table
	var submissions []models.Submission
	if err := database.Database.Db.Where("user_id = ?", session["email"]).Find(&submissions).Error; err != nil {
		return c.SendStatus(401)
	} else {
		return c.JSON(submissions)
	}
}

func Account(c *fiber.Ctx) error {
	// get the user from the cookie
	// return their account
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return c.SendStatus(401)
	}

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
		userResponse := models.ResponseAccount{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UserRole:  user.UserRole,
		}
		return c.JSON(userResponse)
	}
}

func Code(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	cookie := c.Cookies("jwt")

	fmt.Println("data is: ", data)

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
			Submission: data["code"],
			IP:         c.Get("X-Forwarded-For"),
		})
		fmt.Println("\n Potential Malicious Activity Detected (Saved to errors table) - Cookies didnt match")
		return c.SendStatus(401)
	}

	// else the cookie is valid, save the submission to the database & mark it as pending

	submission := models.Submission{
		Code:     data["code"],
		UserID:   session["userID"],
		Language: data["language"],
		IP:       c.Get("X-Forwarded-For"),
	}

	database.Database.Db.Create(&submission)

	//TODO: mark the submission and return the output string

	if questionID, ok := data["QuestionID"]; ok && questionID == "1" && data["language"] == "python" {
		fmt.Println("QuestionID is 1 and language is python")
		// Send submitted code to the Flask API
		flaskAPIEndpoint := os.Getenv("FLASK_API_ENDPOINT")

		responseString, err := sendCodeToFlaskAPI(flaskAPIEndpoint, data["code"])
		if err != nil {
			fmt.Println("Error sending code to Flask API:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to send code to Flask API",
			})
		}

		// parse the response string into a JSON object
		fmt.Println("Response from Flask API:", responseString)
		var results []models.TestResult

		err = json.Unmarshal([]byte(responseString), &results)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to parse response from Flask API",
			})
		}

		fmt.Println("Results from Flask API after parsed into results struct:", results)

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "code was submitted",
			"result":  results,
		})

	} else {
		fmt.Println("we got a submission that isn't meant for GKE :) ")
		// use the local marking system
		output := utils.Marking(data["code"], data["QuestionID"], session["userID"], data["language"])
		fmt.Println("the output from marking is: ", output)
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "code was submitted",
			//"result":  results,
			"result": output,
		})
	}

}

// send as bytes which is more efficient
func sendCodeToFlaskAPI(url, code string) (string, error) {
	fmt.Println("Sending code to Flask API:", code)
	// Prepare a buffer to hold the form data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add the submitted code as a file field to the form
	fw, err := w.CreateFormFile("code", "two_sum.py")
	if err != nil {
		return "", err
	}
	_, err = io.Copy(fw, strings.NewReader(code))
	if err != nil {
		return "", err
	}
	w.Close()

	// Send a POST request with the form data
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code from Flask API: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)

	fmt.Println("BODY STRING from Flask API:", bodyString)

	return bodyString, nil
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

			//fmt.Println("\n The x-forwarded-for header is: ", c.Get("X-Forwarded-For"))

			ip := c.Get("X-Forwarded-For")
			agent := c.Get("User-Agent")
			// we dont pass the C across goroutines because it is not thread safe
			go func(ip, agent string) {

				// make new token to represent the session
				sessionToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
					Issuer:    strconv.Itoa(int(user.ID)),
					ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
				}).SignedString([]byte(RedisSessionKey))

				if err != nil {
					log.Println("Error creating session token:", err) // Log the error instead of returning it
				}

				session := make(map[string]interface{})
				session["sessionID"] = sessionToken
				session["userID"] = user.ID
				session["token"] = token
				session["email"] = user.Email
				session["username"] = user.Name
				session["ip"] = ip
				session["agent"] = agent
				session["role"] = user.UserRole
				session["created_at"] = time.Now().Format("2006-01-02 15:04:05")

				// save the token in redis with their userID as their key
				err = database.Redis.PutHMap(token, session)
				if err != nil {
					log.Println("Error creating session:", err) // Log the error instead of returning it
				} else {
					fmt.Println("\nsuccessfully saved session to redis")
				}
			}(ip, agent)

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
		EmailCode:         utils.GenerateVerficiationCode(),
		Verified:          false,
	}

	if err := database.Database.Db.Create(&user).Error; err != nil {
		return err
	}
	// since we have saved the user, we can now send them an email to verify their email address

	err := utils.SendMail(user.Email)
	if err != nil {
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

	// redirect to the "/" page
	return c.Redirect("https://irvyn.xyz", 200)
}
