package config

import "html/template"

//AppConfig represents the whole application configurations
type AppConfig struct {
	TemplateCache map[string]*template.Template
}
