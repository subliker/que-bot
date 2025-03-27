package validation

import (
	"github.com/go-playground/validator/v10"
)

// Validate is global validator, can be specified with custom registered validators
var Validate = validator.New()
