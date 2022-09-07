package render

import (
	"net/http"
	"testing"

	"github.com/pierbusdev/conferenceRoomBookings/internal/models"
)

func TestNewTemplates(t *testing.T) {
	NewRenderer(&testAppConfig)
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates" //overwriting values in render.go
	templateCache, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	appConfig.TemplateCache = templateCache
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	var w myWriter
	err = Template(&w, r, "home.page.html", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}
	err = Template(&w, r, "not.existent.page.html", &models.TemplateData{})
	if err == nil {
		t.Error("rendering template that does not exist")
	}
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates" //overwriting values in render.go
	templateCache, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	if len(templateCache) == 0 {
		t.Error("template cache is empty")
	}
}

func TestAddDefaultData(t *testing.T) {
	var templateData models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	expectedFlashMessage := "flash message"

	session.Put(r.Context(), "flash", expectedFlashMessage)
	result := AddDefaultData(&templateData, r)
	if result.Flash != expectedFlashMessage {
		t.Errorf("Got flash message %s but expected %s", result.Flash, expectedFlashMessage)
	}

}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}
