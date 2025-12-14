package validator

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/validate"
)

func FormatErrors(errors *validate.Errors) string {
	if errors == nil || len(errors.Errors) == 0 {
		return ""
	}
	errorMessages := make([]string, 0, len(errors.Errors))
	for fieldName, errorsPerField := range errors.Errors {
		if len(errorsPerField) == 0 {
			continue
		}
		errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", fieldName, strings.Join(errorsPerField, ", ")))
	}
	return strings.Join(errorMessages, "\n")
}
