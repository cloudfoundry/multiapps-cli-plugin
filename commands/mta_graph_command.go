package commands

import (
	"flag"
	"fmt"
	"os"

	"code.cloudfoundry.org/cli/v8/plugin"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/util"
	"gopkg.in/yaml.v3"
)

type MtaGraphCommand struct {
	name          string
	cliConnection plugin.CliConnection
}

func NewMtaGraphCommand() *MtaGraphCommand {
	return &MtaGraphCommand{}
}

func (c *MtaGraphCommand) GetPluginCommand() plugin.Command {
	return plugin.Command{
		Name:     "mta-graph",
		HelpText: "Render the dependency graph of an MTA descriptor as DOT",
		UsageDetails: plugin.Usage{
			Usage: "cf mta-graph <path-to-mtad.yaml>",
		},
	}
}

func (c *MtaGraphCommand) Initialize(name string, cliConnection plugin.CliConnection) {
	c.name = name
	c.cliConnection = cliConnection
}

func (c *MtaGraphCommand) Execute(args []string) ExecutionStatus {
	if len(args) == 0 {
		ui.Failed("Usage: cf mta-graph <path-to-mtad.yaml>")
		return Failure
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		ui.Failed("Could not read file %s: %s", args[0], err)
		return Failure
	}

	descriptor, err := parseDescriptor(data)
	if err != nil {
		ui.Failed("Could not parse descriptor: %s", err)
		return Failure
	}

	fmt.Print(buildDot(descriptor))
	return Success
}

// executeInternal satisfies the Command interface (never called for this command).
func (c *MtaGraphCommand) executeInternal(_ []string, _ string, _ *flag.FlagSet, _ util.CloudFoundryTarget) ExecutionStatus {
	return Failure
}

func (c *MtaGraphCommand) defineCommandOptions(_ *flag.FlagSet) {}

// ---- descriptor model ----

type mtadDescriptor struct {
	Modules   []mtadModule   `yaml:"modules"`
	Resources []mtadResource `yaml:"resources"`
}

type mtadModule struct {
	Name     string          `yaml:"name"`
	Provides []mtadProvided  `yaml:"provides"`
	Requires []mtadRequired  `yaml:"requires"`
}

type mtadProvided struct {
	Name string `yaml:"name"`
}

type mtadRequired struct {
	Name string `yaml:"name"`
}

type mtadResource struct {
	Name string `yaml:"name"`
}

func parseDescriptor(data []byte) (*mtadDescriptor, error) {
	var d mtadDescriptor
	if err := yaml.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

// ---- DOT builder ----

func buildDot(d *mtadDescriptor) string {
	providedDepToModule := make(map[string]string)
	for _, m := range d.Modules {
		for _, p := range m.Provides {
			providedDepToModule[p.Name] = m.Name
		}
	}

	resourceNames := make(map[string]bool)
	for _, r := range d.Resources {
		resourceNames[r.Name] = true
	}

	dot := "digraph mta {\n"
	dot += "  rankdir=LR;\n"
	dot += "  node [fontname=\"Helvetica\" fontsize=11];\n"
	dot += "  edge [fontname=\"Helvetica\" fontsize=10];\n"
	for _, m := range d.Modules {
		dot += fmt.Sprintf("  %q [shape=box style=\"rounded,filled\" fillcolor=\"#DCE9F7\" color=\"#4A7AB5\"];\n", m.Name)
	}
	dot += "  { rank=same;"
	for _, r := range d.Resources {
		dot += fmt.Sprintf(" %q", r.Name)
	}
	dot += " }\n"
	for _, r := range d.Resources {
		dot += fmt.Sprintf("  %q [shape=cylinder style=filled fillcolor=\"#EDEDED\" color=\"#8A8A8A\"];\n", r.Name)
	}
	for _, m := range d.Modules {
		for _, req := range m.Requires {
			if provider, ok := providedDepToModule[req.Name]; ok {
				dot += fmt.Sprintf("  %q -> %q [label=\"provides %s\" color=\"#4A7AB5\" fontcolor=\"#4A7AB5\"];\n", provider, m.Name, req.Name)
			} else if resourceNames[req.Name] {
				dot += fmt.Sprintf("  %q -> %q [label=\"requires\" style=dashed color=\"#8A8A8A\" fontcolor=\"#8A8A8A\"];\n", m.Name, req.Name)
			}
		}
	}
	dot += "}\n"
	return dot
}
