package util

import plugin_models "code.cloudfoundry.org/cli/plugin/models"

type CloudFoundryTarget struct {
	Org      plugin_models.Organization
	Space    plugin_models.Space
	Username string
}

func NewCFTarget(org plugin_models.Organization, space plugin_models.Space, username string) CloudFoundryTarget {
	return CloudFoundryTarget{
		Org:      org,
		Space:    space,
		Username: username,
	}
}
