package validator

import (
	"reflect"
)

type Validator interface {
	Clear() Validator
	Add(name string) Field
	Json(data any) Validator
	Ok() bool
	Errors() map[string][]string
}

type validator struct {
	config Config
	errors map[string][]string
	fields map[string]*field
	ok     bool
}

func New(config ...Config) Validator {
	v := &validator{
		errors: make(map[string][]string),
		fields: make(map[string]*field),
	}
	if len(config) > 0 {
		v.config = config[0]
	}
	v.config = v.createDefaultConfig(v.config)
	return v
}

func (v *validator) Add(name string) Field {
	f := &field{}
	v.fields[name] = f
	return f
}

func (v *validator) Errors() map[string][]string {
	return v.errors
}

func (v *validator) Clear() Validator {
	clear(v.errors)
	v.ok = false
	return v
}

func (v *validator) Json(data any) Validator {
	dt := reflect.TypeOf(data)
	dv := reflect.ValueOf(data)
	switch dt.Kind() {
	case reflect.Struct:
		v.validateStruct(dt, dv)
	case reflect.Map:
		v.validateMap(dv)
	}
	v.ok = len(v.errors) == 0
	return v
}

func (v *validator) Ok() bool {
	return v.ok
}

func (v *validator) validateStruct(mt reflect.Type, mv reflect.Value) {
	for i := 0; i < mv.NumField(); i++ {
		jsonKey := mt.Field(i).Tag.Get("json")
		if len(jsonKey) == 0 {
			continue
		}
		errs := v.validateField(jsonKey, mv.Field(i).Interface())
		if len(errs) > 0 {
			v.errors[jsonKey] = errs
		}
	}
}

func (v *validator) validateMap(mv reflect.Value) {
	keys := mv.MapKeys()
	for name := range v.fields {
		exist := false
		for _, k := range keys {
			if k.String() == name {
				exist = true
			}
		}
		if exist {
			continue
		}
		v.errors[name] = []string{v.config.Messages.Required}
	}
	for _, k := range keys {
		errs := v.validateField(k.String(), mv.MapIndex(k).Interface())
		if len(errs) > 0 {
			v.errors[k.String()] = errs
		}
	}
}

func (v *validator) validateField(key string, value any) []string {
	errs := make([]string, 0)
	f, ok := v.fields[key]
	if !ok {
		return errs
	}
	switch val := value.(type) {
	case string:
		if f.required && len(val) == 0 {
			errs = append(errs, v.config.Messages.Required)
		}
		if f.min > 0 && len(val) < f.min {
			errs = append(errs, v.config.Messages.MinText)
		}
		if f.max > 0 && len(val) > f.max {
			errs = append(errs, v.config.Messages.MaxText)
		}
	case int:
		if f.required && val == 0 {
			errs = append(errs, v.config.Messages.Required)
		}
		if f.min > 0 && val < f.min {
			errs = append(errs, v.config.Messages.MinNumber)
		}
		if f.max > 0 && val > f.max {
			errs = append(errs, v.config.Messages.MaxNumber)
		}
	case float32:
		if f.required && val == 0 {
			errs = append(errs, v.config.Messages.Required)
		}
		if f.min > 0 && val < float32(f.min) {
			errs = append(errs, v.config.Messages.MinNumber)
		}
		if f.max > 0 && val > float32(f.max) {
			errs = append(errs, v.config.Messages.MaxNumber)
		}
	case float64:
		if f.required && val == 0 {
			errs = append(errs, v.config.Messages.Required)
		}
		if f.min > 0 && val < float64(f.min) {
			errs = append(errs, v.config.Messages.MinNumber)
		}
		if f.max > 0 && val > float64(f.max) {
			errs = append(errs, v.config.Messages.MaxNumber)
		}
	case bool:
		if f.required && !val {
			errs = append(errs, v.config.Messages.Required)
		}
	}
	return errs
}

func (v *validator) createDefaultConfig(config Config) Config {
	if len(config.Messages.Required) == 0 {
		config.Messages.Required = defaultRequiredMessage
	}
	if len(config.Messages.MinText) == 0 {
		config.Messages.MinText = defaultMinTextMessage
	}
	if len(config.Messages.MaxText) == 0 {
		config.Messages.MaxText = defaultMaxTextMessage
	}
	if len(config.Messages.MinNumber) == 0 {
		config.Messages.MinNumber = defaultMinNumberMessage
	}
	if len(config.Messages.MaxNumber) == 0 {
		config.Messages.MaxNumber = defaultMaxNumberMessage
	}
	return config
}
