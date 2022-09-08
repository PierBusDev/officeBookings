package handlers

import (
	"context"
	"github.com/pierbusdev/conferenceRoomBookings/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
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
	expectedStatus int
}{
	//GET
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"bigoffice", "/big-office", "GET", http.StatusOK},
	{"mediumoffice", "/medium-office", "GET", http.StatusOK},
	{"sharedoffice", "/shared-office", "GET", http.StatusOK},
	{"searchAv", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"makeRes", "/make-reservation", "GET", http.StatusOK},

	//POST
	//{"postSearchAv", "/search-availability", "POST", []postData{
	//	{key: "start", value: "16-09-2020"},
	//	{key: "end", value: "17-09-2020"},
	//}, http.StatusOK},
	//{"postSearchAvJSON", "/search-availability-json", "POST", []postData{
	//	{key: "start", value: "16-09-2020"},
	//	{key: "end", value: "17-09-2020"},
	//}, http.StatusOK},
	//{"makeResPost", "/make-reservation", "POSTs", []postData{
	//	{key: "first_name", value: "Pier"},
	//	{key: "last_name", value: "Paul"},
	//	{key: "email", value: "Pier@pier.com"},
	//	{key: "phone", value: "65656565"},
	//}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range dataForTests {
		res, err := testServer.Client().Get(testServer.URL + test.url)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
		if res.StatusCode != test.expectedStatus {
			t.Errorf("In test %s expected status code %d, got %d", test.testName, test.expectedStatus, res.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		OfficeID: 1,
		Office: models.Office{
			ID:         1,
			OfficeName: "Big office",
		},
	}
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	responseRecorder := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Reservation handler got response code %d, want %d", responseRecorder.Code, http.StatusOK)
	}

	//test cases where the reservation is not in the session:
	//resetting req
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	responseRecorder = httptest.NewRecorder()
	reservation.OfficeID = 5 //see test-repo.go, id greater than 2 will cause errors
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler got response code %d, want %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	//test with not existing office
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler got response code %d, want %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	requestBody := "start_date=01-01-2050"
	requestBody = requestBody + "&end_date=02-01-2050"
	requestBody = requestBody + "&first_name=TestPier"
	requestBody = requestBody + "&last_name=TestPaul"
	requestBody = requestBody + "&email=test@pier.com"
	requestBody = requestBody + "&phone=65656565"
	requestBody = requestBody + "&office_id=1"

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler got response code %d, want %d", responseRecorder.Code, http.StatusSeeOther)
	}

	//test with missing body =============
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler got response code %d, want %d -> Missing response body", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	//test with invalid start date =============
	requestBody = "start_date=99-01-2050" //invalid format
	requestBody = requestBody + "&end_date=02-01-2050"
	requestBody = requestBody + "&first_name=TestPier"
	requestBody = requestBody + "&last_name=TestPaul"
	requestBody = requestBody + "&email=test@pier.com"
	requestBody = requestBody + "&phone=65656565"
	requestBody = requestBody + "&office_id=1"
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler got response code %d, want %d -> Invalid start date", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	//test with invalid start date =============
	requestBody = "start_date=01-01-2050"
	requestBody = requestBody + "&end_date=82-01-2050" //invalid format
	requestBody = requestBody + "&first_name=TestPier"
	requestBody = requestBody + "&last_name=TestPaul"
	requestBody = requestBody + "&email=test@pier.com"
	requestBody = requestBody + "&phone=65656565"
	requestBody = requestBody + "&office_id=1"
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler got response code %d, want %d -> Invalid end date", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	//test with invalid office id =============
	requestBody = "start_date=01-01-2050"
	requestBody = requestBody + "&end_date=02-01-2050"
	requestBody = requestBody + "&first_name=TestPier"
	requestBody = requestBody + "&last_name=TestPaul"
	requestBody = requestBody + "&email=test@pier.com"
	requestBody = requestBody + "&phone=65656565"
	requestBody = requestBody + "&office_id=somenonsense" //invalid format
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler got response code %d, want %d -> Invalid office id", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	//test with invalid data
	requestBody = "start_date=01-01-2050"
	requestBody = requestBody + "&end_date=02-01-2050"
	requestBody = requestBody + "&first_name=T" //first name too short
	requestBody = requestBody + "&last_name=TestPaul"
	requestBody = requestBody + "&email=test@pier.com"
	requestBody = requestBody + "&phone=65656565"
	requestBody = requestBody + "&office_id=1"
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler got response code %d, want %d -> Invalid data", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	//testing failing insert (see related function in test-repo.go)
	requestBody = "start_date=01-01-2050"
	requestBody = requestBody + "&end_date=02-01-2050"
	requestBody = requestBody + "&first_name=TestPier"
	requestBody = requestBody + "&last_name=TestPaul"
	requestBody = requestBody + "&email=test@pier.com"
	requestBody = requestBody + "&phone=65656565"
	requestBody = requestBody + "&office_id=2" //this will make the db insert fail
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler got response code %d, want %d -> failed insertion", responseRecorder.Code, http.StatusTemporaryRedirect)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
