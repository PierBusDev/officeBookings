package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/pierbusdev/conferenceRoomBookings/internal/config"
	"github.com/pierbusdev/conferenceRoomBookings/internal/driver"
	"github.com/pierbusdev/conferenceRoomBookings/internal/forms"
	"github.com/pierbusdev/conferenceRoomBookings/internal/helpers"
	"github.com/pierbusdev/conferenceRoomBookings/internal/models"
	"github.com/pierbusdev/conferenceRoomBookings/internal/render"
	"github.com/pierbusdev/conferenceRoomBookings/internal/repository"
	"github.com/pierbusdev/conferenceRoomBookings/internal/repository/dbrepo"
)

type Repository struct {
	AppConfig *config.AppConfig
	DB        repository.DatabaseRepo
}

var Repo *Repository

// NewRepo creates a new Repository
func NewRepo(c *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		AppConfig: c,
		DB:        dbrepo.NewPostgresDBRepo(db.SQL, c),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (rep *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.html", &models.TemplateData{})
}

func (rep *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.html", &models.TemplateData{})
}

func (rep *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := rep.AppConfig.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("Could not get reservation from session"))
		return
	}

	//getting the office name to show it in the form in the page
	office, err := rep.DB.GetOfficeById(reservation.OfficeID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	reservation.Office.OfficeName = office.OfficeName
	//saving it in session too ( gonna need it for the summary)
	rep.AppConfig.Session.Put(r.Context(), "reservation", reservation)

	//converting dates to format used in the frontend
	startDate := reservation.StartDate.Format("02-01-2006")
	endDate := reservation.EndDate.Format("02-01-2006")
	stringMap := make(map[string]string)
	stringMap["start_date"] = startDate
	stringMap["end_date"] = endDate

	data := make(map[string]any)
	data["reservation"] = reservation

	render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

func (rep *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := rep.AppConfig.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("Could not get reservation from session"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//======= TODO ==========================================================
	//All the part below can be simplifed: now that I get the reservation data via the session I don't need to work with the form
	//casting dates from string to time.Time:

	startDateFormatted, err := convertStringDateIntoTime(r.Form.Get("start_date"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDateFormatted, err := convertStringDateIntoTime(r.Form.Get("end_date"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// ============= OLD to remove after checking convertStringDateIntoTime works
	// startDateString := r.Form.Get("start_date")
	// endDateString := r.Form.Get("end_date")
	// layoutDateInputFormat := "02-01-2006"
	// startDate, err := time.Parse(layoutDateInputFormat, startDateString)
	// if err != nil {
	// 	helpers.ServerError(w, err)
	// 	return
	// }
	// endDate, err := time.Parse(layoutDateInputFormat, endDateString)
	// if err != nil {
	// 	helpers.ServerError(w, err)
	// 	return
	// }
	// //formatting for db
	// layoutDateOutputFormat := "2006-01-02"
	// startDateFormatted, _ := time.Parse(layoutDateOutputFormat, startDate.Format("2006-01-02"))
	// endDateFormatted, _ := time.Parse(layoutDateOutputFormat, endDate.Format("2006-01-02"))

	//converting  office_id from string to int:
	officeId, err := strconv.Atoi(r.Form.Get("office_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	//========================================================================

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]any)
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	//writing to db the reservation
	newReservationId, err := rep.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	//putting back reservation in the session
	rep.AppConfig.Session.Put(r.Context(), "reservation", reservation)

	resctriction := models.OfficeRestriction{
		StartDate:     startDateFormatted,
		EndDate:       endDateFormatted,
		OfficeID:      officeId,
		ReservationID: newReservationId,
		RestrictionID: 1,
	}

	err = rep.DB.InsertOfficeRestriction(resctriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rep.AppConfig.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (rep *Repository) BigOffice(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "big-office.page.html", &models.TemplateData{})
}

func (rep *Repository) MediumOffice(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "medium-office.page.html", &models.TemplateData{})
}

func (rep *Repository) SharedOffice(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "shared-office.page.html", &models.TemplateData{})
}

func (rep *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.html", &models.TemplateData{})
}

func (rep *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	//casting dates from string to time.Time:
	startDateString := r.Form.Get("start")
	endDateString := r.Form.Get("end")
	layoutDateInputFormat := "02-01-2006"
	startDate, err := time.Parse(layoutDateInputFormat, startDateString)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layoutDateInputFormat, endDateString)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	//formatting for db
	layoutDateOutputFormat := "2006-01-02"
	startDateFormatted, _ := time.Parse(layoutDateOutputFormat, startDate.Format("2006-01-02"))
	endDateFormatted, _ := time.Parse(layoutDateOutputFormat, endDate.Format("2006-01-02"))

	offices, err := rep.DB.SearchAvailabilityForAllOffices(startDateFormatted, endDateFormatted)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(offices) == 0 { //no available offices
		rep.AppConfig.Session.Put(r.Context(), "error", "There is no availability")
		rep.AppConfig.InfoLog.Println("no available offices")

		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}
	//else I have offices
	rep.AppConfig.InfoLog.Println("available offices:")
	for _, i := range offices {
		rep.AppConfig.InfoLog.Println(i.OfficeName)
	}

	//pass data to a page
	data := make(map[string]any)
	data["offices"] = offices
	//store some of the data in the session object to allow a better page rerendering
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	rep.AppConfig.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-office.page.html", &models.TemplateData{
		Data: data,
	})

}

type jsonResponse struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	OfficeID  string `json:"office_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (rep *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	startDate := r.Form.Get("start")
	endDate := r.Form.Get("end")
	officeID, _ := strconv.Atoi(r.Form.Get("office_id"))

	startDateFormatted, err := convertStringDateIntoTime(startDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDateFormatted, err := convertStringDateIntoTime(endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	available, err := rep.DB.SearchAvailabilityByDatesByOfficeId(startDateFormatted, endDateFormatted, officeID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	resp := jsonResponse{
		Ok:        available,
		Message:   "",
		StartDate: startDate,
		EndDate:   endDate,
		OfficeID:  strconv.Itoa(officeID),
	}

	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (rep *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

func (rep *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := rep.AppConfig.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		rep.AppConfig.ErrorLog.Println("Cannot get item from session")
		rep.AppConfig.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//cleaning up the session, we already have the reservation data in the reservation variable (see above)
	rep.AppConfig.Session.Remove(r.Context(), "reservation")

	data := make(map[string]any)
	data["reservation"] = reservation

	startDate := reservation.StartDate.Format("02-01-2006")
	endDate := reservation.EndDate.Format("02-01-2006")
	stringMap := make(map[string]string)
	stringMap["start_date"] = startDate
	stringMap["end_date"] = endDate

	render.Template(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

//ChooseOffice is an handler called from the choose-office.page.html template and has to get the passed id in the url parameter
//and update the session reservation object so that it can be used in the make-reservation.page.html template
func (rep *Repository) ChooseOffice(w http.ResponseWriter, r *http.Request) {
	officeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation, ok := rep.AppConfig.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	reservation.OfficeID = officeID

	// put all back in session
	rep.AppConfig.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

//BookOffice takes url parameters, builds a new reservation session object and passes it to make-reservation page
func (rep *Repository) BookOffice(w http.ResponseWriter, r *http.Request) {
	officeID, err := strconv.Atoi(r.URL.Query().Get("office_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var reservation models.Reservation
	reservation.OfficeID = officeID
	reservation.StartDate, err = convertStringDateIntoTime(startDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	reservation.EndDate, err = convertStringDateIntoTime(endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	office, err := rep.DB.GetOfficeById(officeID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	reservation.Office.OfficeName = office.OfficeName

	rep.AppConfig.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

//convertStringDateIntoTime is a helper method to convert a date in a string on the fronted to a time.Time object
func convertStringDateIntoTime(date string) (time.Time, error) {
	layoutDateInputFormat := "02-01-2006" //format of date in the frontend
	dateConverted, err := time.Parse(layoutDateInputFormat, date)
	if err != nil {
		return dateConverted, err
	}

	//formatting for db
	layoutDateOutputFormat := "2006-01-02" //format of date in db
	dateFormatted, err := time.Parse(layoutDateOutputFormat, dateConverted.Format("2006-01-02"))
	if err != nil {
		return dateFormatted, err
	}

	return dateFormatted, nil

}
