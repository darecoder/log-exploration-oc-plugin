package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ViaQ/log-exploration-oc-plugin/pkg/client"
	"github.com/ViaQ/log-exploration-oc-plugin/pkg/k8sresources"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/fake"
)

type HTTPServer struct {
	server *http.Server
	router *http.ServeMux
}

func NewHTTPServer(host string, port uint) *HTTPServer {
	s := &HTTPServer{
		router: http.NewServeMux(),
	}
	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: s,
	}
	return s
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("http request %s: %s\r\n", r.Method, r.URL.Path)
	s.router.ServeHTTP(w, r)
}

func getLogsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is serving from http-HandleFunc")
}

func TestExecute(t *testing.T) {
	tests := []struct {
		TestName      string
		ShouldFail    bool
		TestLogParams map[string]string
		TestResources map[string]string
		Arguments     []string
		Error         error
	}{
		{
			"Logs with no parameters",
			false,
			map[string]string{},
			map[string]string{"Deployment": "openshift-deployment"},
			[]string{"deployment=openshift-deployment"},
			nil,
		},
		// {
		// 	"Logs with podname & namespace",
		// 	false,
		// 	map[string]string{"Podname": "openshift-logging-1234", "Namespace": "openshift-logging"},
		// 	map[string]string{"Deployment": "openshift-deployment"},
		// 	[]string{"deployment=openshift-deployment"},
		// 	nil,
		// },
		// {
		// 	"Logs with tail parameter",
		// 	false,
		// 	map[string]string{"Tail": "30m"},
		// 	map[string]string{"Deployment": "openshift-deployment"},
		// 	[]string{"deployment=openshift-deployment"},
		// 	nil,
		// },
		// {
		// 	"Logs with multiple parameters",
		// 	false,
		// 	map[string]string{
		// 		"Podname":   "openshift-logging-1234",
		// 		"Namespace": "openshift-logging",
		// 		"Tail":      "30m",
		// 		"Limit":     "5",
		// 	},
		// 	map[string]string{"Deployment": "openshift-deployment"},
		// 	[]string{"deployment=openshift-deployment"},
		// 	nil,
		// },
		// {
		// 	"Logs with valid integer limit",
		// 	false,
		// 	map[string]string{"Limit": "5"},
		// 	map[string]string{"Deployment": "openshift-deployment"},
		// 	[]string{"deployment=openshift-deployment"},
		// 	nil,
		// },
		// {
		// 	"Logs with negative limit",
		// 	false,
		// 	map[string]string{"Limit": "-5"},
		// 	map[string]string{"Deployment": "openshift-deployment"},
		// 	[]string{"deployment=openshift-deployment"},
		// 	fmt.Errorf("incorrect \"limit\" value entered, an integer value between 0 and 1000 is required"),
		// },
	}

	for _, tt := range tests {
		t.Log("Running:", tt.TestName)
		logParameters := LogParameters{}
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
		for k, v := range tt.TestResources {
			switch k {
			case "Deployment":
				logParameters.Resources.IsDeployment = true
			case "Daemonset":
				logParameters.Resources.IsDaemonSet = true
			case "Statefulset":
				logParameters.Resources.IsStatefulSet = true
			case "Pod":
				logParameters.Resources.IsPod = true
			}
			logParameters.Resources.Name = v
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
				Labels:      map[string]string{"name": "logging"},
			},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "logging",
						},
					},
				},
			}, metav1.CreateOptions{})

		kubernetesOptions := &client.KubernetesOptions{
			Clientset:        clientset,
			ClusterUrl:       "loclahost.com:8080",
			CurrentNamespace: "openshift-logging",
		}

		server := NewHTTPServer("http://log-exploration-api-route-openshift-logging.apps.com/logs", 0)
		th := http.HandlerFunc(getLogsHandler)
		server.router.Handle("/filter", th)
		// server.RegisterHandler("http://log-exploration-api-route-openshift-logging.apps.com/logs/filter", getLogsHandler)

		req, Error := http.NewRequest(http.MethodGet, "http://log-exploration-api-route-openshift-logging.apps.com/logs/filter", nil)
		if Error != nil {
			t.Errorf("Failed to create HTTP request. E: %v", Error)
		}
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)
		resp := rr.Body.String()
		t.Errorf("Response: %v", resp)

		
		err := logParameters.Execute(kubernetesOptions, genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}, tt.Arguments)
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
		{
			"Invalid limit",
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

// func getLogs(gctx *gin.Context) {
// 	var params LogParameters
// 	err := gctx.Bind(&params)
// 	if err != nil {
// 		gctx.JSON(http.StatusInternalServerError, gin.H{ //If error is not nil, an internal server error might have ocurred
// 			"An error occurred": []string{err.Error()},
// 		})
// 	}

// 	gctx.JSON(http.StatusOK, gin.H{
// 		"Logs": []string{"test-log-1", "test-log-2", "test-log-3", "test-log-4", "test-log-5"},
// 	})
// }
