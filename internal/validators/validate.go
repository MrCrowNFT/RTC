package validators

import(
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator struct wraps the go-playground validator for extensibility.
type Validator struct {
	validate *validator.Validate
}

// NewValidator initializes and returns a Validator instance.
func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func FormatValidationErrors(err error)(string){
	var sb strings.Builder
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			sb.WriteString(fmt.Sprintf("Field '%s' failed validation: %s.\n", fe.Field(), fe.Tag()))
		}
	} else {
		sb.WriteString("Invalid input data.\n")
	}
	return sb.String()
}

// ValidateStruct validates a struct and returns an error if validation fails.
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}