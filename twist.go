package twist

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/go-ini/ini"
	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
)

// Tag name constants
const (
	tagNameDefault = "default"
	tagNameToml    = "toml"
	tagNameIni     = "ini"
	tagNameYaml    = "yaml"
	tagNameJson    = "json"
	tagNameEnv     = "env"
	tagNameCli     = "cli"
)

var isDebug = os.Getenv("TWIST_DEBUG") != ""

// Debug function
// If CC_DEBUG environment is defined, output some logs
func debug(args ...interface{}) {
	if !isDebug {
		return
	}
	fmt.Println(args...)
}

// Dereference reflect.Value
func derefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// Dereference reflect.Type
func derefType(v reflect.Type) reflect.Type {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// Main cascading function
// Note that opts order is important. Configraions will be overrided by options order.
// For example:
//  Mix(v, WithToml(), WithJson())            cascade order is toml -> json
//  Mix(v, WithToml(), WithJson(), WithEnv()) cascade order is toml -> jsoa -> env
func Mix(v interface{}, opts ...Option) error {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		return errors.New("Cascading value must be a struct")
	}
	value := derefValue(reflect.ValueOf(v))
	if !value.CanSet() {
		return errors.New("destination value cannot set values")
	}

	for _, opt := range opts {
		switch opt.name {
		case optionNameToml:
			if err := cascadeToml(opt.value.(string), value, reflect.New(t)); err != nil {
				return errors.Wrap(err, "Failed to cascade toml")
			}
		case optionNameYaml:
			if err := cascadeYaml(opt.value.(string), value, reflect.New(t)); err != nil {
				return errors.Wrap(err, "Failed to cascade yaml")
			}
		case optionNameIni:
			src, err := ini.Load(opt.value.(string))
			if err != nil {
				return errors.Wrap(err, "ini load error")
			}
			if err := cascadeIni(src, src.Section(""), value); err != nil {
				return errors.Wrap(err, "Failed to cascade ini")
			}
		case optionNameJson:
			if err := cascadeJson(opt.value.(string), value, reflect.New(t)); err != nil {
				return errors.Wrap(err, "Failed to cascade json")
			}
		case optionNameEnv:
			if err := cascadeEnv(value); err != nil {
				return errors.Wrap(err, "Failed to cascade env")
			}
		case optionNameCli:
			if err := cascadeCli(value, parseCliArgs(opt.value.([]string))); err != nil {
				return errors.Wrap(err, "Failed to cascase cli")
			}
		}
	}
	if err := cascadeDefault(value); err != nil {
		return errors.Wrap(err, "failed to set default value")
	}
	return nil
}

// Parse toml file and merge to base struct
func cascadeToml(file string, base, clone reflect.Value) error {
	if _, err := toml.DecodeFile(file, clone.Interface()); err != nil {
		return errors.Wrap(err, "toml decode error")
	}
	return mergeConfig(base, derefValue(clone), tagNameToml)
}

// Parse yaml file and merge to base struct
func cascadeYaml(file string, base, clone reflect.Value) error {
	if _, err := os.Stat(file); err != nil {
		return errors.Wrap(err, "yaml file not exists")
	}
	fp, err := os.Open(file)
	if err != nil {
		return errors.Wrap(err, "yaml file open error")
	}
	defer fp.Close()
	if err := yaml.NewDecoder(fp).Decode(clone.Interface()); err != nil {
		return errors.Wrap(err, "yaml decode error")
	}
	return mergeConfig(base, derefValue(clone), tagNameYaml)
}

// Parse JSON file and merge to base struct
func cascadeJson(file string, base, clone reflect.Value) error {
	if _, err := os.Stat(file); err != nil {
		return errors.Wrap(err, "json file not exists")
	}
	fp, err := os.Open(file)
	if err != nil {
		return errors.Wrap(err, "json file open error")
	}
	defer fp.Close()
	if err := json.NewDecoder(fp).Decode(clone.Interface()); err != nil {
		return errors.Wrap(err, "json decode error")
	}
	return mergeConfig(base, derefValue(clone), tagNameJson)
}

// Find INI section value and merge to base struct
// Note that currently we support only single section, so you can't define nested section.
func cascadeIni(cfg *ini.File, s *ini.Section, v reflect.Value) error {
	t := derefType(v.Type())
	v = derefValue(v)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if !value.CanSet() {
			debug("cannot set: ", field.Name)
			continue
		}

		ft := field.Type
		var isPtr bool
		if ft.Kind() == reflect.Ptr {
			isPtr = true
			ft = derefType(ft)
		}

		tag, ok := field.Tag.Lookup(tagNameIni)
		if !ok || tag == "" || tag == "-" {
			continue
		}
		if ft.Kind() == reflect.Struct {
			if ss := cfg.Section(tag); ss != nil {
				debug("subsection: ", tag, ss.KeyStrings())
				cascadeIni(cfg, ss, value)
			}
			continue
		}
		key := s.Key(tag)
		if key == nil {
			debug("key not found in ini: ", tag)
			continue
		}
		debug(field.Name, key.String(), ft.Kind())
		if err := assignValue(ft, value, key.Value(), isPtr, false); err != nil {
			return errors.Wrap(err, "failed to assign values")
		}
		debug("assigned: ", tag, key.Value())
	}
	return nil
}

// Walk struct field and assign from environment variable
func cascadeEnv(v reflect.Value) error {
	t := derefType(v.Type())
	v = derefValue(v)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if !value.CanSet() {
			debug("cannot set: ", field.Name)
			continue
		}

		ft := field.Type
		var isPtr bool
		if ft.Kind() == reflect.Ptr {
			isPtr = true
			ft = derefType(ft)
		}

		if ft.Kind() == reflect.Struct {
			if field.Tag == "" {
				continue
			}
			cascadeEnv(value)
			continue
		}
		tag, ok := field.Tag.Lookup(tagNameEnv)
		if !ok || tag == "" || tag == "-" {
			continue
		}
		envValue := os.Getenv(tag)
		if envValue == "" {
			continue
		}
		debug(field.Name, envValue, ft.Kind())
		if err := assignValue(ft, value, envValue, isPtr, false); err != nil {
			return errors.Wrap(err, "failed to assign values")
		}
		debug("assigned: ", field.Name, envValue)
	}
	return nil
}

// Walk struct field and assign from default tagged value
func cascadeDefault(v reflect.Value) error {
	t := derefType(v.Type())
	v = derefValue(v)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if !value.CanSet() {
			debug("cannot set: ", field.Name)
			continue
		}

		ft := field.Type
		var isPtr bool
		if ft.Kind() == reflect.Ptr {
			isPtr = true
			ft = derefType(ft)
		}

		if ft.Kind() == reflect.Struct {
			if field.Tag == "" {
				continue
			}
			cascadeDefault(value)
			continue
		}
		if !isZeroValue(ft, value) {
			continue
		}
		tag, ok := field.Tag.Lookup(tagNameDefault)
		if !ok || tag == "" || tag == "-" {
			continue
		}
		if err := assignValue(ft, value, tag, isPtr, false); err != nil {
			return errors.Wrap(err, "failed to assign values")
		}
		debug("assigned: ", field.Name, tag)
	}
	return nil
}

// Walk struct field and assign from command-line arguments
func cascadeCli(v reflect.Value, cliOptions map[string]string) error {
	t := derefType(v.Type())
	v = derefValue(v)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if !value.CanSet() {
			debug("cannot set: ", field.Name)
			continue
		}

		ft := field.Type
		var isPtr bool
		if ft.Kind() == reflect.Ptr {
			isPtr = true
			ft = derefType(ft)
		}

		if ft.Kind() == reflect.Struct {
			cascadeCli(value, cliOptions)
			continue
		}
		tag, ok := field.Tag.Lookup(tagNameCli)
		if !ok || tag == "" || tag == "-" {
			continue
		}
		var cliValue string
		var found bool
		for _, name := range strings.Split(tag, ",") {
			if vv, ok := cliOptions[name]; ok {
				cliValue = vv
				found = true
				break
			}
		}
		if !found {
			continue
		}
		if err := assignValue(ft, value, cliValue, isPtr, true); err != nil {
			return errors.Wrap(err, "failed to assign values")
		}
		debug("assigned: ", field.Name, tag)
	}
	return nil
}

// Merge override config
func mergeConfig(v, merge reflect.Value, tagName string) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		debug("processing: ", field.Name, field.Type.Kind())
		target := merge.Field(i)
		if !target.IsValid() {
			debug("target invalid: ", field.Name)
			continue
		}
		tag, ok := field.Tag.Lookup(tagName)
		if !ok || tag == "" || tag == "-" {
			debug("tag not found: ", tagName, field.Name)
			continue
		}
		if field.Type.Kind() == reflect.Struct {
			debug("nested struct: ", field.Name)
			mergeConfig(v.Field(i), derefValue(target), tagName)
		} else {
			v.Field(i).Set(target)
		}
	}
	return nil
}

// Check struct field has already been assigned some value from other cascading
func isZeroValue(ft reflect.Type, value reflect.Value) bool {
	switch ft.Kind() {
	case reflect.String:
		return value.String() == ""
	case reflect.Bool:
		return value.Bool() == false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	}
	return false
}

// Assign value which corresponds to struct fiele type.
// Currently we only support some primitive values like (int, uint, float, string)
// because configurations are enough to use those values.
func assignValue(ft reflect.Type, value reflect.Value, envValue string, isPtr bool, cliAssign bool) error {
	switch ft.Kind() {
	case reflect.String:
		if isPtr {
			value.Set(reflect.ValueOf(&envValue))
		} else {
			value.SetString(envValue)
		}
	case reflect.Bool:
		var b bool
		if cliAssign {
			b = envValue == "" || envValue == "yes"
		} else {
			b = envValue == "true" || envValue == "yes"
		}
		if isPtr {
			value.Set(reflect.ValueOf(&b))
		} else {
			value.SetBool(b)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if cliAssign && envValue == "" {
			return nil
		}

		i, err := strconv.ParseInt(envValue, 10, 64)
		if err != nil {
			return errors.Wrap(err, "failed to convert from string to int")
		}
		if isPtr {
			value.Set(reflect.ValueOf(&i))
		} else {
			value.SetInt(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if cliAssign && envValue == "" {
			return nil
		}
		ui, err := strconv.ParseUint(envValue, 10, 64)
		if err != nil {
			return errors.Wrap(err, "failed to convert from string to uint")
		}
		if isPtr {
			value.Set(reflect.ValueOf(&ui))
		} else {
			value.SetUint(ui)
		}
	case reflect.Float32, reflect.Float64:
		if cliAssign && envValue == "" {
			return nil
		}
		f, err := strconv.ParseFloat(envValue, 64)
		if err != nil {
			return errors.Wrap(err, "failed to convert from string to float")
		}
		if isPtr {
			value.Set(reflect.ValueOf(&f))
		} else {
			value.SetFloat(f)
		}
	}
	return nil
}

// Parse command-line argument strings to map with short/long keys
func parseCliArgs(args []string) map[string]string {
	options := make(map[string]string)
	size := len(args)

	for i := 0; i < size; i++ {
		v := args[i]
		var name, value string

		if v[0] != '-' || len(v) <= 1 {
			continue
		}
		if v[1] == '-' {
			// Parse as long argument
			kv := strings.Split(v, "=")
			name = kv[0][2:]
			if len(kv) > 1 {
				value = kv[1]
			} else if i+1 < size {
				value = args[i+1]
				i++
			} else {
				value = ""
			}
		} else {
			// Parse as short argument
			name = string(v[1])
			if len(v) == 2 {
				if i+1 < size {
					value = args[i+1]
				} else {
					value = ""
				}
			} else {
				value = v[2:]
			}
		}
		options[name] = value
	}

	return options
}
