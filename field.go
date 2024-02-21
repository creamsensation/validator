package validator

type Field interface {
	Required() Field
	Min(min int) Field
	Max(max int) Field
}

type field struct {
	min      int
	max      int
	required bool
}

func (f *field) Required() Field {
	f.required = true
	return f
}

func (f *field) Min(min int) Field {
	f.min = min
	return f
}

func (f *field) Max(max int) Field {
	f.max = max
	return f
}
