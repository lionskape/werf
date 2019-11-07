// +build integration

package git

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alessio/shellescape"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/flant/werf/integration/utils"
)

var _ = Describe("file lifecycle", func() {
	var testDirPath string
	var fixturesPathParts []string
	gitToPath := "/app"

	fileNameToAdd := "test"
	fileNameToAddAndModify := "test2"
	fileDataToAdd := []byte("test")
	fileDataToModify := []byte("test2")

	type fileLifecycleEntry struct {
		name   string
		data   []byte
		perm   os.FileMode
		delete bool
	}

	createFileFunc := func(filePath string, fileData []byte, filePerm os.FileMode) {
		utils.CreateFile(filePath, fileData)
		Ω(os.Chmod(filePath, filePerm)).Should(Succeed())
	}

	fileLifecycleEntryItBody := func(entry fileLifecycleEntry) {
		var commitMsg string

		filePath := filepath.Join(testDirPath, entry.name)
		if entry.delete {
			Ω(os.Remove(filePath)).Should(Succeed())
			commitMsg = "Delete file " + entry.name
		} else {
			createFileFunc(filePath, entry.data, entry.perm)
			commitMsg = "Add/Modify file " + entry.name
		}

		addAndCommitFile(testDirPath, entry.name, commitMsg)

		utils.RunSucceedCommand(
			testDirPath,
			werfBinPath,
			"build",
		)

		var cmd []string
		dockerOptions := []string{"--rm"}

		if entry.delete {
			cmd = append(cmd, checkContainerFileCommand(path.Join(gitToPath, entry.name), false, false))
		} else {
			cmd = append(cmd, checkContainerFileCommand(path.Join(gitToPath, entry.name), false, true))
			cmd = append(cmd, fmt.Sprintf("diff <(stat -c %%a %s) <(echo %s)", path.Join(gitToPath, entry.name), strconv.FormatUint(uint64(entry.perm), 8)))
			cmd = append(cmd, fmt.Sprintf("diff %s %s", path.Join(gitToPath, entry.name), "/source"))

			dockerOptions = append(dockerOptions, fmt.Sprintf("-v %s:%s", filePath, "/source"))
		}

		utils.RunSucceedCommand(
			testDirPath,
			werfBinPath,
			"run", "--docker-options", strings.Join(dockerOptions, " "), "--", "bash", "-ec", strings.Join(cmd, " && "),
		)
	}

	BeforeEach(func() {
		testDirPath = tmpPath()
		fixturesPathParts = []string{"file_lifecycle"}
		commonBeforeEach(testDirPath, fixturePath(fixturesPathParts...))
	})

	AfterEach(func() {
		utils.RunSucceedCommand(
			testDirPath,
			werfBinPath,
			"stages", "purge", "-s", ":local", "--force",
		)
	})

	DescribeTable("processing file with archive apply",
		fileLifecycleEntryItBody,
		Entry("should add file (0755)", fileLifecycleEntry{
			name: fileNameToAdd,
			data: fileDataToAdd,
			perm: 0755,
		}),
		Entry("should add file (0644)", fileLifecycleEntry{
			name: fileNameToAdd,
			data: fileDataToAdd,
			perm: 0644,
		}),
	)

	Context("when gitArchive stage with file is built", func() {
		BeforeEach(func() {
			createFileFunc(filepath.Join(testDirPath, fileNameToAddAndModify), fileDataToAdd, 0755)
			addAndCommitFile(testDirPath, fileNameToAddAndModify, "Add file "+fileNameToAddAndModify)

			utils.RunSucceedCommand(
				testDirPath,
				werfBinPath,
				"build",
			)
		})

		DescribeTable("processing file with patch apply",
			fileLifecycleEntryItBody,
			Entry("should add file (0755)", fileLifecycleEntry{
				name: fileNameToAdd,
				data: fileDataToAdd,
				perm: 0755,
			}),
			Entry("should add file (0644)", fileLifecycleEntry{
				name: fileNameToAdd,
				data: fileDataToAdd,
				perm: 0644,
			}),
			Entry("should modify file", fileLifecycleEntry{
				name: fileNameToAddAndModify,
				data: fileDataToModify,
				perm: 0755,
			}),
			Entry("should change file permission (0755->0644)", fileLifecycleEntry{
				name: fileNameToAddAndModify,
				data: fileDataToAdd,
				perm: 0644,
			}),
			Entry("should modify and change file permission (0755->0644)", fileLifecycleEntry{
				name: fileNameToAddAndModify,
				data: fileDataToModify,
				perm: 0644,
			}),
			Entry("should delete file", fileLifecycleEntry{
				name:   fileNameToAddAndModify,
				delete: true,
			}),
		)
	})

	Context("when file is symlink", func() {
		linkToAdd := "werf.yaml"
		linkToModify := "none"

		type symlinkFileLifecycleEntry struct {
			name   string
			link   string
			delete bool
		}

		symlinkFileLifecycleEntryItBody := func(entry symlinkFileLifecycleEntry) {
			var commitMsg string

			filePath := filepath.Join(testDirPath, entry.name)
			if entry.delete {
				Ω(os.Remove(filePath)).Should(Succeed())
				commitMsg = "Delete file " + entry.name
			} else {
				if _, err := os.Lstat(filePath); err == nil {
					Ω(os.Remove(filePath)).Should(Succeed())
				}

				Ω(os.Symlink(entry.link, filePath)).Should(Succeed())
				commitMsg = "Add/Modify file " + entry.name
			}

			addAndCommitFile(testDirPath, entry.name, commitMsg)

			utils.RunSucceedCommand(
				testDirPath,
				werfBinPath,
				"build",
			)

			var cmd []string
			if entry.delete {
				cmd = append(cmd, checkContainerSymlinkFileCommand(path.Join(gitToPath, entry.name), false))
			} else {
				cmd = append(cmd, checkContainerSymlinkFileCommand(path.Join(gitToPath, entry.name), true))
				readlinkCmd := fmt.Sprintf("readlink %s", path.Join(gitToPath, entry.name))
				cmd = append(cmd, fmt.Sprintf("diff <(%s) <(echo %s)", readlinkCmd, entry.link))
			}

			utils.RunSucceedCommand(
				testDirPath,
				werfBinPath,
				"run", "--docker-options", "--rm", "--", "bash", "-ec", strings.Join(cmd, " && "),
			)
		}

		DescribeTable("processing symlink file with archive apply",
			symlinkFileLifecycleEntryItBody,
			Entry("should add symlink", symlinkFileLifecycleEntry{
				name: fileNameToAdd,
				link: linkToAdd,
			}),
		)

		Context("when gitArchive stage with file is built", func() {
			BeforeEach(func() {
				Ω(os.Symlink(linkToAdd, filepath.Join(testDirPath, fileNameToAddAndModify))).Should(Succeed())
				addAndCommitFile(testDirPath, fileNameToAddAndModify, "Add file "+fileNameToAddAndModify)

				utils.RunSucceedCommand(
					testDirPath,
					werfBinPath,
					"build",
				)
			})

			DescribeTable("processing symlink file with patch apply",
				symlinkFileLifecycleEntryItBody,
				Entry("should add symlink", symlinkFileLifecycleEntry{
					name: fileNameToAdd,
					link: linkToAdd,
				}),
				Entry("should modify file", symlinkFileLifecycleEntry{
					name: fileNameToAddAndModify,
					link: linkToModify,
				}),
				Entry("should delete file", symlinkFileLifecycleEntry{
					name:   fileNameToAddAndModify,
					delete: true,
				}))
		})
	})
})

func checkContainerSymlinkFileCommand(containerDirPath string, exist bool) string {
	var cmd string

	if exist {
		cmd = fmt.Sprintf("test -h %s", shellescape.Quote(containerDirPath))
	} else {
		cmd = fmt.Sprintf("! test -h %s", shellescape.Quote(containerDirPath))
	}

	return cmd
}
