package models

//TemplateData contains data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	DataMap   map[string]any
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}
