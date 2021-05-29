package data

import (
	"regexp"
	"time"

	"github.com/go-playground/validator"
	"github.com/kjunn2000/currency/"
)

// Product defines the structure for an API product
// swagger:model
type Product struct {
	Id        int       `json:"Id" validate:"required"`
	Name      string    `json:"Name" validate:"required,id"`
	Price     int       `json:"Price" validate:"gte=10"`
	CreatedAt time.Time `json:"-"`
}

type Products []Product

func (pro *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("id", idValidate)
	return validate.Struct(pro)
}

func idValidate(f1 validator.FieldLevel) bool {
	reg := regexp.MustCompile(`[a-z]+`)
	matches := reg.FindAllString(f1.Field().String(), -1)
	return len(matches) == 1
}
