package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"
	"time"
	"urlshort/store"

	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Store ...

var linkMap = map[string]*store.Link{}
var linkMapMutex sync.RWMutex

var baseurl *string
var db *sql.DB

// InitDB ...
func init() {

	baseurl = flag.String("url", "127.0.0.1:8080", "The URL (domain) that the server is running on")
	flag.Parse()

	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to DATABASE:%v", err))
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS links(
	id VARCHAR(255) PRIMARY KEY,
	url TEXT NOT NULL
	)
	`)

	if err := LoadLinks_in_Memory(); err != nil {
		panic(fmt.Sprintf("Failed to load links into memory: %v", err))
	}
	if err != nil {
		log.Fatal("Error creating table", err)
	}

}

// Save links into DB ...
func LoadLinks_in_Memory() error {
	rows, err := db.Query("SELECT id, url FROM links")
	if err != nil {
		return fmt.Errorf("query error: %v", err)
	}

	defer rows.Close()
	linkMapMutex.Lock()
	defer linkMapMutex.Unlock()

	for rows.Next() {
		var id, url string
		if err := rows.Scan(&id, &url); err != nil {
			return fmt.Errorf("scan error: %v", err)
		}
		linkMap[id] = &store.Link{Id: id, Url: url}
	}
	return nil
}

// Main ...
func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.GET("/", IndexHandler)
	e.GET("/:id", RedirectHandler)
	e.POST("/submit", SubmitHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

// Generator ...
func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	var result []byte

	for i := 0; i < length; i++ {
		index := seededRand.Intn(len(charset))
		result = append(result, charset[index])
	}

	return string(result)
}

// RedirectHandler ...
func RedirectHandler(c echo.Context) error {
	id := c.Param("id")

	linkMapMutex.RLock()
	link, found := linkMap[id]
	linkMapMutex.RUnlock()

	if !found {
		return c.String(http.StatusNotFound, "Link not found")
	}

	return c.Redirect(http.StatusMovedPermanently, link.Url)
}

// IndexHandler ...
func IndexHandler(c echo.Context) error {
	html := `
		<h1>URL Shorter</h1>
		<form action="/submit" method="POST">
		<label for="url">Enter link here:</label>
		<input type="text" id="url" name="url">
		<input type="submit" value="Shorten URL">
		</form>
		`

	// for _, link := range linkMap {
	// 	html += `<li><a href="/` + link.Id + `">` + link.Id + `</a></li>`
	// }

	return c.HTML(http.StatusOK, html)

}

// SubmitHandler ...
func SubmitHandler(c echo.Context) error {
	url := c.FormValue("url")
	if url == "" {
		back := `<h1> Error field is empty </h1>
		<a href="/">Back and try again</a>
		`
		return c.HTML(http.StatusBadRequest, back)
	}

	validUrl, err := ValidURL(url)
	if err != nil {
		back := `<h1> Invalid URL </h1>
		<a href="/">Back and try again</a>
		`
		return c.HTML(http.StatusBadRequest, back)
	}

	var id string
	for {
		id = generateRandomString(6)
		if _, exists := linkMap[id]; !exists {
			break
		}

	}

	linkMapMutex.Lock()
	linkMap[id] = &store.Link{Id: id, Url: validUrl}
	linkMapMutex.Unlock()

	if _, err := db.Exec("INSERT INTO links (id, url) VALUES ($1, $2)", id, validUrl); err != nil {
		return fmt.Errorf("failed to save link to database: %v", err)
	}

	html := fmt.Sprintf(`
	<h1>Your shortened URL</h1>
	<input type="text" style="width:200px;height:25px;" value="http://localhost:8080/%s" readonly>
	<button onclick="copyToClipboard()">Copy URL</button>
	<script>
	function copyToClipboard(){
		var input = document.querySelector("input");
		input.select();
		document.execCommand("copy");
		alert("URL is coppied " + input.value);
}
		</script>
		<a href="/">Back</a>
	`, id)

	return c.HTML(http.StatusOK, html)
}

func ValidURL(s string) (string, error) {

	parsedURL, err := url.ParseRequestURI(s)
	if err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https") && parsedURL.Host != "" {
		match, _ := regexp.MatchString(`^[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`, parsedURL.Host)
		if match {
			return parsedURL.String(), nil
		}
	}

	s = "https://" + s
	parsedURL, err = url.ParseRequestURI(s)
	if err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https") && parsedURL.Host != "" {
		match, _ := regexp.MatchString(`^[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`, parsedURL.Host)
		if match {
			return parsedURL.String(), nil
		}
	}

	return "", fmt.Errorf("invalid URL: %s", s)
}
