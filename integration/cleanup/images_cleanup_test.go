// +build integration

package cleanup_test

import (
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/flant/werf/integration/utils"
)

var _ = Describe("images cleanup command", func() {
	var testDirPath string

	BeforeEach(func() {
		testDirPath = tmpPath()
		utils.CopyIn(fixturePath("default"), testDirPath)

		utils.RunSucceedCommand(
			testDirPath,
			"git",
			"init", "--bare", "remote.git",
		)

		utils.RunSucceedCommand(
			testDirPath,
			"git",
			"init",
		)

		utils.RunSucceedCommand(
			testDirPath,
			"git",
			"remote", "add", "origin", "remote.git",
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
	})

	AfterEach(func() {
		utils.RunSucceedCommand(
			testDirPath,
			werfBinPath,
			"stages", "purge", "-s", ":local", "--force",
		)
	})

	Context("when deployed images in kubernetes are not taken in account", func() {
		BeforeEach(func() {
			Ω(os.Setenv("WERF_WITHOUT_KUBE", "1")).Should(Succeed())
		})

		It("should work properly with non-existent registry repository", func() {
			utils.RunSucceedCommand(
				testDirPath,
				werfBinPath,
				"images", "cleanup",
			)
		})

		Context("git branch strategy", func() {
			var testBranch = "branchA"

			It("should remove image associated with local branch", func() {
				utils.RunSucceedCommand(
					testDirPath,
					"git",
					"checkout", "-b", testBranch,
				)

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"build-and-publish", "--tag-git-branch", testBranch,
				)

				tags := utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).Should(ContainElement(testBranch))

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"images", "cleanup",
				)

				tags = utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).ShouldNot(ContainElement(testBranch))
			})

			It("should remove image associated with deleted remote branch", func() {
				utils.RunSucceedCommand(
					testDirPath,
					"git",
					"checkout", "-b", testBranch,
				)

				utils.RunSucceedCommand(
					testDirPath,
					"git",
					"push", "--set-upstream", "origin", testBranch,
				)

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"build-and-publish", "--tag-git-branch", testBranch,
				)

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"images", "cleanup",
				)

				tags := utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).Should(ContainElement(testBranch))

				utils.RunSucceedCommand(
					testDirPath,
					"git",
					"push", "origin", "--delete", testBranch,
				)

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"images", "cleanup",
				)

				tags = utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).ShouldNot(ContainElement(testBranch))
			})
		})

		Context("git tag strategy", func() {
			var testTag = "tagA"

			It("should remove image associated with deleted tag", func() {
				utils.RunSucceedCommand(
					testDirPath,
					"git",
					"tag", testTag,
				)

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"build-and-publish", "--tag-git-tag", testTag,
				)

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"images", "cleanup",
				)

				tags := utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).Should(ContainElement(testTag))

				utils.RunSucceedCommand(
					testDirPath,
					"git",
					"tag", "-d", testTag,
				)

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"images", "cleanup",
				)

				tags = utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).ShouldNot(ContainElement(testTag))
			})

			It("should remove image by expiry days policy (WERF_GIT_TAG_STRATEGY_EXPIRY_DAYS)", func() {
				utils.RunSucceedCommand(
					testDirPath,
					"git",
					"tag", testTag,
				)

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"build-and-publish", "--tag-git-tag", testTag,
				)

				tags := utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).Should(HaveLen(1))

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"images", "cleanup", "--git-tag-strategy-expiry-days", "0",
				)

				tags = utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).Should(HaveLen(0))
			})

			It("should remove image by limit policy (WERF_GIT_TAG_STRATEGY_LIMIT)", func() {
				for _, tag := range []string{"tagA", "tagB", "tagC"} {
					utils.RunSucceedCommand(
						testDirPath,
						"git",
						"tag", tag,
					)

					utils.RunSucceedCommand(
						testDirPath,
						werfBinPath,
						"build-and-publish", "--tag-git-tag", tag,
					)
				}

				tags := utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).Should(HaveLen(3))

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"images", "cleanup", "--git-tag-strategy-limit", "1",
				)

				tags = utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).Should(HaveLen(1))
			})
		})

		Context("git commit strategy", func() {
			It("should remove image that associated with non-existent commit", func() {
				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"build-and-publish", "--tag-git-commit", "8a99331ce0f918b411423223f4060e9688e03f6a",
				)

				tags := utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).Should(HaveLen(1))

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"images", "cleanup",
				)

				tags = utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).Should(HaveLen(0))
			})

			It("should not remove image that associated with commit", func() {
				out := utils.SucceedCommandOutput(
					testDirPath,
					"git",
					"rev-parse", "HEAD",
				)
				commit := strings.TrimSpace(out)

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"build-and-publish", "--tag-git-commit", commit,
				)

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"images", "cleanup",
				)

				tags := utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).Should(ContainElement(commit))
			})

			It("should remove image by expiry days policy (WERF_GIT_COMMIT_STRATEGY_EXPIRY_DAYS)", func() {
				out := utils.SucceedCommandOutput(
					testDirPath,
					"git",
					"rev-parse", "HEAD",
				)
				commit := strings.TrimSpace(out)

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"build-and-publish", "--tag-git-commit", commit,
				)

				tags := utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).Should(HaveLen(1))

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"images", "cleanup", "--git-commit-strategy-expiry-days", "0",
				)

				tags = utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).Should(HaveLen(0))
			})

			It("should remove image by limit policy (WERF_GIT_COMMIT_STRATEGY_LIMIT)", func() {
				amount := 4
				for i := 0; i < amount; i++ {
					utils.RunSucceedCommand(
						testDirPath,
						"git",
						"commit", "--allow-empty", "--allow-empty-message", "-m", "",
					)

					out := utils.SucceedCommandOutput(
						testDirPath,
						"git",
						"rev-parse", "HEAD",
					)
					commit := strings.TrimSpace(out)

					utils.RunSucceedCommand(
						testDirPath,
						werfBinPath,
						"build-and-publish", "--tag-git-commit", commit,
					)
				}

				tags := utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).Should(HaveLen(amount))

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"images", "cleanup", "--git-commit-strategy-limit", "2",
				)

				tags = utils.RegistryRepositoryList(registryProjectRepository)
				Ω(tags).Should(HaveLen(2))
			})
		})
	})
})
