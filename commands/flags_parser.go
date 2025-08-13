package commands

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
)

// CommandFlagsParser used for parsing the arguments
type CommandFlagsParser struct {
	flag       *flag.FlagSet
	parser     FlagsParser
	validator  FlagsValidator
	parsedArgs []string
}

// NewCommandFlagsParser creates new command flags parser
func NewCommandFlagsParser(flag *flag.FlagSet, parser FlagsParser, validator FlagsValidator) CommandFlagsParser {
	return CommandFlagsParser{flag: flag, parser: parser, validator: validator, parsedArgs: make([]string, 0)}
}

// Parse parsing the args
func (p *CommandFlagsParser) Parse(args []string) error {
	if unknownFlags := collectUnknownFlags(p.flag, args); len(unknownFlags) > 0 {
		return fmt.Errorf("Unknown or wrong flags: %s", strings.Join(unknownFlags, ", "))
	}

	err := p.parser.ParseFlags(p.flag, args)
	if err != nil {
		return err
	}

	// assume that the parsing of arguments is successful - determine the arguments which are not flagged
	p.parsedArgs = determineParsedNotFlaggedArguments(p.flag, args)

	return p.validator.ValidateParsedFlags(p.flag)
}

func determineParsedNotFlaggedArguments(flag *flag.FlagSet, args []string) []string {
	result := make([]string, 0)
	for _, arg := range args {
		if argument := flag.Lookup(strings.Replace(arg, "-", "", 2)); argument == nil {
			result = append(result, arg)
		} else {
			break
		}
	}
	return result
}

// Args returns the first parsed command line arguments WITHOUT the options
func (p CommandFlagsParser) Args() []string {
	return p.parsedArgs
}

// FlagsParser interface used for parsing the command line arguments using the flag library
type FlagsParser interface {
	ParseFlags(flags *flag.FlagSet, args []string) error
}

// FlagsValidator interface used for validating the parsed flags
type FlagsValidator interface {
	ValidateParsedFlags(flags *flag.FlagSet) error
}

// DefaultCommandFlagsParser defines default implementation of the parser. It uses positional arguments and assumes that the command args will contain arguments
type DefaultCommandFlagsParser struct {
	positionalArgNames []string
}

// NewDefaultCommandFlagsParser initializes DefaultCommandFlagsParser
func NewDefaultCommandFlagsParser(positionalArgNames []string) DefaultCommandFlagsParser {
	return DefaultCommandFlagsParser{positionalArgNames: positionalArgNames}
}

// ParseFlags see DefaultCommandFlagsParser
func (p DefaultCommandFlagsParser) ParseFlags(flags *flag.FlagSet, args []string) error {
	// Check for missing positional arguments
	positionalArgsCount := len(p.positionalArgNames)
	if len(args) < positionalArgsCount {
		return fmt.Errorf("Missing positional argument %q", p.positionalArgNames[len(args)])
	}
	for i := 0; i < positionalArgsCount; i++ {
		if flags.Lookup(strings.Replace(args[i], "-", "", 1)) != nil {
			return fmt.Errorf("Missing positional argument %q", p.positionalArgNames[i])
		}
	}

	// Parse the arguments
	err := flags.Parse(args[positionalArgsCount:])
	if err != nil {
		return errors.New("Parsing of arguments has failed")
	}

	// Check for wrong arguments
	if flags.NArg() > 0 {
		return errors.New("Wrong arguments")
	}
	return nil
}

// DefaultCommandFlagsValidator default implementation of the FlagValidator
type DefaultCommandFlagsValidator struct {
	requiredFlags map[string]bool
}

// NewDefaultCommandFlagsValidator creates a default validator for flags
func NewDefaultCommandFlagsValidator(requiredFlags map[string]bool) DefaultCommandFlagsValidator {
	return DefaultCommandFlagsValidator{requiredFlags: requiredFlags}
}

// ValidateParsedFlags uses a required flags map in order to validate whether the arguments are valid
func (v DefaultCommandFlagsValidator) ValidateParsedFlags(flags *flag.FlagSet) error {
	var missingRequiredOptions []string
	// Check for missing required flags
	flags.VisitAll(func(f *flag.Flag) {
		log.Traceln(f.Name, f.Value)
		if v.requiredFlags[f.Name] && f.Value.String() == "" {
			missingRequiredOptions = append(missingRequiredOptions, f.Name)
		}
	})
	if len(missingRequiredOptions) != 0 {
		return fmt.Errorf("Missing required options '%v'", missingRequiredOptions)
	}

	return nil
}

func collectUnknownFlags(flags *flag.FlagSet, args []string) []string {
	var unknownFlags []string

	for i := 0; i < len(args); i++ {
		currentArgument := args[i]

		if !strings.HasPrefix(currentArgument, "-") {
			continue
		}

		currentFlag := currentArgument
		flagName := strings.TrimLeft(currentFlag, "-")

		if flagName == "" {
			continue
		}

		isFlagKnown := flags.Lookup(flagName)
		if isFlagKnown != nil {
			nextIndex := i + 1
			if nextIndex < len(args) {
				isBoolean := isBoolFlag(isFlagKnown)
				if !isBoolean {
					nextArgument := args[nextIndex]
					nextHasPrefixDash := strings.HasPrefix(nextArgument, "-")
					if !nextHasPrefixDash {
						i = nextIndex
					}
				}
			}
			continue
		}

		unknownFlags = append(unknownFlags, currentFlag)

		nextIndex := i + 1
		if nextIndex < len(args) {
			nextArgument := args[nextIndex]
			nextHasPrefixDash := strings.HasPrefix(nextArgument, "-")
			if !nextHasPrefixDash {
				i = nextIndex
			}
		}
	}

	return unknownFlags
}

func isBoolFlag(flag *flag.Flag) bool {
	type boolFlagInterface interface{ IsBoolFlag() bool }

	boolFlag, isInterfaceImplemented := flag.Value.(boolFlagInterface)
	if !isInterfaceImplemented {
		return false
	}

	return boolFlag.IsBoolFlag()
}
