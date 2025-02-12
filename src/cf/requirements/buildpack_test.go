/*
                       WARNING WARNING WARNING

                Attention all potential contributors

   This testfile is not in the best state. We've been slowly transitioning
   from the built in "testing" package to using Ginkgo. As you can see, we've
   changed the format, but a lot of the setup, test body, descriptions, etc
   are either hardcoded, completely lacking, or misleading.

   For example:

   Describe("Testing with ginkgo"...)      // This is not a great description
   It("TestDoesSoemthing"...)              // This is a horrible description

   Describe("create-user command"...       // Describe the actual object under test
   It("creates a user when provided ..."   // this is more descriptive

   For good examples of writing Ginkgo tests for the cli, refer to

   src/cf/commands/application/delete_app_test.go
   src/cf/terminal/ui_test.go
   src/github.com/cloudfoundry/loggregator_consumer/consumer_test.go
*/

package requirements_test

import (
	"cf/models"
	. "cf/requirements"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testapi "testhelpers/api"
	testassert "testhelpers/assert"
	testterm "testhelpers/terminal"
)

var _ = Describe("Testing with ginkgo", func() {
	It("TestBuildpackReqExecute", func() {

		buildpack := models.Buildpack{}
		buildpack.Name = "my-buildpack"
		buildpack.Guid = "my-buildpack-guid"
		buildpackRepo := &testapi.FakeBuildpackRepository{FindByNameBuildpack: buildpack}
		ui := new(testterm.FakeUI)

		buildpackReq := NewBuildpackRequirement("foo", ui, buildpackRepo)
		success := buildpackReq.Execute()

		Expect(success).To(BeTrue())
		Expect(buildpackRepo.FindByNameName).To(Equal("foo"))
		Expect(buildpackReq.GetBuildpack()).To(Equal(buildpack))
	})
	It("TestBuildpackReqExecuteWhenBuildpackNotFound", func() {

		buildpackRepo := &testapi.FakeBuildpackRepository{FindByNameNotFound: true}
		ui := new(testterm.FakeUI)

		buildpackReq := NewBuildpackRequirement("foo", ui, buildpackRepo)

		testassert.AssertPanic(testterm.FailedWasCalled, func() {
			buildpackReq.Execute()
		})
	})
})
