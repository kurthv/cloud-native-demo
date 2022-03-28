/*
Copyright © 2022 Lutz Behnke <lutz.behnke@gmx.de>
This file is part of cloud-native demo
*/
package cmd

import (
	"fmt"
	"html/template"

	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"

	"github.com/cypherfox/cloud-native-demo/pkg/k8s"
)

var Port int16
var k8sClient *k8s.K8sClient
var Namespace string
var Deployment string

var templ *template.Template

const root_tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Welcome to BugSim</title>
	</head>
	<body>
		{{range .Items}}<div>Pod: <a href="/api/delete/{{ .Name }}">{{ .Name }}</a></div>{{else}}<div><strong>no pods</strong></div>{{end}}
	</body>
</html>`

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run the web listener",
	Long:  `run the web listener`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("server called")

		err := doServer()
		if err != nil {
			return err
		}
		fmt.Println("server ended")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().Int16VarP(&Port, "port", "p", 80, "port on which to listen for web requests")
	serverCmd.Flags().StringVarP(&Namespace, "namespace", "n", "default", "Namespace in which to look for pods.")
	serverCmd.Flags().StringVarP(&Deployment, "deployment", "d", "web", "deployment from which to delete pods.")
}

func doServer() error {
	var err error

	fmt.Println("Setting up Kubernetes Client")
	k8sClient, err = k8s.NewKubeClient()
	if err != nil {
		fmt.Printf("Initializing Kubernetes client failed: %s", err.Error())
		os.Exit(1)
	}

	templ, err = template.New("rootPage").Parse(root_tpl)
	if err != nil {
		fmt.Printf("Initializing root template failed: %s", err.Error())
		return err
	}

	fmt.Println("Setting up pages")
	http.HandleFunc("/api/delete/{id}", deleteSinglePod)
	http.HandleFunc("/", rootPage2)

	fmt.Printf("Starting to serve on port %d\n", Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), nil))

	return nil
}

func rootPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := respPrintf(w, "Hello, Bug from %s \n", r.RemoteAddr)
	if err != nil {
		return
	}

	pods, err := k8sClient.GetPods(Namespace, Deployment)
	if err != nil {
		fmt.Printf("reading pods failed: %s\n", err.Error())
		return
	}

	err = respPrintf(w, "Current Number of Pods: %d \n", len(pods.Items))
	if err != nil {
		return
	}

	for i, pod := range pods.Items {

		err = printLink(w, i, pod.Name)
		if err != nil {
			return
		}

	}
}

func rootPage2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := respPrintf(w, "Hello, Bug from %s \n", r.RemoteAddr)
	if err != nil {
		return
	}

	pods, err := k8sClient.GetPods(Namespace, Deployment)
	if err != nil {
		fmt.Printf("reading pods failed: %s\n", err.Error())
		return
	}

	data := struct {
		Items []v1.Pod
	}{
		Items: pods.Items,
	}

	err = templ.Execute(w, data)
	if err != nil {
		fmt.Printf("generating root page from template failed: %s\n", err.Error())
		return
	}

}

func printLink(w http.ResponseWriter, num int, podName string) error {
	respPrintf(w, "Pod %d: <a href=\"/api/delete/%s\">%s</a>\n", num, podName, podName)
	return nil
}

func deleteSinglePod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Printf("Key: %s\n", key)
}

func respPrintf(w http.ResponseWriter, format string, a ...interface{}) error {
	_, err := io.WriteString(w, fmt.Sprintf(format, a...))
	if err != nil {
		fmt.Printf("cannot write to response: %s", err.Error())
		return err
	}
	return nil
}
