package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	//GET
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"bigoffice", "/big-office", "GET", []postData{}, http.StatusOK},
	{"mediumoffice", "/medium-office", "GET", []postData{}, http.StatusOK},
	{"sharedoffice", "/shared-office", "GET", []postData{}, http.StatusOK},
	{"searchAv", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"makeRes", "/make-reservation", "GET", []postData{}, http.StatusOK},
	//POST
	{"postSearchAv", "/search-availability", "POST", []postData{
		{key: "start", value: "16-09-2020"},
		{key: "end", value: "17-09-2020"},
	}, http.StatusOK},
	{"postSearchAvJSON", "/search-availability-json", "POST", []postData{
		{key: "start", value: "16-09-2020"},
		{key: "end", value: "17-09-2020"},
	}, http.StatusOK},
	{"makeResPost", "/make-reservation", "POSTs", []postData{
		{key: "first_name", value: "Pier"},
		{key: "last_name", value: "Paul"},
		{key: "email", value: "Pier@pier.com"},
		{key: "phone", value: "65656565"},
	}, http.StatusOK},
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
			values := url.Values{}
			for _, param := range test.params {
				values.Add(param.key, param.value)
			}
			res, err := testServer.Client().PostForm(testServer.URL+test.url, values)
			if err != nil {
				t.Errorf("Error: %s", err)
			}
			if res.StatusCode != test.expectedStatus {
				t.Errorf("In test %s expected status code %d, got %d", test.testName, test.expectedStatus, res.StatusCode)
			}
		}
	}
}
