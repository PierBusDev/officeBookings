package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

//Form embeds url.Values object + errors
type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks if the form submitted by the user has a field named field and it is not empty
func (f *Form) Has(field string, r *http.Request) bool {
	value := r.Form.Get(field)
	if value == "" {
		return false
	}
	return true
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

//Valid returns true if there are no errors, false otherwise
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

//MinLegth checks that the length of the content of field is at least length
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	fieldValue := r.Form.Get(field)
	if len(fieldValue) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}

	return true
}

//IsEmail checks for valid emails
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
