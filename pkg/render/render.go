package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/pierbusdev/conferenceRoomBookings/pkg/config"
	"github.com/pierbusdev/conferenceRoomBookings/pkg/models"
)

var appConfig *config.AppConfig

// NewTemplate  set the app configuration
func NewTemplate(c *config.AppConfig) {
	appConfig = c
}

func AddDefaultData(templData *models.TemplateData) *models.TemplateData {
	//TODO add default data useful for all the pages in the future
	return templData
}

// RenderTemplate renders given tmpl template using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, templData *models.TemplateData) {
	var templateCache map[string]*template.Template
	var err error
	if appConfig.UseCache {
		//obtaining template cache from the appConfig
		templateCache = appConfig.TemplateCache
	} else {
		//just rebuild a new instance right now and use that
		templateCache, err = CreateTemplateCache()
		if err != nil {
			log.Fatal("Can't create a template cache: ", err)
		}
	}

	templ, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("Can't find template in the cache")
	}

	buf := new(bytes.Buffer)
	err = templ.Execute(buf, AddDefaultData(templData))
	if err != nil {
		log.Println(err)
	}

	//render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("Some error occurred while writing the template to the browser:", err)
	}
}

//CreateTemplateCache will return a map containing all the templates which are present inside the `templates` folder
// indexed by the filename, or an error if it occurs.
func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := make(map[string]*template.Template)

	//find all the files with .page pseudoextension
	pages, err := filepath.Glob("./templates/*page.html")
	if err != nil {
		return templateCache, err
	}

	//we need to check if there are layouts to add to the various pages
	for _, page := range pages {
		fileName := filepath.Base(page) //basically the filename stripped of the path
		templ, err := template.New(fileName).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		//checking for layouts presence
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return templateCache, err
		}

		if len(matches) > 0 { //if we have layouts we have to incorporate them in the template
			templ, err = templ.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[fileName] = templ
	}

	return templateCache, nil
}
