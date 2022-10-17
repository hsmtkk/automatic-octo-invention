package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port, err := getPort()
	if err != nil {
		log.Fatal(err)
	}

	t := &Template{templates: template.Must(template.ParseGlob("templates/*.html"))}

	// Echo instance
	e := echo.New()
	e.Renderer = t

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", getRoot)
	e.POST("/", postRoot)
	e.GET("/api/examples", examples)
	e.GET("/api/greetings", greetings)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

type param struct {
	Message string
}

// Handler
func getRoot(c echo.Context) error {
	req := c.Request()
	reqBytes, err := httputil.DumpRequest(req, false)
	if err != nil {
		return fmt.Errorf("httputil.DumpRequest failed; %w", err)
	}
	fmt.Printf("request: %s\n", string(reqBytes))

	return c.Render(http.StatusOK, "index", param{Message: "App Engine 勉強会にようこそ"})
}

func postRoot(c echo.Context) error {
	req := c.Request()
	reqBytes, err := httputil.DumpRequest(req, false)
	if err != nil {
		return fmt.Errorf("httputil.DumpRequest failed; %w", err)
	}
	fmt.Printf("request: %s\n", string(reqBytes))

	return c.Render(http.StatusOK, "index", param{Message: req.FormValue("message")})
}

type Guest struct {
	Author string `json:"author"`
	ID     int    `json:"id"`
}

func examples(ectx echo.Context) error {
	guests := []Guest{
		{"Tsuyoshi Igarashi", 1},
		{"Ryutaro Miyayama", 2},
		{"Mai Shirakawa", 3},
	}
	return ectx.JSON(http.StatusOK, guests)
}

type Greeting struct {
	ID      int    `json:"id"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

func greetings(ectx echo.Context) error {
	greets := []Greeting{
		{1, "Tsuyoshi Igarashi", "Hello"},
		{2, "Ryutaro Miyayama", "Looks good to me"},
	}
	return ectx.JSON(http.StatusOK, greets)
}

func getPort() (int, error) {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		return 0, fmt.Errorf("you must define PORT env var")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s as int; %w", portStr, err)
	}
	return port, nil
}
