package validator

import (
	"net/url"
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) CheckURL(urlString string, key string) {
	if _, err := url.ParseRequestURI(urlString); err != nil {
		v.AddError(key, "must be a valid URL")
	}
}

func PermittedValues(values []string, permitted []string) (string, bool) {
	for _, value := range values {
		isPermitted := false
		for _, permittedValue := range permitted {
			if value == permittedValue {
				isPermitted = true
				break
			}
		}
		if !isPermitted {
			return value, false
		}
	}

	return "", true
}
