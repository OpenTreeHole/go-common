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
}
