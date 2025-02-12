package quota

import (
	"cf/api"
	"cf/command_metadata"
	"cf/configuration"
	"cf/formatters"
	"cf/requirements"
	"cf/terminal"
	"fmt"
	"github.com/codegangsta/cli"
)

type showQuota struct {
	ui        terminal.UI
	config    configuration.Reader
	quotaRepo api.QuotaRepository
}

func NewShowQuota(ui terminal.UI, config configuration.Reader, quotaRepo api.QuotaRepository) *showQuota {
	return &showQuota{
		ui:        ui,
		config:    config,
		quotaRepo: quotaRepo,
	}
}

func (command *showQuota) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "quota",
		Usage:       "CF_NAME quota QUOTA",
		Description: "Show quota info",
	}
}

func (cmd *showQuota) GetRequirements(requirementsFactory requirements.Factory, context *cli.Context) ([]requirements.Requirement, error) {
	if len(context.Args()) != 1 {
		cmd.ui.FailWithUsage(context, "quotas")
	}

	return []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}, nil
}

func (cmd *showQuota) Run(context *cli.Context) {
	quotaName := context.Args()[0]
	cmd.ui.Say("Getting quota %s info as %s...", quotaName, cmd.config.Username())

	quota, err := cmd.quotaRepo.FindByName(quotaName)
	if err != nil {
		cmd.ui.Failed(err.Error())
	}

	cmd.ui.Ok()

	table := terminal.NewTable(cmd.ui, []string{"", ""})
	table.Print([][]string{
		[]string{"Memory", formatters.ByteSize(quota.MemoryLimit * formatters.MEGABYTE)},
		[]string{"Routes", fmt.Sprintf("%d", quota.RoutesLimit)},
		[]string{"Services", fmt.Sprintf("%d", quota.ServicesLimit)},
		[]string{"Paid service plans", formatters.Allowed(quota.NonBasicServicesAllowed)},
	})
}
