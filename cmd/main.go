package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Link struct {
	Id  string
	Url string
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var linkMap = map[string]*Link{"example": {Id: "example", Url: "https://example.com"}}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.GET("/:id", RedirectHandler)
	e.GET("/", IndexHandler)
	e.POST("/submit", SubmitHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	var result []byte

	for i := 0; i < length; i++ {
		index := seededRand.Intn(len(charset))
		result = append(result, charset[index])
	}

	return string(result)
}

func RedirectHandler(c echo.Context) error {
	id := c.Param("id")
	link, found := linkMap[id]

	if !found {
		return c.String(http.StatusNotFound, "Link not found")
	}

	return c.Redirect(http.StatusMovedPermanently, link.Url)
}

func IndexHandler(c echo.Context) error {
	html := `
		<h1>Submit a new website</h1>
		<form action="/submit" method="POST">
		<label for="url">Website URL:</label>
		<input type="text" id="url" name="url">
		<input type="submit" value="Submit">
		</form>
		<h2>Existing Links </h2>
		<ul>`

	for _, link := range linkMap {
		html += `<li><a href="/` + link.Id + `">` + link.Id + `</a></li>`
	}
	html += `</ul>`

	return c.HTML(http.StatusOK, html)
}

func SubmitHandler(c echo.Context) error {
	url := c.FormValue("url")
	if url == "" {
		return c.String(http.StatusBadRequest, "URL is required")
	}

	if !(len(url) >= 4 && (url[:4] == "http" || url[:5] == "https")) {
		url = "https://" + url
	}

	id := generateRandomString(8)

	linkMap[id] = &Link{Id: id, Url: url}

	return c.Redirect(http.StatusSeeOther, "/")
}
