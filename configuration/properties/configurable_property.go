package properties

type ConfigurableProperty struct {
	Name                  string
	Parser                Parser
	ParsingSuccessMessage string
	ParsingFailureMessage string
	DefaultValue          interface{}
}

type Parser interface {
	Parse(value string) (interface{}, error)
}

type noOpParser struct {
}

func (p noOpParser) Parse(value string) (interface{}, error) {
	return value, nil
}
