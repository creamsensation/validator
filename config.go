package validator

type Config struct {
	Messages Messages
}

type Messages struct {
	Required  string `json:"required" toml:"required" yaml:"required"`
	MinText   string `json:"minText" toml:"minText" yaml:"minText"`
	MaxText   string `json:"maxText" toml:"maxText" yaml:"maxText"`
	MinNumber string `json:"minNumber" toml:"minNumber" yaml:"minNumber"`
	MaxNumber string `json:"maxNumber" toml:"maxNumber" yaml:"maxNumber"`
}

const (
	defaultRequiredMessage  = "field is required"
	defaultMinTextMessage   = "field length is smaller than should be"
	defaultMaxTextMessage   = "field length is higher than should be"
	defaultMinNumberMessage = "field value is smaller than should be"
	defaultMaxNumberMessage = "field value is higher than should be"
)
