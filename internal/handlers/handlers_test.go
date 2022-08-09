package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

var dataForTests = []struct {
	testName       string
	url            string
	method         string
	params         []postData
	expectedStatus int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"bigoffice", "/big-office", "GET", []postData{}, http.StatusOK},
	{"mediumoffice", "/medium-office", "GET", []postData{}, http.StatusOK},
	{"sharedoffice", "/shared-office", "GET", []postData{}, http.StatusOK},
	{"searchAv", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"makeRes", "/make-reservation", "GET", []postData{}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range dataForTests {
		if test.method == "GET" {
			res, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Errorf("Error: %s", err)
			}
			if res.StatusCode != test.expectedStatus {
				t.Errorf("In test %s expected status code %d, got %d", test.testName, test.expectedStatus, res.StatusCode)
			}
		} else { //POST

		}
	}
}
