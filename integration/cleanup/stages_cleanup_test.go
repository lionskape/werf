// +build integration

package cleanup_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/flant/werf/integration/utils"
)

var _ = Describe("stages cleanup command", func() {
	var testDirPath string

	BeforeEach(func() {
		testDirPath = tmpPath()
		utils.CopyIn(fixturePath("stages_cleanup"), testDirPath)

		utils.RunSucceedCommand(
			testDirPath,
			"git",
			"init",
		)

		utils.RunSucceedCommand(
			testDirPath,
			"git",
			"add", "werf.yaml",
		)

		utils.RunSucceedCommand(
			testDirPath,
			"git",
			"commit", "-m", "Initial commit",
		)

		Ω(os.Setenv("WERF_IMAGES_REPO", registryProjectRepository)).Should(Succeed())
		Ω(os.Setenv("WERF_STAGES_STORAGE", ":local")).Should(Succeed())

		Ω(os.Setenv("FROM_CACHE_VERSION", "x")).Should(Succeed())
	})

	AfterEach(func() {
		utils.RunSucceedCommand(
			testDirPath,
			werfBinPath,
			"stages", "purge", "-s", ":local", "--force",
		)
	})

	It("should work properly with non-existent registry repository", func() {
		utils.RunSucceedCommand(
			testDirPath,
			werfBinPath,
			"stages", "cleanup",
		)
	})

	for _, disableStageCleanupDatePeriodPolicy := range []string{"0", "1"} {
		if disableStageCleanupDatePeriodPolicy == "1" {
			It("should not remove stages images related with built images in repository (WERF_DISABLE_STAGES_CLEANUP_DATE_PERIOD_POLICY=1)", func() {
				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"build-and-publish", "--tag-git-branch", "test",
				)

				count := LocalProjectStagesCount()
				Ω(count).Should(Equal(count))

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"stages", "cleanup",
				)

				Ω(LocalProjectStagesCount()).Should(Equal(count))
			})

			Context("when there is running container based on werf image", func() {
				BeforeEach(func() {
					utils.RunSucceedCommand(
						testDirPath,
						werfBinPath,
						"build",
					)

					utils.RunSucceedCommand(
						testDirPath,
						werfBinPath,
						"run", "--docker-options", "-d", "--", "/bin/sleep", "30",
					)

					Ω(os.Setenv("WERF_LOG_PRETTY", "0")).Should(Succeed())
				})

				It("should skip stage image with related running container", func() {
					out, err := utils.RunCommand(
						testDirPath,
						werfBinPath,
						"stages", "cleanup",
					)
					Ω(err).Should(Succeed())
					Ω(string(out)).Should(ContainSubstring("Skip image "))
					Ω(string(out)).Should(ContainSubstring("used by container"))
				})
			})
		}

		boundedPolicyValue := disableStageCleanupDatePeriodPolicy

		var itMsg string
		if disableStageCleanupDatePeriodPolicy == "0" {
			itMsg = fmt.Sprintf("should not remove unused stages images (WERF_DISABLE_STAGES_CLEANUP_DATE_PERIOD_POLICY=0)")
		} else {
			itMsg = fmt.Sprintf("should remove unused stages images (WERF_DISABLE_STAGES_CLEANUP_DATE_PERIOD_POLICY=1)")
		}

		It(itMsg, func() {
			Ω(os.Setenv("WERF_DISABLE_STAGES_CLEANUP_DATE_PERIOD_POLICY", boundedPolicyValue)).Should(Succeed())

			utils.RunSucceedCommand(
				testDirPath,
				werfBinPath,
				"build-and-publish", "--tag-git-branch", "test",
			)

			countAfterFirstBuild := LocalProjectStagesCount()
			Ω(countAfterFirstBuild).Should(Equal(countAfterFirstBuild))

			Ω(os.Setenv("FROM_CACHE_VERSION", "fully rebuild")).Should(Succeed())

			utils.RunSucceedCommand(
				testDirPath,
				werfBinPath,
				"build-and-publish", "--tag-git-branch", "test",
			)

			countAfterSecondBuild := LocalProjectStagesCount()

			utils.RunSucceedCommand(
				testDirPath,
				werfBinPath,
				"stages", "cleanup",
			)

			Ω(LocalProjectStagesCount()).Should(Equal(countAfterSecondBuild - countAfterFirstBuild))
		})
	}
})
