/*
Copyright © 2022 Lutz Behnke <lutz.behnke@gmx.de>
This file is part of cloud-native demo
*/
package cmd

import (
	"fmt"

	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/cypherfox/cloud-native-demo/pkg/k8s"
)

var Port int16
var k8sClient *k8s.K8sClient
var Namespace string
var Deployment string

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

	fmt.Println("Setting up pages")
	http.HandleFunc("/", rootPage)

	fmt.Printf("Starting to serve on port %d\n", Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), nil))

	return nil
}

func rootPage(w http.ResponseWriter, r *http.Request) {
	err := respPrintf(w, "Hello, Bug from %s \n", r.RemoteAddr)
	if err != nil {
		return
	}

	pods, err := k8sClient.GetPods(Namespace)
	if err != nil {
		fmt.Printf("reading pods failed: %s\n", err.Error())
		return
	}

	err = respPrintf(w, "Current Number of Pods: %d \n", len(pods.Items))
	if err != nil {
		return
	}

	for i, pod := range pods.Items {

		err = respPrintf(w, "Pod %d: %s\n", i, pod.Name)
		if err != nil {
			return
		}

	}
}

func respPrintf(w http.ResponseWriter, format string, a ...interface{}) error {
	_, err := io.WriteString(w, fmt.Sprintf(format, a...))
	if err != nil {
		fmt.Printf("cannot write to response: %s", err.Error())
		return err
	}
	return nil
}
