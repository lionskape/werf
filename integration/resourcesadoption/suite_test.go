// +build integration integration_k8s

package resourcesadoption

import (
	"fmt"

	"github.com/flant/kubedog/pkg/kube"
	"github.com/flant/werf/integration/utils/werfexec"
	"github.com/ghodss/yaml"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func newDeployment(manifestYaml string) *v1.Deployment {
	obj := &v1.Deployment{}
	Expect(yaml.Unmarshal([]byte(manifestYaml), &obj)).To(Succeed())
	return obj
}

func newNamespace(manifestYaml string) *corev1.Namespace {
	obj := &corev1.Namespace{}
	Expect(yaml.Unmarshal([]byte(manifestYaml), &obj)).To(Succeed())
	return obj
}

var _ = Describe("Werf release installer and updater", func() {
	JustBeforeEach(func() {
		Expect(kube.Init(kube.InitOptions{})).To(Succeed())
	})

	Context("when installing a new release with resources that already exist in cluster", func() {
		namespace := "resourcesadoption-app1-dev"

		It("should fail to install release", func(done Done) {
			d := newDeployment(`
kind: Deployment
apiVersion: apps/v1
metadata:
  name: mydeploy2
  labels:
    service: mydeploy2
spec:
  replicas: 1
  selector:
    matchLabels:
      service: mydeploy2
  template:
    metadata:
      labels:
        service: mydeploy2
    spec:
      containers:
      - name: main 
        command: [ "/bin/bash", "-c", "while true; do date; sleep 1; done" ]
        image: ubuntu:18.04
`)

			Expect(kube.Kubernetes.AppsV1().Deployments(namespace).Create(d)).To(Succeed())

			ns := &corev1.Namespace{}

			Expect(yaml.Unmarshal([]byte(fmt.Sprintf(`
apiVersion: v1
kind: Namespace
metadata:
  name: %s
`, namespace)), &ns)).To(Succeed())

			kube.Kubernetes.CoreV1().Namespaces().Create(ns)
			Expect(kube.Kubernetes.CoreV1().Namespaces().Create(ns)).To(Succeed())

			//if mydeploy1, err := kube.Kubernetes.AppsV1().Deployments(namespace).Get("mydeploy1", metav1.GetOptions{}); err == nil {
			//}

			Expect(werfDeploy("app1", werfexec.CommandOptions{})).To(Succeed())

			close(done)
		})

		It("should not delete already existing resources on failed release removal", func(done Done) {
			close(done)
		})

		It("should delete new resources created during failed release installation on failed release removal", func(done Done) {
			close(done)
		})
	})
})

func werfDeploy(dir string, opts werfexec.CommandOptions) error {
	return werfexec.ExecWerfCommand(dir, werfBinPath, opts, "deploy", "--env", "dev")
}

func werfDismiss(dir string, opts werfexec.CommandOptions) error {
	return werfexec.ExecWerfCommand(dir, werfBinPath, opts, "dismiss", "--env", "dev", "--with-namespace")
}
