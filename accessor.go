package googp

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// accessor is an interface for writing the value of ogp to variables.
type accessor interface {
	Set(key string, val string) error
}

// valueAccessor is an accessor for writing the value of ogp to single variable.
type valueAccessor struct {
	value  reflect.Value
	didSet bool
}

// arrayAccessor is an accessor for writing the values of ogp to an array or a slice.
type arrayAccessor struct {
	tag     *tag
	value   reflect.Value
	idx     int
	current accessor
}

// structAccessor is an accessor for writing the values of ogp to a struct.
type structAccessor struct {
	value  *reflect.Value
	fields map[string]*field
}

type field struct {
	structField *reflect.StructField
	tag         *tag
	value       *reflect.Value
	accessor    accessor
}

func newAccessor(tag *tag, v reflect.Value) accessor {
	iv := reflect.Indirect(v)
	switch iv.Kind() {
	case reflect.Array, reflect.Slice:
		return &arrayAccessor{tag: tag, value: iv}
	case reflect.Struct:
		if iv.CanAddr() {
			if iv.Addr().Type().Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()) {
				return newValueAccessor(v)
			}
		}

		t := iv.Type()
		fields := make(map[string]*field)
		for i := 0; i < t.NumField(); i++ {
			structField := t.Field(i)
			fieldValue := iv.Field(i)
			if !fieldValue.CanSet() {
				continue
			}

			tag := newTag(structField)
			field := &field{
				structField: &structField,
				tag:         newTag(structField),
				value:       &fieldValue,
			}

			for _, name := range tag.names {
				if _, ok := fields[name]; !ok {
					fields[name] = field
				}
			}
		}

		if len(fields) == 0 {
			return newValueAccessor(v)
		}
		return &structAccessor{value: &v, fields: fields}
	default:
		return newValueAccessor(v)
	}
}

func newValueAccessor(v reflect.Value) *valueAccessor {
	return &valueAccessor{value: v}
}

func (f *valueAccessor) Set(key string, val string) error {
	// NOTE: The first tag (from top to bottom) is given preference during conflicts.
	if f.didSet {
		return nil
	}

	if !f.value.IsValid() {
		return fmt.Errorf("invalid reflect.Value")
	}

	ty := f.value.Type()
	v := f.value
	if f.value.Kind() == reflect.Ptr {
		if f.value.IsNil() && f.value.CanSet() {
			v = reflect.New(f.value.Type().Elem())
			f.value.Set(v)
			v = v.Elem()
		} else {
			v = reflect.Indirect(v)
		}
	}
	if !v.CanSet() {
		return fmt.Errorf("Cannot set to value")
	}

	switch v.Kind() {
	case reflect.String:
		v.Set(reflect.ValueOf(val).Convert(v.Type()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil || reflect.Zero(v.Type()).OverflowInt(i) {
			return convertErr(key, val, ty)
		}
		v.Set(reflect.ValueOf(i).Convert(v.Type()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		u, err := strconv.ParseUint(val, 10, 64)
		if err != nil || reflect.Zero(v.Type()).OverflowUint(u) {
			return convertErr(key, val, ty)
		}
		v.Set(reflect.ValueOf(u).Convert(v.Type()))
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(val, 64)
		if err != nil || reflect.Zero(v.Type()).OverflowFloat(n) {
			return convertErr(key, val, ty)
		}
		v.Set(reflect.ValueOf(n).Convert(v.Type()))
	default:
		unmarshaler, ok := v.Addr().Interface().(encoding.TextUnmarshaler)
		if !ok {
			return unsupportedErr(key, ty)
		}
		if err := unmarshaler.UnmarshalText([]byte(val)); err != nil {
			return convertErr(key, val, ty)
		}
	}

	f.didSet = true
	return nil
}

func (f *arrayAccessor) Set(key string, val string) error {
	if f.current != nil && f.tag != nil && !f.tag.isContainsName(key) {
		return f.current.Set(key, val)
	}

	if f.current != nil {
		f.idx += 1
	}

	if f.idx >= f.value.Len() {
		if f.value.Kind() == reflect.Slice {
			if !f.value.CanSet() {
				return fmt.Errorf("Cannot set to value")
			}
			v := reflect.New(f.value.Type().Elem()).Elem()
			f.value.Set(reflect.Append(f.value, v))
		} else {
			return nil
		}
	}

	f.current = newAccessor(nil, f.value.Index(f.idx))
	return f.current.Set(key, val)
}

func (ac *structAccessor) Set(key string, val string) error {
	parts := strings.Split(key, ":")
	for i := len(parts); i >= 0; i-- {
		k := strings.Join(parts[0:i], ":")
		if f := ac.fields[k]; f != nil {
			if f.accessor == nil {
				if f.value.Kind() == reflect.Ptr && f.value.IsNil() {
					f.value.Set(reflect.New(f.value.Type().Elem()))
				}
				f.accessor = newAccessor(f.tag, *f.value)
			}
			return f.accessor.Set(key, val)
		}
	}
	return nil
}

func convertErr(key string, val string, ty reflect.Type) error {
	return fmt.Errorf("%s field is invalid. (type = %s, value = %s)", key, ty.Name(), val)
}

func unsupportedErr(key string, ty reflect.Type) error {
	return fmt.Errorf("%s is unsupported type (field = %s)", ty.Name(), key)
}
