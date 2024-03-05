package main

import (
	"errors"
	"reflect"
	"strconv"
)

// castStringToType dynamically casts a string to a field's type and assigns it.
func castStringToType(s string, fieldPtr any) error {
	if s == "" {
		return nil
	}
	// Ensure fieldPtr is indeed a pointer, to be able to modify the original data
	ptrValue := reflect.ValueOf(fieldPtr)
	if ptrValue.Kind() != reflect.Ptr {
		return errors.New("fieldPtr must be a pointer")
	}

	// Dereference the pointer to work with the actual value
	value := ptrValue.Elem()

	switch value.Kind() {
	case reflect.Int, reflect.Int64:
		// Convert string to int
		if intValue, err := strconv.ParseInt(s, 10, 64); err == nil {
			value.SetInt(intValue)
		} else {
			return err
		}
	case reflect.Float64:
		// Convert string to float64
		if floatValue, err := strconv.ParseFloat(s, 64); err == nil {
			value.SetFloat(floatValue)
		} else {
			return err
		}
	case reflect.String:
		// Assign the string directly
		value.SetString(s)
	case reflect.Ptr:
		// For pointer types, you would need more specialized logic depending on what the pointer is pointing to.
		// This is a simplistic example for *string types.
		if value.Type().Elem().Kind() == reflect.String {
			value.Set(reflect.ValueOf(&s))
		}
	default:
		return errors.New("unsupported field type")
	}

	return nil
}
