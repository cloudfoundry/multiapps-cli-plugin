package commands

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

// CommandFlagsParser used for parsing the arguments
type CommandFlagsParser struct {
	flags     *flag.FlagSet
	parser    FlagsParser
	validator FlagsValidator
}

// FlagsParser interface used for parsing the command line arguments using the flag library
type FlagsParser interface {
	ParseFlags(flags *flag.FlagSet, args []string) error
}

// FlagsValidator interface used for validating the parsed flags
type FlagsValidator interface {
	ValidateFlags(flags *flag.FlagSet, args []string) error
	IsBeforeParsing() bool
}

// CommandOption defines an option for a command
type CommandOption struct {
	Value        interface{}
	DefaultValue interface{}
	Usage        string
	IsShortOpt   bool
}

func NewCommandFlagsParserWithValidator(flags *flag.FlagSet, parser FlagsParser, validator FlagsValidator) CommandFlagsParser {
	return CommandFlagsParser{flags, parser, validator}
}

func NewCommandFlagsParser(flags *flag.FlagSet, parser FlagsParser) CommandFlagsParser {
	return CommandFlagsParser{flags: flags, parser: parser}
}

// Parse parses the args
func (p *CommandFlagsParser) Parse(args []string) error {
	validated := false
	if p.validator != nil && p.validator.IsBeforeParsing() {
		if err := p.validator.ValidateFlags(p.flags, args); err != nil {
			return err
		}
		validated = true
	}
	if err := p.parser.ParseFlags(p.flags, args); err != nil {
		return err
	}
	if !validated && p.validator != nil {
		if err := p.validator.ValidateFlags(p.flags, args); err != nil {
			return err
		}
	}
	return nil
}

// DefaultCommandFlagsParser defines default implementation of the parser
// It assumes that the command args will contain arguments
type DefaultCommandFlagsParser struct {
	offset int
}

// NewDefaultCommandFlagsParser initializes DefaultCommandFlagsParser
func NewDefaultCommandFlagsParser(offset int) DefaultCommandFlagsParser {
	return DefaultCommandFlagsParser{offset}
}

// ParseFlags see DefaultCommandFlagsParser
func (p DefaultCommandFlagsParser) ParseFlags(flags *flag.FlagSet, args []string) error {
	// Parse the arguments
	err := flags.Parse(args[p.offset:])
	if err != nil {
		return errors.New("Unknown or wrong flag")
	}
	// Check for wrong arguments
	if flags.NArg() > 0 {
		return errors.New("Wrong arguments")
	}
	return nil
}

type ProcessActionExecutorCommandArgumentsParser struct {
	offset int
}

func NewProcessActionExecutorCommandArgumentsParser(offset int) ProcessActionExecutorCommandArgumentsParser {
	return ProcessActionExecutorCommandArgumentsParser{offset}
}

func (p ProcessActionExecutorCommandArgumentsParser) ParseFlags(flags *flag.FlagSet, args []string) error {
	executeActionOptCount := make(map[string]int)
	for _, arg := range args {
		optionFlag := flags.Lookup(strings.Replace(arg, "-", "", 1))
		if optionFlag != nil && (operationIDOpt == optionFlag.Name || actionOpt == optionFlag.Name) {
			executeActionOptCount[optionFlag.Name]++
		}
	}

	if len(executeActionOptCount) > 2 || p.areOptionsSpecifiedMoreThanOnce(executeActionOptCount) {
		return fmt.Errorf("Options %s and %s should be specified only once", operationIDOpt, actionOpt)
	}
	if len(executeActionOptCount) == 1 {
		return errors.New("All the a i options should be specified together")
	}

	offset := p.offset
	if len(executeActionOptCount) == 2 {
		offset = 0
	}
	return NewDefaultCommandFlagsParser(offset).ParseFlags(flags, args)
}

func (p *ProcessActionExecutorCommandArgumentsParser) areOptionsSpecifiedMoreThanOnce(executeActionOptCount map[string]int) bool {
	for _, num := range executeActionOptCount {
		if num > 1 {
			return true
		}
	}
	return false
}

type PositionalArgumentsFlagsValidator struct {
	positionalArgs []string
}

func NewPositionalArgumentsFlagsValidator(positionalArgs []string) *PositionalArgumentsFlagsValidator {
	return &PositionalArgumentsFlagsValidator{positionalArgs}
}

func (v *PositionalArgumentsFlagsValidator) ValidateFlags(flags *flag.FlagSet, args []string) error {
	// Check for missing positional arguments
	positionalArgsCount := len(v.positionalArgs)
	if len(args) < positionalArgsCount {
		return fmt.Errorf("Missing positional argument '%s'", v.positionalArgs[len(args)])
	}
	for i := 0; i < positionalArgsCount; i++ {
		if flags.Lookup(strings.Replace(args[i], "-", "", 1)) != nil {
			return fmt.Errorf("Missing positional argument '%s'", v.positionalArgs[i])
		}
	}
	return nil
}

func (*PositionalArgumentsFlagsValidator) IsBeforeParsing() bool {
	return true
}
