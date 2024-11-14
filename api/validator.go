package api

import (
	"github.com/RahmounFares2001/simple-bank-backend/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// check if currency supported
		return util.IsSupportedCurrency(currency)
	}
	return false
}
