package cmd

import (
	"testing"

	"github.com/ViaQ/log-exploration-oc-plugin/pkg/client"
	"github.com/ViaQ/log-exploration-oc-plugin/pkg/k8sresources"
)

func TestProcessLogParameters(t *testing.T) {
	tests := []struct {
		TestName      string
		ShouldFail    bool
		TestLogParams map[string]string
		Arguments     []string
		Error         error
	}{
		{
			"Logs with no parameters",
			false,
			map[string]string{},
			[]string{},
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
			}
		}
		logParameters.Resources = k8sresources.Resources{}
		// clientset := fake.NewSimpleClientset()
		// Have to update this later
		kubernetesOptions, _ := client.KubernetesClient()
		logParameters.ProcessLogParameters(kubernetesOptions, tt.Arguments)
		// if err != tt.Error {
			// t.Errorf("Expected error is %v, found %v", tt.Error, err)
		// }
	}
}
