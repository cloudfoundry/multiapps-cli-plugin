package configuration

type configurableProperty struct {
	Name                  string
	DeprecatedNames       []string
	Parser                configurablePropertyParser
	ParsingSuccessMessage string
	ParsingFailureMessage string
	DefaultValue          interface{}
}

type configurablePropertyParser interface {
	Parse(value string) (interface{}, error)
}
