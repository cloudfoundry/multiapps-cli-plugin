package commands

import (
	"fmt"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	plugin_models "github.com/cloudfoundry/cli/plugin/models"
)

type Context struct {
	Username string
	Org      string
	Space    string
}

func CreateContext(cliConnection plugin.CliConnection) (Context, error) {
	username, err := GetUsername(cliConnection)
	if err != nil {
		return Context{}, err
	}
	org, err := GetOrg(cliConnection)
	if err != nil {
		return Context{}, err
	}
	space, err := GetSpace(cliConnection)
	if err != nil {
		return Context{}, err
	}
	return Context{username, org.Name, space.Name}, nil
}

func GetOrg(cliConnection plugin.CliConnection) (plugin_models.Organization, error) {
	org, err := cliConnection.GetCurrentOrg()
	if err != nil {
		return plugin_models.Organization{}, fmt.Errorf("Could not get current org: %s", err)
	}
	if org.Name == "" {
		return plugin_models.Organization{}, fmt.Errorf("No org and space targeted, use '%s' to target an org and a space", terminal.CommandColor("cf target -o ORG -s SPACE"))
	}
	return org, nil
}

func GetSpace(cliConnection plugin.CliConnection) (plugin_models.Space, error) {
	space, err := cliConnection.GetCurrentSpace()
	if err != nil {
		return plugin_models.Space{}, fmt.Errorf("Could not get current space: %s", err)
	}

	if space.Name == "" || space.Guid == "" {
		return plugin_models.Space{}, fmt.Errorf("No space targeted, use '%s' to target a space", terminal.CommandColor("cf target -s"))
	}
	return space, nil
}

func GetUsername(cliConnection plugin.CliConnection) (string, error) {
	username, err := cliConnection.Username()
	if err != nil {
		return "", fmt.Errorf("Could not get username: %s", err)
	}
	if username == "" {
		return "", fmt.Errorf("Not logged in. Use '%s' to log in.", terminal.CommandColor("cf login"))
	}
	return username, nil
}
