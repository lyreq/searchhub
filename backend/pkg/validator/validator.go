package validator

import (
	"errors"
	"fmt"
	"reflect"

	"io"
	"mime/multipart"

	"github.com/Lexographics/go-openapigen"
	"github.com/Lexographics/go-postmangen"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var GenerateDocs = false

type Validator interface {
	ParseAndValidate(c echo.Context, o interface{}) error
}

type validatorImpl struct {
	validate   *validator.Validate
	pggen      *postmangen.PostmanGen
	openapigen *openapigen.OpenAPIGen
}

func New(db *gorm.DB, docgen *postmangen.PostmanGen, openapigen *openapigen.OpenAPIGen) Validator {
	v := validator.New()
	RegisterCustomValidators(v, db)
	return &validatorImpl{
		validate:   v,
		pggen:      docgen,
		openapigen: openapigen,
	}
}

func (v *validatorImpl) ParseAndValidate(c echo.Context, o interface{}) error {
	if GenerateDocs {
		data := map[string]any{
			"method":    c.Request().Method,
			"path":      c.Request().URL.Path,
			"inputType": reflect.TypeOf(o),
		}

		err := v.openapigen.Register(data)
		if err != nil {
			return err
		}

		err = v.pggen.Register(data)
		if err != nil {
			return err
		}
		return errors.New("docs generated")
	}

	if err := c.Bind(o); err != nil {
		return err
	}

	if err := v.LocalsParser(c, o); err != nil {
		return err
	}

	if err := v.FormFileParser(c, o); err != nil {
		return err
	}

	if err := v.formArrayParser(c, o); err != nil {
		return err
	}

	// Validate the struct
	if err := v.validate.Struct(o); err != nil {
		return err
	}

	return nil
}

func (v *validatorImpl) LocalsParser(c echo.Context, o interface{}) error {
	val := reflect.ValueOf(o)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("invalid object")
	}

	val = val.Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("locals")
		if tag == "" {
			continue
		}

		value := c.Get(tag)
		if value == nil {
			continue
		}

		fieldValue := val.Field(i)
		if !fieldValue.CanSet() {
			return fmt.Errorf("field %s is not settable", field.Name)
		}

		localValue := reflect.ValueOf(value)
		if localValue.Type().ConvertibleTo(fieldValue.Type()) {
			fieldValue.Set(localValue.Convert(fieldValue.Type()))
		} else {
			return fmt.Errorf("invalid type '%s' for field %s", localValue.Type().String(), field.Name)
		}
	}

	return nil
}

func (v *validatorImpl) FormFileParser(c echo.Context, o interface{}) error {
	val := reflect.ValueOf(o)
	if val.Kind() != reflect.Ptr {
		return errors.New("must be ptr")
	}

	if val.IsNil() {
		return errors.New("must be non-nil")
	}

	val = val.Elem()
	typ := val.Type()
	if val.Kind() != reflect.Struct {
		return errors.New("must be pointer to a struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("formFile")
		if tag == "" || tag == "-" {
			continue
		}

		fieldValue := val.Field(i)
		if !fieldValue.CanSet() {
			return fmt.Errorf("field %s is not settable", field.Name)
		}

		if fieldValue.Type() == reflect.PointerTo(reflect.TypeOf(multipart.FileHeader{})) {
			file, err := c.FormFile(tag)
			if err != nil {
				continue
			}
			fieldValue.Set(reflect.ValueOf(file))
			continue
		}

		if fieldValue.Type() == reflect.TypeOf([]*multipart.FileHeader{}) {
			form, err := c.MultipartForm()
			if err != nil {
				continue
			}

			files, ok := form.File[tag]
			if !ok {
				continue
			}
			fieldValue.Set(reflect.ValueOf(files))
			continue
		}

		if fieldValue.Type() == reflect.TypeOf([]byte{}) || fieldValue.Type() == reflect.TypeOf(string("")) {
			file, err := c.FormFile(tag)
			if err != nil {
				continue
			}

			f, err := file.Open()
			if err != nil {
				continue
			}
			defer f.Close()

			fileBytes, err := io.ReadAll(f)
			if err != nil {
				continue
			}

			if fieldValue.Type() == reflect.TypeOf(string("")) {
				fieldValue.SetString(string(fileBytes))
			} else {
				fieldValue.Set(reflect.ValueOf(fileBytes))
			}
			continue
		}

		if fieldValue.Type() == reflect.TypeOf([][]byte{}) || fieldValue.Type() == reflect.TypeOf([]string{}) {
			form, err := c.MultipartForm()
			if err != nil {
				continue
			}

			files, ok := form.File[tag]
			if !ok {
				continue
			}

			for _, file := range files {
				f, err := file.Open()
				if err != nil {
					continue
				}
				defer f.Close()

				fileBytes, err := io.ReadAll(f)
				if err != nil {
					continue
				}
				if fieldValue.Type() == reflect.TypeOf([]string{}) {
					fieldValue.Set(reflect.Append(fieldValue, reflect.ValueOf(string(fileBytes))))
				} else {
					fieldValue.Set(reflect.Append(fieldValue, reflect.ValueOf(fileBytes)))
				}
			}
			continue
		}
	}

	return nil
}
