package test

import (
	"testing"
	"urlshort/store"
	"urlshort/utils"
)

func TestEmptyField(t *testing.T) {

	link := store.Link{}

	if link.Url != "" {
		t.Error("Field empty", link.Url)
	}

}

func TestValidURL(t *testing.T) {
	validURL := "http://example.com"
	invalidURL := "ex-a-mple"

	if !utils.IsURL(validURL) {
		t.Errorf("Expected valid URL, but got false: %s", validURL)
	}
	if utils.IsURL(invalidURL) {
		t.Errorf("Expected invalid URL, but got true: %s", invalidURL)

	}
}
