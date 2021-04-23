package cmd

import (
	"net/http"
	"os"
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestMakeHttpRequest(t *testing.T) {
	tests := []struct {
		TestName        string
		ShouldFail      bool
		TestApiUrl      string
		TestLogParams   map[string]string
		TestResponseUrl string
	}{
		{
			"Logs with no parameters",
			false,
			"http://localhost:8080/logs/filter",
			map[string]string{},
			"http://localhost:8080/logs/filter",
		},
		{
			"Logs by podname",
			false,
			"http://localhost:8080/logs/filter",
			map[string]string{"Podname": "openshift-kube-scheduler"},
			"http://localhost:8080/logs/filter?podname=openshift-kube-scheduler",
		},
		{
			"Logs by given time interval",
			false,
			"http://localhost:8080/logs/filter",
			map[string]string{"Tail": "00h30m"},
			"http://localhost:8080/logs/filter?podname=openshift-kube-scheduler",
		},
		{
			"Logs with max log limit",
			false,
			"http://localhost:8080/logs/filter",
			map[string]string{"Limit": "5"},
			"http://localhost:8080/logs/filter?podname=openshift-kube-scheduler&maxlogs=5",
		},
	}

	NewCmdLogFilter(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	logParameters := LogParameters{}
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
		if res.Request.URL.String() != tt.TestResponseUrl {
			t.Errorf("Response url expected to be %s and found %s", tt.TestResponseUrl, res.Request.URL)
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be %d and found %d", http.StatusOK, res.StatusCode)
		}
	}
}
