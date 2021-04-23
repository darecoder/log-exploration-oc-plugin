package cmd

import (
	"testing"
)

func TestFetchLogs(t *testing.T) {
	tests := []struct {
		TestName      string
		ShouldFail    bool
		TestLogList   []string
		TestApiUrl    string
		TestLogParams map[string]string
		Error         error
	}{
		{
			"Logs with no parameters",
			false,
			[]string{},
			"http://localhost:8080/logs/filter",
			map[string]string{},
			nil,
		},
		{
			"Logs by podname",
			false,
			[]string{},
			"http://localhost:8080/logs/filter",
			map[string]string{"Podname": "openshift-kube-scheduler"},
			nil,
		},
		{
			"Logs by given time interval",
			false,
			[]string{},
			"http://localhost:8080/logs/filter",
			map[string]string{"Tail": "00h30m"},
			nil,
		},
		{
			"Logs with max log limit",
			false,
			[]string{},
			"http://localhost:8080/logs/filter",
			map[string]string{"Limit": "5"},
			nil,
		},
	}

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
		err := fetchLogs(&tt.TestLogList, tt.TestApiUrl, &logParameters)
		if err == nil && tt.Error != nil {
			t.Errorf("Expected error is %v, found %v", tt.Error, err)
		}
		if err != nil && tt.Error == nil {
			t.Errorf("Expected error is %v, found %v", tt.Error, err)
		}
		if err != nil && tt.Error != nil && err.Error() != tt.Error.Error() {
			t.Errorf("Expected error is %v, found %v", tt.Error, err)
		}
	}
}
