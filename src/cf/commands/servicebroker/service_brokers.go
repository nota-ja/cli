package servicebroker

import (
	"cf/api"
	"cf/command_metadata"
	"cf/configuration"
	"cf/models"
	"cf/requirements"
	"cf/terminal"
	"github.com/codegangsta/cli"
)

type ListServiceBrokers struct {
	ui     terminal.UI
	config configuration.Reader
	repo   api.ServiceBrokerRepository
}

func NewListServiceBrokers(ui terminal.UI, config configuration.Reader, repo api.ServiceBrokerRepository) (cmd ListServiceBrokers) {
	cmd.ui = ui
	cmd.config = config
	cmd.repo = repo
	return
}

func (command ListServiceBrokers) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "service-brokers",
		Description: "List service brokers",
		Usage:       "CF_NAME service-brokers",
	}
}

func (cmd ListServiceBrokers) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	reqs = append(reqs, requirementsFactory.NewLoginRequirement())
	return
}

func (cmd ListServiceBrokers) Run(c *cli.Context) {
	cmd.ui.Say("Getting service brokers as %s...\n", terminal.EntityNameColor(cmd.config.Username()))

	table := cmd.ui.Table([]string{"name", "url"})
	foundBrokers := false
	apiErr := cmd.repo.ListServiceBrokers(func(serviceBroker models.ServiceBroker) bool {
		table.Print([][]string{{serviceBroker.Name, serviceBroker.Url}})
		foundBrokers = true
		return true
	})

	if apiErr != nil {
		cmd.ui.Failed("Failed fetching service brokers.\n%s", apiErr)
		return
	}

	if !foundBrokers {
		cmd.ui.Say("No service brokers found")
	}
}
