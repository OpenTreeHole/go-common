package common

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestTesting(t *testing.T) {
	app := fiber.New(fiber.Config{
		AppName:               "test",
		ErrorHandler:          CommonErrorHandler,
		JSONDecoder:           json.Unmarshal,
		JSONEncoder:           json.Marshal,
		DisableStartupMessage: true,
	})

	type User struct {
		ID int `json:"id" query:"id"`
	}

	users := map[int]User{}
	for i := 1; i <= 10; i++ {
		users[i] = User{ID: i}
	}

	app.Use(MiddlewareGetUserID)
	app.Use(MiddlewareCustomLogger)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		userQuery, err := ValidateQuery[User](c)
		if err != nil {
			return err
		}
		user, ok := users[userQuery.ID]
		if !ok {
			return gorm.ErrRecordNotFound
		}
		return c.JSON(user)
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		user, err := ValidateBody[User](c)
		if err != nil {
			return err
		}
		users[user.ID] = *user
		return c.Status(fiber.StatusCreated).JSON(user)
	})

	app.Get("/jwt", func(c *fiber.Ctx) error {
		userID, err := GetUserID(c)
		if err != nil {
			return err
		}
		return c.JSON(Map{"user_id": userID})
	})

	app.Post("/form", func(c *fiber.Ctx) error {
		value := c.FormValue("data")
		return c.Status(201).SendString(value)
	})

	RegisterApp(app)

	// Test GET /
	DefaultTester.Get(t, RequestConfig{Route: "/", ExpectedBody: "Hello, World!"})

	// Test GET /users
	var user User
	DefaultTester.Get(t, RequestConfig{Route: "/users", RequestQuery: Map{"id": 1}, ResponseModel: &user})
	assert.EqualValues(t, 1, user.ID)

	// Test GET /users with a non-existing user
	DefaultTester.Get(t, RequestConfig{Route: "/users", RequestQuery: Map{"id": 11}, ExpectedStatus: fiber.StatusNotFound})

	// Test POST /users
	var newUser = User{ID: 11}
	DefaultTester.Post(t, RequestConfig{
		Route:         "/users",
		RequestBody:   newUser,
		ResponseModel: &newUser,
	})

	// Test POST /users with invalid body
	DefaultTester.Post(t, RequestConfig{
		Route:          "/users",
		RequestBody:    map[string]any{"id": "1"},
		ExpectedStatus: fiber.StatusBadRequest,
	})

	/* Test GET /jwt */
	// Test GET /jwt without token
	DefaultTester.Get(t, RequestConfig{Route: "/jwt", ExpectedStatus: fiber.StatusUnauthorized})

	// Test GET /jwt with invalid token
	UserTester.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2ODgzOTQwMDEsInVzZXJfaWQiOjEsInR5cGUiOiJhY2Nlc3MiLCJ1aWQiOjF9"
	UserTester.Get(t, RequestConfig{Route: "/jwt", ExpectedStatus: fiber.StatusUnauthorized})

	// Test GET /jwt with valid token and id field
	UserTester.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2ODgzOTQwMDEsImlkIjoxLCJ0eXBlIjoiYWNjZXNzIiwidWlkIjoxfQ.JQZdPizvZyI7-Fg8uHN45t4URShtVYtFvt9Mif7ArQk"
	UserTester.Get(t, RequestConfig{Route: "/jwt", ExpectedBody: `{"user_id":1}`})

	// Test GET /jwt with valid token and user_id field
	UserTester.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2ODgzOTQwMDEsInVzZXJfaWQiOjEsInR5cGUiOiJhY2Nlc3MiLCJ1aWQiOjF9.LBPhM9rK4zMctR1_-TTfOtrXmtXaAlAUzTwIGuJJhgI"
	UserTester.Get(t, RequestConfig{Route: "/jwt", ExpectedBody: `{"user_id":1}`})

	// Test GET /jwt with header X-CONSUMER-USERNAME
	UserTester.Token = ""
	UserTester.Get(t, RequestConfig{Route: "/jwt", ExpectedBody: `{"user_id":1}`, RequestHeaders: map[string]string{"X-CONSUMER-USERNAME": "1"}})

	// Test POST /form
	DefaultTester.Post(t, RequestConfig{
		Route:        "/form",
		RequestBody:  Map{"data": "test"},
		ExpectedBody: "test",
		ContentType:  fiber.MIMEApplicationForm,
	})
}
