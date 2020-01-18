package configuration

type noOpParser struct {
}

func (p noOpParser) Parse(value string) (interface{}, error) {
	return value, nil
}
