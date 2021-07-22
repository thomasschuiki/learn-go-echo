package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/thomasschuiki/learn-go-echo/models"
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: "method=${method}, uri=${uri}, status=${status}\n"}))
	e.Use(serverHeader)
	e.GET("/health-check", HealthCheck)
	e.GET("/cats/:data", GetCats)
	e.POST("/cats", AddCat)
	e.Logger.Fatal(e.Start(":8000"))
}

func GetCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("your cat name is: %s\nand cat type is: %s\n", catName, catType))
	} else if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType})
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Please specify the data type as string or json"})
	}
}

func AddCat(c echo.Context) error {
	type Cat struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	cat := Cat{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&cat)
	if err != nil {
		log.Fatalf("Failed reading the request body %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	log.Printf("this is your cat %#v", cat)
	return c.String(http.StatusOK, "we got your cat!")
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
