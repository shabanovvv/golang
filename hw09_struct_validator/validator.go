package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("field '%s': %v", v.Field, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errs := make([]string, 0, len(v))
	for _, e := range v {
		errs = append(errs, fmt.Sprintf("field '%s': %s", e.Field, e.Err))
	}
	return strings.Join(errs, ", ")
}

func Validate(v interface{}) error {
	up := UtilPars{}
	var validationErrors ValidationErrors
	uc := UtilCheck{}
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return errors.New("struct type failed")
	}

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := strings.TrimSpace(typeField.Tag.Get("validate"))
		if tag == "" {
			continue
		}

		up.rr = RuleResult{}
		err := up.ParseTag(typeField.Type, tag)
		if err != nil {
			return err
		}

		valid := uc.CheckTag(valueField, up.rr)
		if !valid {
			validationErrors = append(validationErrors, ValidationError{
				Field: typeField.Name,
				Err:   fmt.Errorf("field %s is invalid", typeField.Name),
			})
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}
