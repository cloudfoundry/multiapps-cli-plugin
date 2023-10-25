package properties

const DefaultDisableProgressBar = false

var DisableProgressBar = ConfigurableProperty{
	Name:                  "MULTIAPPS_DISABLE_UPLOAD_PROGRESS_BAR",
	Parser:                booleanParser{},
	ParsingSuccessMessage: "Attention: You've specified %v for the environment variable %s.\n",
	ParsingFailureMessage: "Invalid boolean value (%s) for environment variable %s. Using default value %v.\n",
	DefaultValue:          DefaultDisableProgressBar,
}
