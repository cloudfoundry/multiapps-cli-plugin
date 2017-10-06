package util

import (
	"bytes"
	"strings"
)

// CfCommandStringBuilder ...
type CfCommandStringBuilder struct {
	name string
	args bytes.Buffer
	opts bytes.Buffer
}

const longOptionPrefix = "--"
const optionPrefix = "-"

// NewCfCommandStringBuilder creates a new CfCommandStringBuilder
func NewCfCommandStringBuilder() *CfCommandStringBuilder {
	return &CfCommandStringBuilder{}
}

// SetName sets the name of the command that will be generated from the Build method
func (builder *CfCommandStringBuilder) SetName(name string) *CfCommandStringBuilder {
	builder.name = name
	return builder
}

// AddBooleanOption adds a short boolean option to the command that will be generated from the Build method
func (builder *CfCommandStringBuilder) AddBooleanOption(option string) *CfCommandStringBuilder {
	return builder.addBooleanOption(option, optionPrefix)
}

// AddOption adds a short option to the command that will be generated from the Build method
func (builder *CfCommandStringBuilder) AddOption(option, value string) *CfCommandStringBuilder {
	return builder.addOption(option, value, optionPrefix)
}

func (builder *CfCommandStringBuilder) addBooleanOption(option, prefix string) *CfCommandStringBuilder {
	builder.opts.WriteString(prefix)
	builder.opts.WriteString(option)
	builder.opts.WriteRune(' ')
	return builder
}

func (builder *CfCommandStringBuilder) addOption(option, value, prefix string) *CfCommandStringBuilder {
	builder.opts.WriteString(prefix)
	builder.opts.WriteString(option)
	builder.opts.WriteRune(' ')
	builder.opts.WriteString(value)
	builder.opts.WriteRune(' ')
	return builder
}

// AddLongBooleanOption adds a long boolean option to the command that will be generated from the Build method
func (builder *CfCommandStringBuilder) AddLongBooleanOption(option string) *CfCommandStringBuilder {
	return builder.addBooleanOption(option, longOptionPrefix)
}

// AddLongOption adds a long option to the command that will be generated from the Build method
func (builder *CfCommandStringBuilder) AddLongOption(option, value string) *CfCommandStringBuilder {
	return builder.addOption(option, value, longOptionPrefix)
}

// AddArgument adds an argument to the command that will be generated from the Build method
func (builder *CfCommandStringBuilder) AddArgument(argument string) *CfCommandStringBuilder {
	builder.args.WriteString(argument)
	builder.args.WriteRune(' ')
	return builder
}

// Build generates a command string, in which the arguments always preceed the options
func (builder *CfCommandStringBuilder) Build() string {
	var commandBuilder bytes.Buffer

	commandBuilder.WriteString("cf ")
	commandBuilder.WriteString(builder.name)
	commandBuilder.WriteRune(' ')
	commandBuilder.WriteString(builder.args.String())
	commandBuilder.WriteString(builder.opts.String())

	return strings.Trim(commandBuilder.String(), " ")
}
