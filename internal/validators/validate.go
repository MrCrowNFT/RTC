package validators

import(

)

var validate = validator.New()

func formatValidationErrors(err error)(string){
	var sb strings.Builder
	for _, err :=  range err.(validator.ValidationErrors){
		sb.WriteString(fmt.Sprintf("Field '%s' failed validation: %s.", err.Field(), err.Tag()))
	}
	return sb.String()
}