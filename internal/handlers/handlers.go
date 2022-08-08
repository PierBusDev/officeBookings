package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pierbusdev/conferenceRoomBookings/internal/config"
	"github.com/pierbusdev/conferenceRoomBookings/internal/forms"
	"github.com/pierbusdev/conferenceRoomBookings/internal/models"
	"github.com/pierbusdev/conferenceRoomBookings/internal/render"
)

type Repository struct {
	AppConfig *config.AppConfig
}

var Repo *Repository

// NewRepo creates a new Repository
func NewRepo(c *config.AppConfig) *Repository {
	return &Repository{
		AppConfig: c,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (rep *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	rep.AppConfig.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
}

func (rep *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIp := rep.AppConfig.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	render.RenderTemplate(w, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (rep *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]any)
	data["reservation"] = emptyReservation
	render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (rep *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("error in parsing data:", err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]any)
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	rep.AppConfig.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (rep *Repository) BigOffice(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "big-office.page.html", &models.TemplateData{})
}

func (rep *Repository) MediumOffice(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "medium-office.page.html", &models.TemplateData{})
}

func (rep *Repository) SharedOffice(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "shared-office.page.html", &models.TemplateData{})
}

func (rep *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.html", &models.TemplateData{})
}

func (rep *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start, end := r.Form.Get("start"), r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (rep *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		Ok:      true,
		Message: "available!",
	}

	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (rep *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.html", &models.TemplateData{})
}

func (rep *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := rep.AppConfig.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("Cannot get item from session")
		rep.AppConfig.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//cleaning up the session, we already have the reservation data in the reservation variable (see above)
	rep.AppConfig.Session.Remove(r.Context(), "reservation")

	data := make(map[string]any)
	data["reservation"] = reservation
	render.RenderTemplate(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data,
	})
}
