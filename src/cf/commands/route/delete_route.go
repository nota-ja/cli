package route

import (
	"cf/api"
	"cf/command_metadata"
	"cf/configuration"
	"cf/errors"
	"cf/flag_helpers"
	"cf/requirements"
	"cf/terminal"
	"github.com/codegangsta/cli"
)

type DeleteRoute struct {
	ui        terminal.UI
	config    configuration.Reader
	routeRepo api.RouteRepository
}

func NewDeleteRoute(ui terminal.UI, config configuration.Reader, routeRepo api.RouteRepository) (cmd *DeleteRoute) {
	cmd = new(DeleteRoute)
	cmd.ui = ui
	cmd.config = config
	cmd.routeRepo = routeRepo
	return
}

func (command *DeleteRoute) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "delete-route",
		Description: "Delete a route",
		Usage:       "CF_NAME delete-route DOMAIN [-n HOSTNAME] [-f]",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "f", Usage: "Force deletion without confirmation"},
			flag_helpers.NewStringFlag("n", "Hostname"),
		},
	}
}

func (cmd *DeleteRoute) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {

	if len(c.Args()) != 1 {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c, "delete-route")
		return
	}

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}
	return
}

func (cmd *DeleteRoute) Run(c *cli.Context) {
	host := c.String("n")
	domainName := c.Args()[0]

	url := domainName
	if host != "" {
		url = host + "." + domainName
	}
	if !c.Bool("f") {
		if !cmd.ui.ConfirmDelete("route", url) {
			return
		}
	}

	cmd.ui.Say("Deleting route %s...", terminal.EntityNameColor(url))

	route, apiErr := cmd.routeRepo.FindByHostAndDomain(host, domainName)

	switch apiErr.(type) {
	case nil:
	case *errors.ModelNotFoundError:
		cmd.ui.Warn("Unable to delete, route '%s' does not exist.", url)
		return
	default:
		cmd.ui.Failed(apiErr.Error())
		return
	}

	apiErr = cmd.routeRepo.Delete(route.Guid)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
}
