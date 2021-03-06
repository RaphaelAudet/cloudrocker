package stager_test

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudcredo/cloudrocker/Godeps/_workspace/src/github.com/cloudfoundry-incubator/buildpack_app_lifecycle/buildpackrunner"
	. "github.com/cloudcredo/cloudrocker/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/cloudcredo/cloudrocker/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/cloudcredo/cloudrocker/Godeps/_workspace/src/github.com/onsi/gomega/gbytes"
	"github.com/cloudcredo/cloudrocker/config"
	"github.com/cloudcredo/cloudrocker/stager"
)

type TestRunner struct {
	RunCalled bool
}

func (f *TestRunner) Run() error {
	f.RunCalled = true
	return nil
}

var _ = Describe("Stager", func() {
	Describe("Running a buildpack", func() {
		It("should tell a buildpack runner to run", func() {
			buffer := gbytes.NewBuffer()
			testrunner := new(TestRunner)
			err := stager.RunBuildpack(buffer, testrunner)
			Expect(testrunner.RunCalled).To(Equal(true))
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(buffer).Should(gbytes.Say(`Running Buildpacks...`))
		})
	})

	Describe("Getting a buildpack runner", func() {
		It("should return the address of a valid buildpack runner, with a correct buildpack list", func() {
			buildpackDir, _ := ioutil.TempDir(os.TempDir(), "crocker-buildpackrunner-test")
			os.Mkdir(buildpackDir+"/test-buildpack", 0755)
			ioutil.WriteFile(buildpackDir+"/test-buildpack"+"/testfile", []byte("test"), 0644)
			runner := stager.NewBuildpackRunner(buildpackDir)
			var runnerVar *buildpackrunner.Runner
			Expect(runner).Should(BeAssignableToTypeOf(runnerVar))
			md5BuildpackName := fmt.Sprintf("%x", md5.Sum([]byte("test-buildpack")))
			md5BuildpackDir, err := os.Open("/tmp/buildpacks")
			contents, err := md5BuildpackDir.Readdirnames(0)
			Expect(contents, err).Should(ContainElement(md5BuildpackName))
			md5Buildpack, err := os.Open(md5BuildpackDir.Name() + "/" + md5BuildpackName)
			buildpackContents, err := md5Buildpack.Readdirnames(0)
			Expect(buildpackContents, err).Should(ContainElement("testfile"))
			os.RemoveAll(buildpackDir)
			os.RemoveAll("/tmp/buildpacks")
		})
	})

	Describe("Validating a staged application", func() {
		var cfhome string
		BeforeEach(func() {
			cfhome, _ = ioutil.TempDir(os.TempDir(), "stager-test-staged")
		})
		AfterEach(func() {
			os.RemoveAll(cfhome)
		})

		Context("with something that looks like a staged application", func() {
			It("should not return an error", func() {
				dropletDir := config.NewDirectories(cfhome).Tmp()
				os.MkdirAll(dropletDir+"/tmp", 0755)
				ioutil.WriteFile(dropletDir+"/result.json", []byte("test-staging-info"), 0644)
				ioutil.WriteFile(dropletDir+"/droplet", []byte("test-droplet"), 0644)
				err := stager.ValidateStagedApp(config.NewDirectories(cfhome))
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("without something that looks like a staged application", func() {
			Context("because we have no droplet", func() {
				It("should return an error about a missing droplet", func() {
					err := stager.ValidateStagedApp(config.NewDirectories(cfhome))
					Expect(err).Should(MatchError("Staging failed - have you added a buildpack for this type of application?"))
				})
			})
			Context("because we have no staging_info.yml", func() {
				It("should return an error about missing staging info", func() {
					os.MkdirAll(cfhome+"/tmp/droplet/app", 0755)
					err := stager.ValidateStagedApp(config.NewDirectories(cfhome))
					Expect(err).Should(MatchError("Staging failed - no result json was produced by the matching buildpack!"))
				})
			})
		})
	})
})
