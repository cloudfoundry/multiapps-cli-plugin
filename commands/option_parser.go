package commands

import "flag"

type OptionParser interface {
	parseOption(name string, option CommandOption, flags *flag.FlagSet) bool
	additionalParse(name string, option CommandOption, flags *flag.FlagSet)
}

type AbstractOptionParser struct {
}

func (AbstractOptionParser) parseOption(name string, option CommandOption, flags *flag.FlagSet) bool {
	parsed := true
	switch val := option.Value.(type) {
	case *string:
		flags.StringVar(val, name, option.DefaultValue.(string), "")
	case *uint:
		flags.UintVar(val, name, uint(option.DefaultValue.(int)), "")
	case *bool:
		flags.BoolVar(val, name, option.DefaultValue.(bool), "")
	default:
		parsed = false
	}
	return parsed
}

type DefaultOptionParser struct {
	AbstractOptionParser
}

func (DefaultOptionParser) additionalParse(name string, option CommandOption, flags *flag.FlagSet) {
	//no-op
}
