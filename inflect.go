package inflect

import (
	"errors"
	"reflect"
)

//-----------------------------------------------------------------------------
// utilities

func valueOf(data interface{}) *reflect.Value {
	var val reflect.Value

	switch reflect.TypeOf(data).Kind() {
	case reflect.Ptr:
		val = reflect.ValueOf(data).Elem()
	case reflect.Struct:
		val = reflect.ValueOf(data)
	default:
		return nil
	}

	return &val
}

func ptrOf(data interface{}) *reflect.Value {
	var val reflect.Value

	switch reflect.TypeOf(data).Kind() {
	case reflect.Ptr:
		val = reflect.ValueOf(data).Elem()
	default:
		return nil
	}

	return &val
}

// Errors
var (
	ErrNotFound     = errors.New("field does not exist")
	ErrNotMatched   = errors.New("CAS old value does not match")
	ErrInvalidType  = errors.New("data has invalid type - not a pointer nor a struct")
	ErrNonPointer   = errors.New("a pointer was expected")
	ErrNoSet        = errors.New("can not set field value")
	ErrTypeMismatch = errors.New("field type & value type are different")
)

//-----------------------------------------------------------------------------
// API

// Get extracts the value of a field
func Get(data interface{}, fieldName string) (interface{}, error) {
	val := valueOf(data)
	if val == nil {
		return nil, ErrInvalidType
	}
	field := val.FieldByName(fieldName)
	if field.IsValid() {
		return field.Interface(), nil
	}
	return nil, ErrNotFound
}

// Set accepts a pointer and sets the field value
func Set(data interface{}, fieldName string, value interface{}) error {
	ptr := ptrOf(data)
	if ptr == nil {
		return ErrNonPointer
	}
	field := ptr.FieldByName(fieldName)
	if !field.IsValid() {
		return ErrNotFound
	}
	if !field.CanSet() {
		return ErrNoSet
	}
	newValue := reflect.ValueOf(value)
	if newValue.Type() != field.Type() {
		return ErrTypeMismatch
	}
	field.Set(newValue)
	return nil
}

// CAS compares field's value with an old value and sets it to new value if it matches
func CAS(data interface{}, fieldName string, oldValue, newValue interface{}) error {
	ptr := ptrOf(data)
	if ptr == nil {
		return ErrNonPointer
	}
	field := ptr.FieldByName(fieldName)
	if !field.IsValid() {
		return ErrNotFound
	}
	if !field.CanSet() {
		return ErrNoSet
	}
	_newValue := reflect.ValueOf(newValue)
	_oldValue := reflect.ValueOf(oldValue)
	ft := field.Type()
	if _oldValue.Type() != ft || _newValue.Type() != ft {
		return ErrTypeMismatch
	}
	if field.Interface() != oldValue {
		return ErrNotMatched
	}
	field.Set(_newValue)
	return nil
}

// Tag extract a tag & returns it's values for a field
func Tag(data interface{}, fieldName string, tag string) (string, error) {
	val := valueOf(data)
	if val == nil {
		return "", ErrInvalidType
	}
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return "", ErrNotFound
	}
	fieldType, _ := val.Type().FieldByName(fieldName)
	tv, _ := fieldType.Tag.Lookup(tag)
	if tv == "" {
		return "", ErrNotFound
	}
	return tv, nil
}
