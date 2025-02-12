package space

import (
	"cf/api"
	"cf/command_metadata"
	"cf/configuration"
	"cf/models"
	"cf/requirements"
	"cf/terminal"
	"github.com/codegangsta/cli"
)

type ListSpaces struct {
	ui        terminal.UI
	config    configuration.Reader
	spaceRepo api.SpaceRepository
}

func NewListSpaces(ui terminal.UI, config configuration.Reader, spaceRepo api.SpaceRepository) (cmd ListSpaces) {
	cmd.ui = ui
	cmd.config = config
	cmd.spaceRepo = spaceRepo
	return
}

func (command ListSpaces) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "spaces",
		Description: "List all spaces in an org",
		Usage:       "CF_NAME spaces",
	}
}

func (cmd ListSpaces) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewTargetedOrgRequirement(),
	}
	return
}

func (cmd ListSpaces) Run(c *cli.Context) {
	cmd.ui.Say("Getting spaces in org %s as %s...\n",
		terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
		terminal.EntityNameColor(cmd.config.Username()))

	foundSpaces := false
	table := cmd.ui.Table([]string{"name"})
	apiErr := cmd.spaceRepo.ListSpaces(func(space models.Space) bool {
		table.Print([][]string{{space.Name}})
		foundSpaces = true
		return true
	})

	if apiErr != nil {
		cmd.ui.Failed("Failed fetching spaces.\n%s", apiErr.Error())
		return
	}

	if !foundSpaces {
		cmd.ui.Say("No spaces found")
	}
}
