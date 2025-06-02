package domain

import "github.com/go-playground/validator"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ResponseURL struct {
	URL   string `json:"url"`
	Alias string `json:"alias"`
}

type RequestURL struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias"`
}

func (r *RequestURL) Validate() error {
	return validate.Struct(r)
}
