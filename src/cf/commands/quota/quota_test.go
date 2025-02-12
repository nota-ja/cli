package quota_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "cf/commands/quota"
	"cf/errors"
	"cf/models"
	testapi "testhelpers/api"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"

	. "testhelpers/matchers"
)

var _ = Describe("quota", func() {
	var (
		ui                  *testterm.FakeUI
		requirementsFactory *testreq.FakeReqFactory
		quotaRepo           *testapi.FakeQuotaRepository
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		requirementsFactory = &testreq.FakeReqFactory{}
		quotaRepo = &testapi.FakeQuotaRepository{}
	})

	runCommand := func(args ...string) {
		cmd := NewShowQuota(ui, testconfig.NewRepositoryWithDefaults(), quotaRepo)
		testcmd.RunCommand(cmd, testcmd.NewContext("quotas", args), requirementsFactory)
	}

	Context("When not logged in", func() {
		It("fails requirements", func() {
			runCommand("quota-name")
			Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
		})
	})

	Context("When logged in", func() {
		BeforeEach(func() {
			requirementsFactory.LoginSuccess = true
		})

		Context("When not providing a quota name", func() {
			It("fails with usage", func() {
				runCommand()
				Expect(ui.FailedWithUsage).To(BeTrue())
			})
		})

		Context("When providing a quota name", func() {
			Context("that exists", func() {
				BeforeEach(func() {
					quotaRepo.FindByNameReturns.Quota = models.QuotaFields{
						Guid:                    "my-quota-guid",
						Name:                    "muh-muh-muh-my-qua-quota",
						MemoryLimit:             512,
						RoutesLimit:             2000,
						ServicesLimit:           47,
						NonBasicServicesAllowed: true,
					}
				})

				It("shows you that quota", func() {
					runCommand("muh-muh-muh-my-qua-quota")

					Expect(ui.Outputs).To(ContainSubstrings(
						[]string{"Getting quota", "muh-muh-muh-my-qua-quota", "my-user"},
						[]string{"OK"},
						[]string{"Memory", "512M"},
						[]string{"Routes", "2000"},
						[]string{"Services", "47"},
						[]string{"Paid service plans", "allowed"},
					))
				})
			})

			Context("that doesn't exist", func() {
				BeforeEach(func() {
					quotaRepo.FindByNameReturns.Error = errors.New("oops i accidentally a quota")
				})

				It("gives an error", func() {
					runCommand("an-quota")

					Expect(ui.Outputs).To(ContainSubstrings(
						[]string{"FAILED"},
						[]string{"oops"},
					))
				})
			})
		})
	})
})
