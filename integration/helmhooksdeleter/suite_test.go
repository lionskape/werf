// +build integration integration_k8s

package helmhooksdeleter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/flant/kubedog/pkg/kube"
	"github.com/flant/werf/integration/utils"
	"github.com/flant/werf/integration/utils/werfexec"
	"sigs.k8s.io/yaml"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func unmarshalObject(manifestYaml string, obj interface{}) {
	manifestJson, err := yaml.YAMLToJSON([]byte(manifestYaml))
	Expect(err).To(Succeed())
	Expect(json.Unmarshal(manifestJson, obj)).To(Succeed())
}

func NewNamespace(manifestYaml string) *corev1.Namespace {
	obj := &corev1.Namespace{}
	unmarshalObject(manifestYaml, &obj)
	return obj
}

func NewJob(manifestYaml string) *batchv1.Job {
	obj := &batchv1.Job{}
	unmarshalObject(manifestYaml, &obj)
	return obj
}

var _ = Describe("Helm hooks deleter", func() {
	Context("when installing chart with post-install Job hook and hook-succeeded delete policy", func() {
		AfterEach(func() {
			werfDismiss("app1", werfexec.CommandOptions{})
		})

		It("should delete hook and wait till hook deleted without timeout https://github.com/flant/werf/issues/1885", func(done Done) {
			gotDeletingHookLine := false

			Expect(werfDeploy("app1", werfexec.CommandOptions{
				OutputLineHandler: func(line string) {
					Expect(strings.HasPrefix(line, "│ NOTICE Will not delete Job/migrate: resource does not belong to the helm release")).ShouldNot(BeTrue(), fmt.Sprintf("Got unexpected output line: %v", line))

					if strings.HasPrefix(line, "│ Deleting resource Job/migrate from release") {
						gotDeletingHookLine = true
					}
				},
			})).Should(Succeed())

			Expect(gotDeletingHookLine).Should(BeTrue())

			close(done)
		}, 120)
	})

	Context("when releasing a chart containing a hook with before-hook-creation delete policy and the hook already exists in the cluster before release", func() {
		var namespace string
		var err error

		BeforeEach(func() {
			namespace = fmt.Sprintf("%s-dev", utils.ProjectName())
		})

		BeforeEach(func() {
			Expect(kube.Init(kube.InitOptions{})).To(Succeed())
		})

		AfterEach(func() {
			werfDismiss("app2", werfexec.CommandOptions{})
		})

		FIt("", func(done Done) {
			_, err = kube.Kubernetes.CoreV1().Namespaces().Create(NewNamespace(fmt.Sprintf(`
apiVersion: core/v1
kind: Namespace
metadata:
  name: %s
`, namespace)))
			Expect(err).NotTo(HaveOccurred())

			originalJob, err := kube.Kubernetes.BatchV1().Jobs(namespace).Create(NewJob(`
apiVersion: batch/v1
kind: Job
metadata:
  name: migrate
  annotations:
    "helm.sh/hook": post-upgrade
    "helm.sh/hook-delete-policy": before-hook-creation
    "helm.sh/hook-weight": "10"
    "werf.io/recreate": "false"
spec:
  backoffLimit: 1
  activeDeadlineSeconds: 600
  template:
    metadata:
      name: migrate
    spec:
      restartPolicy: Never
      containers:
      - name: main
        image: ubuntu:18.04
        command: [ "/bin/bash", "-lec", "for i in {1..3}; do date; sleep 1; done" ]
`))
			Expect(err).NotTo(HaveOccurred())

			// Install release, hook should not be touched because it is post-update hook
			Expect(werfDeploy("app2", werfexec.CommandOptions{})).Should(Succeed())

			jobAfterInstall, err := kube.Kubernetes.BatchV1().Jobs(namespace).Get(originalJob.Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(jobAfterInstall.UID).To(Equal(originalJob.UID))

			// Update release, hook should be deleted by before-hook-creation policy and created again
			Expect(werfDeploy("app2", werfexec.CommandOptions{})).Should(Succeed())
			jobAfterUpdate, err := kube.Kubernetes.BatchV1().Jobs(namespace).Get(originalJob.Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(jobAfterUpdate.UID).NotTo(Equal(originalJob.UID))

			close(done)
		}, 120)
	})
})

func werfDeploy(dir string, opts werfexec.CommandOptions) error {
	return werfexec.ExecWerfCommand(dir, werfBinPath, opts, "deploy", "--env", "dev")
}

func werfDismiss(dir string, opts werfexec.CommandOptions) error {
	return werfexec.ExecWerfCommand(dir, werfBinPath, opts, "dismiss", "--env", "dev", "--with-namespace")
}
