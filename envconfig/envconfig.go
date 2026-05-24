// Package envconfig loads environment variables into struct fields using tags.
// It supports required fields, default values, and common scalar types.
package envconfig

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Load parses environment variables into the struct pointed to by dst. Each
// exported field may carry an `env` tag:
//
//	`env:"VAR_NAME"`                 // use VAR_NAME
//	`env:"VAR_NAME,required"`        // must be set (env or default)
//	`env:"VAR_NAME,default=value"`   // fallback value
//
// Load returns an error if dst is not a pointer to a struct, if a required
// variable is missing, or if a value cannot be converted to the field type.
func Load(dst any, opts ...Option) error {
	c := &config{}
	for _, fn := range opts {
		fn(c)
	}
	return c.load(dst)
}

func (c *config) load(dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("envconfig: dst must be a non-nil pointer to struct, got %T", dst)
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("envconfig: dst must be a non-nil pointer to struct, got %T", dst)
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		if !ft.IsExported() {
			continue
		}
		fv := v.Field(i)
		tag, ok := ft.Tag.Lookup("env")
		if !ok || tag == "-" {
			continue
		}
		spec := parseTag(tag)
		if spec.name == "" {
			spec.name = strings.ToUpper(ft.Name)
		}

		varName := spec.name
		if c.prefix != "" {
			varName = c.prefix + spec.name
		}

		val, found := os.LookupEnv(varName)
		if !found {
			if spec.def != "" {
				val = spec.def
			} else if spec.required {
				return &MissingError{VarName: varName}
			} else {
				continue
			}
		}

		if err := setField(fv, val, varName); err != nil {
			return err
		}
	}
	return nil
}

type tagSpec struct {
	name     string
	required bool
	def      string
}

func parseTag(tag string) tagSpec {
	parts := strings.Split(tag, ",")
	spec := tagSpec{name: strings.TrimSpace(parts[0])}
	for _, p := range parts[1:] {
		p = strings.TrimSpace(p)
		switch {
		case p == "required":
			spec.required = true
		case strings.HasPrefix(p, "default="):
			spec.def = p[len("default="):]
		}
	}
	return spec
}

func setField(fv reflect.Value, val, varName string) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(val)

	case reflect.Bool:
		v, err := strconv.ParseBool(val)
		if err != nil {
			return &ParseError{VarName: varName, Value: val, Err: err}
		}
		fv.SetBool(v)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(val, 0, fv.Type().Bits())
		if err != nil {
			return &ParseError{VarName: varName, Value: val, Err: err}
		}
		fv.SetInt(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(val, 0, fv.Type().Bits())
		if err != nil {
			return &ParseError{VarName: varName, Value: val, Err: err}
		}
		fv.SetUint(v)

	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(val, fv.Type().Bits())
		if err != nil {
			return &ParseError{VarName: varName, Value: val, Err: err}
		}
		fv.SetFloat(v)

	default:
		return &ParseError{
			VarName: varName,
			Value:   val,
			Err:     fmt.Errorf("unsupported type %s", fv.Type()),
		}
	}
	return nil
}
