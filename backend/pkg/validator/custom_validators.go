package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// RegisterCustomValidators registers all custom validators
func RegisterCustomValidators(v *validator.Validate, db *gorm.DB) {
	v.RegisterValidation("exists", func(fl validator.FieldLevel) bool {
		// Get the field value
		fieldValue := fl.Field().Interface()

		// Skip validation if field is empty (use required tag for that)
		if fieldValue == nil || (reflect.ValueOf(fieldValue).Kind() == reflect.String && fieldValue.(string) == "") {
			return true
		}

		// Get the parameters from the tag (format: exists:table,column)
		params := strings.Split(fl.Param(), ".")
		if len(params) != 2 {
			return false
		}

		tableName := params[0]
		columnName := params[1]

		// Create a dynamic query
		var count int64
		result := db.Table(tableName).Where(columnName+" = ?", fieldValue).Count(&count)

		if result.Error != nil {
			return false
		}

		return count > 0
	})
}
