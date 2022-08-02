package handlers

import (
	"net/http"

	"github.com/pierbusdev/basicWeb/pkg/config"
	"github.com/pierbusdev/basicWeb/pkg/render"
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
	render.RenderTemplate(w, "home.page.html")
}

func (rep *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.html")
}
