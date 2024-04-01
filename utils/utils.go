package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// ResultsToFields maps ExtractResult from ExtractAll to fields in a struct.
func ResultsToFields(result ExtractResult, fieldsToSet map[string]any) error {
	for key, fieldPtr := range fieldsToSet {
		if v, ok := result[key]; ok {
			err := CastStringToType(v, fieldPtr)
			if err != nil {
				return fmt.Errorf("error casting %s to type: %w", key, err)
			}
		}
	}
	return nil
}

// CastStringToType dynamically casts a string to a field's type and assigns it.
func CastStringToType(s string, fieldPtr any) error {
	if s == "" {
		return nil
	}
	// Ensure fieldPtr is indeed a pointer, to be able to modify the original data
	ptrValue := reflect.ValueOf(fieldPtr)
	if ptrValue.Kind() != reflect.Pointer {
		return errors.New("fieldPtr must be a pointer")
	}

	// Dereference the pointer to work with the actual value
	value := ptrValue.Elem()

	switch value.Kind() {
	case reflect.Int, reflect.Int64:
		if intValue, err := strconv.ParseInt(s, 10, 64); err == nil {
			value.SetInt(intValue)
		} else {
			return err
		}
	case reflect.Float64:
		if floatValue, err := strconv.ParseFloat(s, 64); err == nil {
			value.SetFloat(floatValue)
		} else {
			return err
		}
	case reflect.String:
		value.SetString(s)
	case reflect.Bool:
		if boolValue, err := strconv.ParseBool(s); err == nil {
			value.SetBool(boolValue)
		} else {
			return err
		}
	case reflect.Pointer:
		return CastStringToType(s, value.Addr().Interface())
	default:
		return errors.New("unsupported field type")
	}

	return nil
}
