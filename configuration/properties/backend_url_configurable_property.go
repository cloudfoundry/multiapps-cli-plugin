package properties

var BackendURL = ConfigurableProperty{
	Name:                  "MULTIAPPS_CONTROLLER_URL",
	Parser:                noOpParser{},
	ParsingSuccessMessage: "Attention: You've specified a custom backend URL (%s) via the environment variable \"%s\". The application listening on that URL may be outdated, contain bugs or unreleased features or may even be modified by a potentially untrused person. Use at your own risk.\n",
	ParsingFailureMessage: "No validation implemented for custom backend URLs. If you're seeing this message then something has gone horribly wrong.\n",
	DefaultValue:          "",
}
