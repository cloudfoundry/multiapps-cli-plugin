package util

import (
	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
	plugin_models "code.cloudfoundry.org/cli/plugin/models"
	"fmt"
)

type CloudFoundryContext struct {
	cliConnection plugin.CliConnection
}

func NewCloudFoundryContext(cliConnection plugin.CliConnection) CloudFoundryContext {
	return CloudFoundryContext{cliConnection: cliConnection}
}

// GetOrg gets the current org from the CLI connection
func (c *CloudFoundryContext) GetOrg() (plugin_models.Organization, error) {
	org, err := c.cliConnection.GetCurrentOrg()
	if err != nil {
		return plugin_models.Organization{}, fmt.Errorf("Could not get current org: %s", err)
	}
	if org.Name == "" {
		return plugin_models.Organization{}, fmt.Errorf("No org and space targeted, use '%s' to target an org and a space", terminal.CommandColor("cf target -o ORG -s SPACE"))
	}
	return org, nil
}

// GetSpace gets the current space from the CLI connection
func (c *CloudFoundryContext) GetSpace() (plugin_models.Space, error) {
	space, err := c.cliConnection.GetCurrentSpace()
	if err != nil {
		return plugin_models.Space{}, fmt.Errorf("Could not get current space: %s", err)
	}

	if space.Name == "" || space.Guid == "" {
		return plugin_models.Space{}, fmt.Errorf("No space targeted, use '%s' to target a space", terminal.CommandColor("cf target -s SPACE"))
	}
	return space, nil
}

// GetUsername gets the username from the CLI connection
func (c *CloudFoundryContext) GetUsername() (string, error) {
	username, err := c.cliConnection.Username()
	if err != nil {
		return "", fmt.Errorf("Could not get username: %s", err)
	}
	if username == "" {
		return "", fmt.Errorf("Not logged in. Use '%s' to log in.", terminal.CommandColor("cf login"))
	}
	return username, nil
}
