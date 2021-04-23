package cmd

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ViaQ/log-exploration-oc-plugin/pkg/k8sresources"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/fake"
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
			t.Errorf("Expected error is %v, found %v", tt.Error, err)
		}
	}
}

func TestPrintLogs(t *testing.T) {
	tests := []struct {
		TestName    string
		ShouldFail  bool
		TestLogList []string
		TestLimit   string
		Error       error
	}{
		{
			"Test correct LogList",
			false,
			[]string{"test log-1", "test log-2", "test log-3"},
			"5",
			nil,
		},
		{
			"Empty LogList",
			false,
			[]string{},
			"5",
			fmt.Errorf("no logs present, or input parameters were invalid"),
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
			fmt.Errorf("incorrect \"limit\" value entered, an integer value between 0 and 1000 is required"),
		},
	}

	for _, tt := range tests {
		t.Log("Running:", tt.TestName)
		err := printLogs(tt.TestLogList, genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}, tt.TestLimit)
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
			"Daemonset doesn't exist",
			false,
			"dummy-daemon",
			"openshift-logging",
			[]string{},
			fmt.Errorf("daemon set \"dummy-daemon\" not found in namespace \"openshift-logging\""),
		},
		{
			"Daemonset is present",
			false,
			"openshift-daemon",
			"openshift-logging",
			[]string{},
			nil,
		},
	}

	clientset := fake.NewSimpleClientset(
		&appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "openshift-daemon",
				Namespace:   "openshift-logging",
				Annotations: map[string]string{},
			},
			Spec: appsv1.DaemonSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"name": "logging"},
				},
			},
		})

	clientset.CoreV1().Pods("openshift-logging").Create(context.TODO(),
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{
			Name:        "openshift-daemon",
			Namespace:   "openshift-logging",
			Annotations: map[string]string{},
		},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "logging",
					},
				},
			},
		}, metav1.CreateOptions{})

	for _, tt := range tests {
		t.Log("Running:", tt.TestName)

		err := k8sresources.GetDaemonSetPodsList(clientset, &tt.TestLogList, tt.DaemonSet, tt.Namespace)
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

func TestGetStatefulSetPodsList(t *testing.T) {
	tests := []struct {
		TestName    string
		ShouldFail  bool
		StatefulSet   string
		Namespace   string
		TestLogList []string
		Error       error
	}{
		{
			"Statefulset doesn't exist",
			false,
			"dummy-statefulset",
			"openshift-logging",
			[]string{},
			fmt.Errorf("stateful set \"dummy-statefulset\" not found in namespace \"openshift-logging\""),
		},
		{
			"Statefulset is present",
			false,
			"openshift-stateful",
			"openshift-logging",
			[]string{},
			nil,
		},
	}

	clientset := fake.NewSimpleClientset(
		&appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "openshift-stateful",
				Namespace:   "openshift-logging",
				Annotations: map[string]string{},
			},
			Spec: appsv1.StatefulSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"name": "logging"},
				},
			},
		})

	clientset.CoreV1().Pods("openshift-logging").Create(context.TODO(),
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{
			Name:        "openshift-stateful",
			Namespace:   "openshift-logging",
			Annotations: map[string]string{},
		},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "logging",
					},
				},
			},
		}, metav1.CreateOptions{})

	for _, tt := range tests {
		t.Log("Running:", tt.TestName)

		err := k8sresources.GetStatefulSetPodsList(clientset, &tt.TestLogList, tt.StatefulSet, tt.Namespace)
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

func TestGetDeploymentPodsList(t *testing.T) {
	tests := []struct {
		TestName    string
		ShouldFail  bool
		Deployment   string
		Namespace   string
		TestLogList []string
		Error       error
	}{
		{
			"Deployment doesn't exist",
			false,
			"dummy-deployment",
			"openshift-logging",
			[]string{},
			fmt.Errorf("deployment \"dummy-deployment\" not found in namespace \"openshift-logging\""),
		},
		{
			"Deployment is present",
			false,
			"openshift-deployment",
			"openshift-logging",
			[]string{},
			nil,
		},
	}

	clientset := fake.NewSimpleClientset(
		&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "openshift-deployment",
				Namespace:   "openshift-logging",
				Annotations: map[string]string{},
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"name": "logging"},
				},
			},
		})

	clientset.CoreV1().Pods("openshift-logging").Create(context.TODO(),
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{
			Name:        "openshift-deployment",
			Namespace:   "openshift-logging",
			Annotations: map[string]string{},
		},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "logging",
					},
				},
			},
		}, metav1.CreateOptions{})

	for _, tt := range tests {
		t.Log("Running:", tt.TestName)

		err := k8sresources.GetDeploymentPodsList(clientset, &tt.TestLogList, tt.Deployment, tt.Namespace)
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
