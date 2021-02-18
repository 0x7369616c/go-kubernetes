// +build kubeall kubernetes
// NOTE: See the notes in the other Kubernetes example tests for why this build tag is included.

package test

import (
	"fmt"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

func TestKubernetesHelloWorld(t *testing.T) {
	t.Parallel()

	// website::tag::1:: Path to the Kubernetes resource config we will test.
	kubeResourcePath := "../k8s-deployment.yml"

	// website::tag::2:: Setup the kubectl config and context.
	options := k8s.NewKubectlOptions("", "", "default")

	// website::tag::6:: At the end of the test, run "kubectl delete" to clean up any resources that were created.
	defer k8s.KubectlDelete(t, options, kubeResourcePath)

	// website::tag::3:: Run `kubectl apply` to deploy. Fail the test if there are any errors.
	k8s.KubectlApply(t, options, kubeResourcePath)

	// website::tag::4:: Verify the service is available and get the URL for it.
	k8s.WaitUntilServiceAvailable(t, options, "go-hello-world-service", 10, 1*time.Second)
	tunnel := k8s.NewTunnel(options, k8s.ResourceTypeService, "go-hello-world-service", 0, 8080)
	defer tunnel.Close()
	tunnel.ForwardPort(t)
	url := fmt.Sprintf("http://%s", tunnel.Endpoint())

	// website::tag::5:: Make an HTTP request to the URL and make sure it returns a 200 OK with the body "Hello, World".
	http_helper.HttpGetWithRetry(t, url, nil, 200, "Hello, Guest", 30, 3*time.Second)
}
