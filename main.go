package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"

	"github.com/thomasschuiki/learn-go-echo/db"
	"github.com/thomasschuiki/learn-go-echo/models"
)

func main() {
	db.Init("./app.db")
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: "method=${method}, uri=${uri}, status=${status}\n"}))
	e.Use(serverHeader)
	e.GET("/health-check", HealthCheck)
	e.GET("/users", GetUsers)
	e.GET("/users/:name", GetUser)
	e.POST("/users", SignUp)
	e.POST("/signin", SignIn)
	e.Logger.Fatal(e.Start(":8000"))
}

func GetUsers(c echo.Context) error {
	users, err := db.AllUsers()
	if err != nil {
		log.Fatal("failed to get all users from db")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	return c.JSON(http.StatusOK, users)
}

func GetUser(c echo.Context) error {
	name := c.Param("name")
	user, err := db.GetUser(name)
	if err != nil {
		log.Fatal("failed to get user from db")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	return c.JSON(http.StatusOK, user)
}
func SignUp(c echo.Context) error {
	user := models.User{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		log.Fatalf("Failed reading the request body. Error: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}

	log.Printf("this is your user %#v", user)
	userid, err := db.CreateUser(user.Name, user.Password)
	log.Printf("created user with id %d", userid)
	return c.String(http.StatusOK, "user created successfully")
}

func SignIn(c echo.Context) error {

	user := models.User{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		log.Fatalf("Failed reading the request body. Error: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	storedUser, err := db.GetUser(user.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error)
	}

	return c.String(http.StatusOK, "user logged in successfully")
}

func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Custom-Header", "foobar")
		return next(c)
	}
}
func HealthCheck(c echo.Context) error {
	type HealthCheckResponse struct {
		Message string `json:"message"`
	}
	resp := HealthCheckResponse{Message: "Everything looking good!"}
	return c.JSON(http.StatusOK, resp)
}
