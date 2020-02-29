package crepe

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const tagName = "crepe"

func Unmarshal(b []byte, i interface{}) error {
	reader := bytes.NewReader(b)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return fmt.Errorf("failed to parse HTML document: %w", err)
	}
	return unmarshal(doc.Selection, i)
}

func unmarshal(s *goquery.Selection, i interface{}) error {
	rValueOf := reflect.ValueOf(i)
	if rValueOf.IsNil() {
		return fmt.Errorf("%w: is nil", ErrInvalidValue)
	}
	switch rValueOf.Kind() {
	case reflect.Ptr:
		return unmarshalPtr(s, rValueOf)
	}
	return fmt.Errorf("%w: not a ptr", ErrInvalidValue)
}

func unmarshalPtr(s *goquery.Selection, rValueOf reflect.Value) error {
	if rValueOf.Kind() != reflect.Ptr {
		return fmt.Errorf("%w: not a ptr", ErrInvalidValue)
	}

	switch rValueOf := rValueOf.Elem(); rValueOf.Kind() {
	case reflect.Struct:
		return unmarshalStruct(s, rValueOf)
	}
	return fmt.Errorf("%w: not a ptr to struct", ErrInvalidValue)
}

func unmarshalStruct(s *goquery.Selection, rValueOf reflect.Value) error {
	if rValueOf.Kind() != reflect.Struct {
		return fmt.Errorf("%w: not a struct", ErrInvalidValue)
	}

	// Iterate over all available fields and read the tag value
	for i := 0; i < rValueOf.NumField(); i++ {
		rType, rValue := rValueOf.Type().Field(i), rValueOf.Field(i)

		var err error
		switch kind := rValue.Kind(); kind {
		// Primitive type
		case reflect.Bool, reflect.String, reflect.Float32, reflect.Float64,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			err = unmarshalPrimitive(s, rValue, rType)

		case reflect.Slice:
			err = unmarshalSlice(s, rValue, rType)

		// Map is not supported
		case reflect.Map, reflect.Interface:
			err = fmt.Errorf("%w: %v", ErrUnimplemented, kind)

		case reflect.Ptr:
			err = unmarshalPtr(s, rValue)

		// TODO
		case reflect.Array:
			fallthrough

		// TODO
		case reflect.Struct:
			fallthrough

		default:
			err = fmt.Errorf("%w: %v", ErrUnimplemented, kind)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func unmarshalSlice(s *goquery.Selection, rValueOf reflect.Value, rField reflect.StructField) error {
	tag := rField.Tag.Get(tagName)
	if tag == "" {
		return nil // Ignore this field
	}

	if !rValueOf.CanSet() {
		return ErrUnexportedField
	}

	statement, err := newDecoderFromTag(tag)
	if err != nil {
		return fmt.Errorf("failed to parse tag: %w", err)
	}
	s = statement.Delegate(s)

	slice := reflect.MakeSlice(rValueOf.Type(), rValueOf.Len(), rValueOf.Cap())
	for i := 0; i < s.Length(); i++ {
		var elem reflect.Value
		switch typ := rValueOf.Type().Elem(); typ.Kind() {
		case reflect.Ptr:
			elem = reflect.New(typ.Elem())
			if err := unmarshalStruct(s, elem.Elem()); err != nil {
				return err
			}
		case reflect.Struct:
			elem = reflect.New(typ).Elem()
			if err := unmarshalStruct(s, elem); err != nil {
				return err
			}
		default:
			return ErrUnsupportedType
		}
		slice = reflect.Append(slice, elem)
	}
	rValueOf.Set(slice)
	return nil
}

func unmarshalPrimitive(s *goquery.Selection, rValueOf reflect.Value, rField reflect.StructField) error {
	tag := rField.Tag.Get(tagName)
	if tag == "" {
		return nil // Ignore this field
	}

	if !rValueOf.CanSet() {
		return ErrUnexportedField
	}

	statement, err := newDecoderFromTag(tag)
	if err != nil {
		return fmt.Errorf("failed to parse tag: %w", err)
	}

	result, err := statement.Decode(s)
	if err != nil {
		return err
	}

	switch rField.Type.Kind() {
	case reflect.String:
		rValueOf.SetString(result.raw)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(result.raw, 0, 64)
		if err != nil {
			return fmt.Errorf("failed to parse %v: %w", result.raw, err)
		}
		rValueOf.SetInt(v)

	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(result.raw, 64)
		if err != nil {
			return fmt.Errorf("failed to parse %v: %w", result.raw, err)
		}
		rValueOf.SetFloat(v)

	case reflect.Bool:
		v, err := strconv.ParseBool(result.raw)
		if err != nil {
			return fmt.Errorf("failed to parse %v: %w", result.raw, err)
		}
		rValueOf.SetBool(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(result.raw, 0, 64)
		if err != nil {
			return fmt.Errorf("failed to parse %v: %w", result.raw, err)
		}
		rValueOf.SetUint(v)

	default:
		return fmt.Errorf("%w: not a primitive type: %v", ErrInvalidValue, rField.Type.Kind())
	}
	return nil
}
