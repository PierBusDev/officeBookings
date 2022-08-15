package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/something", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("Form is not valid but should be valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/something", nil)
	form := New(r.PostForm)

	form.Required("name", "email")
	if form.Valid() {
		t.Error("Form is valid but should be invalid because required fields are missing")
	}

	postData := url.Values{}
	postData.Add("name", "John")
	postData.Add("email", "john@johnmail.com")

	r, _ = http.NewRequest("POST", "/something", nil)
	r.PostForm = postData
	form = New(r.PostForm)
	form.Required("name", "email")
	if !form.Valid() {
		t.Error("Form is not valid but should be valid because required data is present")
	}

}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/something", nil)
	form := New(r.PostForm)

	has := form.Has("name")
	if has {
		t.Error("Form has field but should not have it")
	}

	postData := url.Values{}
	postData.Add("name", "Pier")
	form = New(postData)
	has = form.Has("name")
	if !has {
		t.Error("Form does not have field but should have it")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/something", nil)
	form := New(r.PostForm)

	form.MinLength("name", 3)
	if form.Valid() {
		t.Error("Form is valid but should be invalid because field is not existent")
	}

	isError := form.Errors.Get("name")
	if isError == "" {
		t.Error("Should have an error, but did not get one")
	}

	postData := url.Values{}
	postData.Add("name", "Pier")
	form = New(postData)
	form.MinLength("name", 50)
	if form.Valid() {
		t.Error("Form is valid but should not be valid because field is NOT long enough")
	}

	postData = url.Values{}
	postData.Add("name", "Pier")
	form = New(postData)
	form.MinLength("name", 3)
	if !form.Valid() {
		t.Error("Form is not valid but should be valid because field is long enough")
	}
	isError = form.Errors.Get("name")
	if isError != "" {
		t.Error("Got an error, but should not have one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("Form is valid but should be invalid because field is not existent")
	}

	postData = url.Values{}
	postData.Add("email", "pierpier.pier")
	form = New(postData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("Form is valid but should not be valid because field is not in a correct email format")
	}

	postData = url.Values{}
	postData.Add("email", "pier@pier.com")
	form = New(postData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("Form is not valid but should be valid because field is in a correct email format")
	}
}
