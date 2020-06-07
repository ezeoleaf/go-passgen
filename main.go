package main

import (
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	letters     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers     = "0123456789"
	letterType  = 1
	numberType  = 2
	upperLetter = 1
	lowerLetter = 2
)

type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func initRoutes(e *echo.Echo) {
	e.GET("/", home())
	e.GET("/:length", getPass())
}

func startServer() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e.Renderer = renderer

	initRoutes(e)

	e.Logger.Fatal(e.Start(os.Getenv("APP_PORT")))
}

func home() echo.HandlerFunc {
	return func(c echo.Context) error {
		length, _ := strconv.Atoi(os.Getenv("PASS_LENGTH"))
		password := generatePassword(length)

		return c.Render(http.StatusOK, "template.html", map[string]interface{}{
			"password": password,
		})
	}
}

func getPass() echo.HandlerFunc {
	return func(c echo.Context) error {
		length, _ := strconv.Atoi(c.Param("length"))
		password := generatePassword(length)

		return c.Render(http.StatusOK, "template.html", map[string]interface{}{
			"password": password,
		})
	}
}

func generatePassword(length int) string {
	rand.Seed(time.Now().UnixNano())

	p := make([]byte, length)
	prevLetter := 0

	for i := 0; i < length; i++ {
		randType := rand.Intn(2) + 1
		if randType == letterType {
			if prevLetter == 0 {
				prevLetter = rand.Intn(2) + 1
			}
			ls := letters
			if prevLetter == upperLetter {
				ls = strings.ToLower(ls)
				prevLetter = lowerLetter
			} else if prevLetter == lowerLetter {
				ls = strings.ToUpper(ls)
				prevLetter = upperLetter
			}

			loc := rand.Intn(len(ls)) + 0
			p[i] = ls[loc]
		} else if randType == numberType {
			loc := rand.Intn(len(numbers)) + 0
			p[i] = numbers[loc]
		}
	}

	return string(p)
}

func main() {
	startServer()
}
