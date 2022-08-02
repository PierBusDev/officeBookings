package config

import "html/template"

//AppConfig represents the whole application configurations
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
}
