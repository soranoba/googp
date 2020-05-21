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
	values map[string]accessor
}

func newAccessor(tag *tag, v reflect.Value) accessor {
	switch reflect.Indirect(v).Kind() {
	case reflect.Array, reflect.Slice:
		return &arrayAccessor{tag: tag, value: reflect.Indirect(v)}
	case reflect.Struct:
		sv := reflect.Indirect(v)
		t := sv.Type()
		values := make(map[string]accessor)
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			// NOTE: private field
			if !isUpperPrefix(f.Name) {
				continue
			}

			fieldValue := sv.Field(i)
			if !fieldValue.CanSet() {
				continue
			}

			tag := newTag(f)
			for _, name := range tag.names {
				if _, ok := values[name]; !ok {
					values[name] = newAccessor(tag, fieldValue)
				}
			}
		}

		if len(values) == 0 {
			return newValueAccessor(v)
		} else {
			return &structAccessor{values: values}
		}
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

func (f *structAccessor) Set(key string, val string) error {
	parts := strings.Split(key, ":")
	for i := len(parts); i >= 1; i-- {
		k := strings.Join(parts[0:i], ":")
		if v := f.values[k]; v != nil {
			return v.Set(key, val)
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
