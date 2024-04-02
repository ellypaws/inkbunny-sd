package utils

import (
	"errors"
	"fmt"
	"github.com/ellypaws/inkbunny/api"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
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

func HasTxtFile(submission api.Submission) bool {
	for _, file := range submission.Files {
		if strings.HasPrefix(file.MimeType, "text") {
			return true
		}
	}
	return false
}

func HasJsonFile(submission api.Submission) bool {
	for _, file := range submission.Files {
		if strings.HasSuffix(file.MimeType, "json") {
			return true
		}
	}
	return false
}

func HasMetadata(submission api.Submission) bool {
	for _, file := range submission.Files {
		if strings.HasPrefix(file.MimeType, "text") {
			return true
		}
		if strings.HasSuffix(file.MimeType, "json") {
			return true
		}
	}
	return false
}

func FilterMetadata(submission api.Submission) (files []api.File) {
	for _, file := range submission.Files {
		if strings.HasPrefix(file.MimeType, "text") {
			files = append(files, file)
		}
		if strings.HasSuffix(file.MimeType, "json") {
			files = append(files, file)
		}
	}
	return files
}

type MetadataContent struct {
	Blob []byte
	api.File
}

func GetMetadataBytes(submission api.Submission) ([]MetadataContent, error) {
	var metadata []MetadataContent
	for _, file := range submission.Files {
		if !strings.HasPrefix(file.MimeType, "text") && !strings.HasSuffix(file.MimeType, "json") {
			continue
		}
		if file.FileURLFull == "" {
			continue
		}
		r, err := http.Get(file.FileURLFull)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()
		if r.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %v", r.Status)
		}
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		metadata = append(metadata, MetadataContent{b, file})
	}
	return metadata, nil
}
