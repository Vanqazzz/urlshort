package utils

import (
	"net/url"
	"regexp"
)

func IsURL(s string) bool {

	parsedURL, err := url.ParseRequestURI(s)
	if err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https") && parsedURL.Host != "" {
		match, _ := regexp.MatchString(`^[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`, parsedURL.Host)
		if match {
			return true
		}
	}

	s = "https://" + s
	parsedURL, err = url.ParseRequestURI(s)
	if err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https") && parsedURL.Host != "" {
		match, _ := regexp.MatchString(`^[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`, parsedURL.Host)
		if match {
			return true
		}
	}

	return false
}
