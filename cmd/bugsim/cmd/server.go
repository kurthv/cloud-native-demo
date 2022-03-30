/*
Copyright © 2022 Lutz Behnke <lutz.behnke@gmx.de>
This file is part of cloud-native demo
*/
package cmd

import (
	"fmt"
	"html/template"
	"math/rand"
	"time"

	"io"
	"log"
	"net/http"
	"os"

	mux "github.com/gorilla/mux"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"

	"github.com/cypherfox/cloud-native-demo/pkg/k8s"
)

var Port uint16
var k8sClient *k8s.K8sClient
var Namespace string
var Deployment string
var SuccessRate uint8

var root_templ *template.Template
var failed_templ *template.Template
var router *mux.Router

const root_tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Welcome to BugSim</title>
	</head>
	<body>
	    <div class=welcome-msg>Willkommen zum Bugsimulator!<p>Möchtest du Bug spielen? Du hast eine 15% Wahrscheinlichkeit, den Pod zu erschießen, den du auswählst. Klicke einfach einen der Links mit den Namen der Pods unten</div>
		<table>
		<tr><th><div>Name</div></th><th><div>Status</div></th><th><div>Alter</div></th></tr>
		{{range .Items}}
		    <tr>
			<td><div><a href="/api/delete/{{ .Name }}">{{ .Name }}</a></div></td>
			<td><div>{{ .State }}</div></td>
			<td><div>{{ .AgeString }}</div></td>
			</tr>
		{{else}}<div><strong>no pods</strong></div>{{end}}
		</table>
	</body>
</html>`

const failed_tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Welcome to BugSim</title>
	</head>
	<body>
	    <div>Dein niederträchtiger Angriff ist fehlgeschlagen!<p>Möchtest du noch einmal spielen?</div>
		<div><p>Zurück zur <a href="/">Startseite</a></p></div>
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

	serverCmd.Flags().Uint16VarP(&Port, "port", "p", 80, "Port on which to listen for web requests")
	serverCmd.Flags().StringVarP(&Namespace, "namespace", "n", "default", "Namespace in which to look for pods.")
	serverCmd.Flags().StringVarP(&Deployment, "deployment", "d", "web", "Deployment from which to delete pods.")
	serverCmd.Flags().Uint8VarP(&SuccessRate, "success-rate", "r", 15, "Success rate of the simulated bug attack on the pods in percent [1-100] (default are 15%")
}

func doServer() error {
	var err error

	if SuccessRate < 1 || SuccessRate > 100 {
		fmt.Printf("Invalid success rate percentage: %d. Must be between 1 and 100")
		return os.ErrInvalid
	}

	fmt.Println("Setting up Kubernetes Client")
	k8sClient, err = k8s.NewKubeClient()
	if err != nil {
		fmt.Printf("Initializing Kubernetes client failed: %s", err.Error())
		return err
	}

	root_templ, err = template.New("rootPage").Parse(root_tpl)
	if err != nil {
		fmt.Printf("Initializing root template failed: %s", err.Error())
		return err
	}
	failed_templ, err = template.New("failedPage").Parse(failed_tpl)
	if err != nil {
		fmt.Printf("Initializing failed attack template failed: %s", err.Error())
		return err
	}

	router = mux.NewRouter()

	fmt.Println("Setting up pages")
	router.HandleFunc("/", rootPage)
	router.HandleFunc("/api/delete/{id}", deleteSinglePod)
	http.Handle("/", router)

	fmt.Printf("Starting to serve on port %d\n", Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), nil))

	return nil
}

type podData struct {
	Name      string
	State     string
	AgeString string
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

	podDataArr := []podData{}

	for _, pod := range pods.Items {
		podDataArr = append(podDataArr, podData{
			Name:      pod.GetName(),
			State:     statusMessage(pod),
			AgeString: time.Now().Sub(pod.GetCreationTimestamp().Time).String(),
		})
		fmt.Printf("%d: Setting %s to state %s, age %s",
			len(podDataArr),
			podDataArr[len(podDataArr)-1].Name,
			podDataArr[len(podDataArr)-1].State,
			podDataArr[len(podDataArr)-1].AgeString,
		)
	}

	data := struct {
		Items []podData
	}{
		Items: podDataArr,
	}

	err = root_templ.Execute(w, data)
	if err != nil {
		fmt.Printf("generating root page from template failed: %s\n", err.Error())
		return
	}

}

func statusMessage(pod v1.Pod) string {
	if pod.DeletionTimestamp != nil {
		return "Terminating"
	}
	return string(pod.Status.Phase)
}

func deleteSinglePod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	podName := vars["id"]

	// Cast a vote
	probability := rand.Float32()
	if probability > (float32(SuccessRate) / 100.0) {
		data := struct{}{}

		err := failed_templ.Execute(w, data)
		if err != nil {
			fmt.Printf("generating root page from template failed: %s\n", err.Error())
		}

		return
	}

	// vote was successful -> kill the pod
	err := k8sClient.DeletePod(podName, Namespace)
	if err != nil {
		fmt.Printf("deleting pod %s failed: %s\n", podName, err.Error())
		return
	}

	// TODO: change this into a redirect, in order to clean up the URL.
	rootPage(w, r)
}

func respPrintf(w http.ResponseWriter, format string, a ...interface{}) error {
	_, err := io.WriteString(w, fmt.Sprintf(format, a...))
	if err != nil {
		fmt.Printf("cannot write to response: %s", err.Error())
		return err
	}
	return nil
}
