package teste2e

import (
	"fmt"
	"os"
	"testing"

	logs "github.com/ViaQ/log-exploration-oc-plugin/pkg/cmd/historical_logs"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestMakeHttpRequest(t *testing.T) {
	tests := []struct {
		TestName      string
		ShouldFail    bool
		TestApiUrl    string
		TestLogParams map[string]string
		TestURL       string
		Response      map[string][]string
	}{
		{
			"Logs with no parameters",
			false,
			"log-exploration-api-route-openshift-logging.apps.test.devcluster.openshift.com",
			map[string]string{},
			"",
			map[string][]string{"Logs": {"test-log-1", "test-log-2", "test-log-3"}},
		},
		{
			"Logs by podname",
			false,
			"log-exploration-api-route-openshift-logging.apps.test.devcluster.openshift.com",
			map[string]string{"Podname": "openshift-kube-scheduler"},
			"",
			map[string][]string{"Logs": {"test-log-1", "test-log-2", "test-log-3"}},
		},
		{
			"Logs by given time interval",
			false,
			"log-exploration-api-route-openshift-logging.apps.test.devcluster.openshift.com",
			map[string]string{"Tail": "00h30m"},
			"",
			map[string][]string{"Logs": {"test-log-1", "test-log-2", "test-log-3"}},
		},
		{
			"Logs with max log limit",
			false,
			"log-exploration-api-route-openshift-logging.apps.emishra-test121.devcluster.openshift.com",
			map[string]string{"Limit": "5"},
			"",
			map[string][]string{"Logs": {"test-log-1", "test-log-2", "test-log-3"}},
		},
	}

	logs.NewCmdLogFilter(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	logParameters := logs.LogParameters{}
	for _, tt := range tests {
		t.Log("Running:", tt.TestName)
		for k, v := range tt.TestLogParams {
			switch k {
			case "Namespace":
				logParameters.Namespace = v
			case "Podname":
				logParameters.Podname = v
			case "Tail":
				logParameters.Tail = v
			case "StartTime":
				logParameters.StartTime = v
			case "EndTime":
				logParameters.EndTime = v
			case "Level":
				logParameters.Level = v
			case "Limit":
				logParameters.Limit = v
			case "Deployment":
				logParameters.Deployment = v
			case "StatefulSet":
				logParameters.StatefulSet = v
			case "DaemonSet":
				logParameters.DaemonSet = v
			}
		}
		res, _ := logParameters.makeHttpRequest(tt.TestApiUrl)
		fmt.Print(res)
	}
}
