// +build integration

package git

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/alessio/shellescape"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	"github.com/flant/werf/integration/utils"
)

var _ = Describe("cleanup empty directories with git patch apply", func() {
	var testDirPath string
	var fixturesPathParts []string
	gitToPath := "/app"

	type removingEmptyDirectoriesEntry struct {
		dirToAdd        string
		shouldBeDeleted []string
		shouldBeSkipped []string
	}

	removingEmptyDirectoriesItBody := func(fixturePathFolder string) func(removingEmptyDirectoriesEntry) {
		return func(entry removingEmptyDirectoriesEntry) {
			commonBeforeEach(testDirPath, fixturePath(append(fixturesPathParts, fixturePathFolder)...))

			projectAddedFilePath := filepath.Join(entry.dirToAdd, "file")
			containerAddedDirPath := path.Join(gitToPath, entry.dirToAdd)

			By(fmt.Sprintf("Add file %s", shellescape.Quote(projectAddedFilePath)))
			createAndCommitFile(filepath.Join(testDirPath, entry.dirToAdd), "file", 12)

			By("Build and cache source code in gitArchive stage")
			utils.RunSucceedCommand(
				testDirPath,
				werfBinPath,
				"build",
			)

			By(fmt.Sprintf("Check container directory %s exists", shellescape.Quote(containerAddedDirPath)))
			checkContainerDirectoryExists(testDirPath, containerAddedDirPath)

			By(fmt.Sprintf("Remove file %s", shellescape.Quote(projectAddedFilePath)))

			utils.RunSucceedCommand(
				testDirPath,
				"git",
				"rm", projectAddedFilePath,
			)

			utils.RunSucceedCommand(
				testDirPath,
				"git",
				"commit", "-m", "Remove file "+projectAddedFilePath,
			)

			utils.RunSucceedCommand(
				testDirPath,
				werfBinPath,
				"build",
			)

			for _, relDirPath := range entry.shouldBeDeleted {
				containerDirPath := path.Join(gitToPath, relDirPath)
				By(fmt.Sprintf("Check container directory %s does not exist", shellescape.Quote(containerDirPath)))
				checkContainerDirectoryDoesNotExist(testDirPath, containerDirPath)
			}

			for _, relDirPath := range entry.shouldBeSkipped {
				containerDirPath := path.Join(gitToPath, relDirPath)
				By(fmt.Sprintf("Check container directory %s exists", shellescape.Quote(containerDirPath)))
				checkContainerDirectoryExists(testDirPath, containerDirPath)
			}
		}
	}

	BeforeEach(func() {
		testDirPath = tmpPath()
		fixturesPathParts = []string{"cleanup_empty_directories_with_git_patch_apply"}
	})

	AfterEach(func() {
		utils.RunSucceedCommand(
			testDirPath,
			werfBinPath,
			"stages", "purge", "-s", ":local", "--force",
		)
	})

	DescribeTable("base",
		removingEmptyDirectoriesItBody("base"),
		Entry("should remove empty directory (dir)", removingEmptyDirectoriesEntry{
			dirToAdd:        "dir",
			shouldBeDeleted: []string{"dir"},
			shouldBeSkipped: []string{},
		}),
		Entry("should remove empty directories (dir/sub_dir)", removingEmptyDirectoriesEntry{
			dirToAdd:        "dir/sub_dir",
			shouldBeDeleted: []string{"dir/sub_dir", "dir"},
			shouldBeSkipped: []string{},
		}),
		Entry("should remove empty directories (dir/sub dir/sub dir with special ch@ra(c)ters? ())", removingEmptyDirectoriesEntry{
			dirToAdd:        "dir/sub dir/sub dir with special ch@ra(c)ters? ()",
			shouldBeDeleted: []string{"dir/sub dir/sub dir with special ch@ra(c)ters? ()", "dir/sub dir", "dir"},
			shouldBeSkipped: []string{},
		}),
	)

	DescribeTable("processing directory created by user 'dir/dir_created_by_user'",
		removingEmptyDirectoriesItBody("skipping_user_directory"),
		Entry("should not remove directory (dir)", removingEmptyDirectoriesEntry{
			dirToAdd:        "dir",
			shouldBeDeleted: []string{},
			shouldBeSkipped: []string{"dir"},
		}),
		Entry("should remove only empty directory (dir/sub_dir)", removingEmptyDirectoriesEntry{
			dirToAdd:        "dir/sub_dir",
			shouldBeDeleted: []string{"dir/sub_dir"},
			shouldBeSkipped: []string{"dir"},
		}),
		Entry("should remove empty directories (dir/dir_created_by_user)", removingEmptyDirectoriesEntry{
			dirToAdd:        "dir/dir_created_by_user",
			shouldBeDeleted: []string{"dir/dir_created_by_user", "dir"},
			shouldBeSkipped: []string{},
		}),
	)
})

func checkContainerDirectoryExists(projectPath, containerDirPath string) {
	checkContainerDirectory(projectPath, containerDirPath, true)
}

func checkContainerDirectoryDoesNotExist(projectPath, containerDirPath string) {
	checkContainerDirectory(projectPath, containerDirPath, false)
}

func checkContainerDirectory(projectPath, containerDirPath string, exist bool) {
	cmd := checkContainerFileCommand(containerDirPath, true, exist)

	utils.RunSucceedCommand(
		projectPath,
		werfBinPath,
		"run", "--docker-options", "--rm", "--", "bash", "-ec", cmd,
	)
}

func checkContainerFileCommand(containerDirPath string, directory bool, exist bool) string {
	var cmd string
	var flag string

	if directory {
		flag = "-d"
	} else {
		flag = "-f"
	}

	if exist {
		cmd = fmt.Sprintf("test %s %s", flag, shellescape.Quote(containerDirPath))
	} else {
		cmd = fmt.Sprintf("! test %s %s", flag, shellescape.Quote(containerDirPath))
	}

	return cmd
}
