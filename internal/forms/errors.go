package forms

type errors map[string][]string

// Add an error message to a given field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error message
func (e errors) Get(field string) string {
	errorString := e[field]
	if len(errorString) == 0 {
		return ""
	}
	return errorString[0] //returns first error (because anyway we will want to show just one)
}
