package application_test

import (
	. "cf/commands/application"
	"cf/configuration"
	"cf/errors"
	"cf/formatters"
	"cf/models"
	testapi "testhelpers/api"
	testassert "testhelpers/assert"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
	testtime "testhelpers/time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("app Command", func() {
	var (
		ui                  *testterm.FakeUI
		configRepo          configuration.ReadWriter
		appSummaryRepo      *testapi.FakeAppSummaryRepo
		appInstancesRepo    *testapi.FakeAppInstancesRepo
		requirementsFactory *testreq.FakeReqFactory
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		appSummaryRepo = &testapi.FakeAppSummaryRepo{}
		appInstancesRepo = &testapi.FakeAppInstancesRepo{}
		configRepo = testconfig.NewRepositoryWithDefaults()
		requirementsFactory = &testreq.FakeReqFactory{
			LoginSuccess:         true,
			TargetedSpaceSuccess: true,
		}
	})

	runCommand := func(args ...string) {
		cmd := NewShowApp(ui, configRepo, appSummaryRepo, appInstancesRepo)
		testcmd.RunCommand(cmd, testcmd.NewContext("apps", args), requirementsFactory)
	}

	Describe("requirements", func() {
		It("fails if not logged in", func() {
			requirementsFactory.LoginSuccess = false
			runCommand("cf-plays-dwarf-fortress")
			Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
		})

		It("fails if a space is not targeted", func() {
			requirementsFactory.TargetedSpaceSuccess = false
			runCommand("cf-plays-dwarf-fortress")
			Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
		})

		It("fails with usage when no arguments are given", func() {
			runCommand()
			Expect(ui.FailedWithUsage).To(BeTrue())
			Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
		})

	})

	Describe("displaying a summary of an app", func() {
		BeforeEach(func() {
			app := makeAppWithRoute("my-app")
			appInstance := models.AppInstanceFields{
				State:     models.InstanceRunning,
				Since:     testtime.MustParse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 15:04:05 -0700 MST 2012"),
				CpuUsage:  1.0,
				DiskQuota: 1 * formatters.GIGABYTE,
				DiskUsage: 32 * formatters.MEGABYTE,
				MemQuota:  64 * formatters.MEGABYTE,
				MemUsage:  13 * formatters.BYTE,
			}

			appInstance2 := models.AppInstanceFields{
				State: models.InstanceDown,
				Since: testtime.MustParse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Apr 1 15:04:05 -0700 MST 2012"),
			}

			instances := []models.AppInstanceFields{appInstance, appInstance2}

			appSummaryRepo.GetSummarySummary = app
			appInstancesRepo.GetInstancesResponses = [][]models.AppInstanceFields{instances}
			requirementsFactory.Application = app
		})

		It("displays a summary of the app", func() {
			runCommand("my-app")

			Expect(appSummaryRepo.GetSummaryAppGuid).To(Equal("app-guid"))

			testassert.SliceContains(ui.Outputs, testassert.Lines{
				{"Showing health and status", "my-app"},
				{"state", "started"},
				{"instances", "2/2"},
				{"usage", "256M x 2 instances"},
				{"urls", "my-app.example.com", "foo.example.com"},
				{"#0", "running", "2012-01-02 03:04:05 PM", "100.0%", "13 of 64M", "32M of 1G"},
				{"#1", "down", "2012-04-01 03:04:05 PM", "0%", "0 of 0", "0 of 0"},
			})
		})
	})

	Describe("when the app is not running", func() {
		BeforeEach(func() {
			application := models.Application{}
			application.Name = "my-app"
			application.Guid = "my-app-guid"
			application.State = "stopped"
			application.InstanceCount = 2
			application.RunningInstances = 0
			application.Memory = 256

			appSummaryRepo.GetSummarySummary = application
			requirementsFactory.Application = application
		})

		It("displays nice output when the app is stopped", func() {
			appSummaryRepo.GetSummaryErrorCode = errors.APP_STOPPED
			runCommand("my-app")

			Expect(appSummaryRepo.GetSummaryAppGuid).To(Equal("my-app-guid"))
			Expect(appInstancesRepo.GetInstancesAppGuid).To(Equal("my-app-guid"))

			testassert.SliceContains(ui.Outputs, testassert.Lines{
				{"Showing health and status", "my-app", "my-org", "my-space", "my-user"},
				{"state", "stopped"},
				{"instances", "0/2"},
				{"usage", "256M x 2 instances"},
				{"no running instances"},
			})
		})

		It("displays nice output when the app has not yet finished staging", func() {
			appSummaryRepo.GetSummaryErrorCode = errors.APP_NOT_STAGED
			runCommand("my-app")

			Expect(appSummaryRepo.GetSummaryAppGuid).To(Equal("my-app-guid"))
			Expect(appInstancesRepo.GetInstancesAppGuid).To(Equal("my-app-guid"))

			testassert.SliceContains(ui.Outputs, testassert.Lines{
				{"Showing health and status", "my-app", "my-org", "my-space", "my-user"},
				{"state", "stopped"},
				{"instances", "0/2"},
				{"usage", "256M x 2 instances"},
				{"no running instances"},
			})
		})
	})
})

func makeAppWithRoute(appName string) models.Application {
	application := models.Application{}
	application.Name = appName
	application.Guid = "app-guid"

	domain := models.DomainFields{}
	domain.Name = "example.com"

	route := models.RouteSummary{Host: "foo", Domain: domain}
	secondRoute := models.RouteSummary{Host: appName, Domain: domain}

	application.State = "started"
	application.InstanceCount = 2
	application.RunningInstances = 2
	application.Memory = 256
	application.Routes = []models.RouteSummary{route, secondRoute}

	return application
}
