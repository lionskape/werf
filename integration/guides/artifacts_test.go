// +build integration

package guides_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"

	"github.com/flant/werf/integration/utils"
	utilsDocker "github.com/flant/werf/integration/utils/docker"
)

var _ = Describe("Advanced build/Artifacts", func() {
	var testDirPath string

	BeforeEach(func() {
		testDirPath = tmpPath()
		utils.CopyIn(fixturePath("artifacts"), testDirPath)
	})

	AfterEach(func() {
		utils.RunSucceedCommand(
			testDirPath,
			werfBinPath,
			"stages", "purge", "-s", ":local", "--force",
		)
	})

	It("application should be built and checked", func() {
		utils.RunSucceedCommand(
			testDirPath,
			werfBinPath,
			"build", "-s", ":local",
		)

		containerHostPort := utils.GetFreeTCPHostPort()
		containerName := fmt.Sprintf("go_booking_artifacts_%s", utils.GetRandomString(10))
		utils.RunSucceedCommand(
			testDirPath,
			werfBinPath,
			"run", "-s", ":local", "--docker-options", fmt.Sprintf("-d -p %d:9000 --name %s", containerHostPort, containerName), "go-booking", "--", "/app/run.sh",
		)
		defer func() { utilsDocker.ContainerStopAndRemove(containerName) }()

		url := fmt.Sprintf("http://localhost:%d", containerHostPort)
		waitTillHostReadyAndCheckResponseBody(
			url,
			utils.DefaultWaitTillHostReadyToRespondMaxAttempts,
			"revel framework booking demo",
		)
	})
})
