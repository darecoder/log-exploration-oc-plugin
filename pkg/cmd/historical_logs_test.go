package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		TestName      string
		ShouldFail    bool
		TestLogParams map[string]string
		Error         error
	}{
		{
			"Logs with no parameters",
			false,
			map[string]string{},
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
		err := logParameters.Execute(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
		if err != tt.Error {
			t.Fail()
		}
	}
}

func TestPrintLogs(t *testing.T) {
	tests := []struct {
		TestName    string
		ShouldFail  bool
		TestLogList []string
		TestLimit   string
		Response    error
	}{
		{
			"Empty LogList",
			false,
			[]string{},
			"5",
			nil,
		},
		{
			"Limit equals to 0",
			false,
			[]string{"test log-1", "test log-2", "test log-3"},
			"0",
			nil,
		},
		{
			"Negative limit",
			false,
			[]string{"test log-1", "test log-2", "test log-3"},
			"-2",
			nil,
		},
	}

	for _, tt := range tests {
		err := printLogs(tt.TestLogList, genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}, tt.TestLimit)
		fmt.Printf("\nError: %v\n", err)
		t.Fail()
	}
}

func TestGetLogs(t *testing.T) {
	tests := []struct {
		TestName         string
		ShouldFail       bool
		TestLogList      []string
		TestResponseBody []byte
		Error            error
	}{
		{
			"Empty loglist",
			false,
			[]string{},
			[]byte{},
			nil,
		},
		{
			"Empty response",
			false,
			[]string{},
			[]byte{},
			nil,
		},
	}
	for _, tt := range tests {
		t.Log("Running:", tt.TestName)
		err := getLogs(tt.TestResponseBody, &tt.TestLogList)
		if tt.Error != err {
			t.Fail()
		}
	}

}

func TestGetDaemonSetPodsList(t *testing.T) {
	tests := []struct {
		TestName    string
		ShouldFail  bool
		DaemonSet   string
		Namespace   string
		TestLogList []string
		Error       error
	}{
		{
			"Empty loglist",
			false,
			"",
			"",
			[]string{},
			nil,
		},
		{
			"Empty response",
			false,
			"",
			"",
			[]string{},
			nil,
		},
	}

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Errorf("kubeconfig Error: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Errorf("an error occurred while creating a kubernetes client: %v", err)
	}

	for _, tt := range tests {
		t.Log("Running:", tt.TestName)

		err := GetDaemonSetPodsList(clientset, tt.DaemonSet, tt.Namespace, &tt.TestLogList)
		if tt.Error != err {
			t.Fail()
		}
	}
}

func TestGetStatefulSetPodsList(t *testing.T) {
	tests := []struct {
		TestName    string
		ShouldFail  bool
		StatefulSet string
		Namespace   string
		TestLogList []string
		Error       error
	}{
		{
			"Empty loglist",
			false,
			"",
			"",
			[]string{},
			nil,
		},
		{
			"Empty response",
			false,
			"",
			"",
			[]string{},
			nil,
		},
	}

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Errorf("kubeconfig Error: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Errorf("an error occurred while creating a kubernetes client: %v", err)
	}

	for _, tt := range tests {
		t.Log("Running:", tt.TestName)

		err := GetStatefulSetPodsList(clientset, tt.StatefulSet, tt.Namespace, &tt.TestLogList)
		if tt.Error != err {
			t.Fail()
		}
	}
}

func TestGetDeploymentPodsList(t *testing.T) {
	tests := []struct {
		TestName    string
		ShouldFail  bool
		Deployment  string
		Namespace   string
		TestLogList []string
		Error       error
	}{
		{
			"Empty loglist",
			false,
			"",
			"",
			[]string{},
			nil,
		},
		{
			"Empty response",
			false,
			"",
			"",
			[]string{},
			nil,
		},
	}

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Errorf("kubeconfig Error: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Errorf("an error occurred while creating a kubernetes client: %v", err)
	}

	for _, tt := range tests {
		t.Log("Running:", tt.TestName)

		err := GetDaemonSetPodsList(clientset, tt.Deployment, tt.Namespace, &tt.TestLogList)
		if tt.Error != err {
			t.Fail()
		}
	}
}
