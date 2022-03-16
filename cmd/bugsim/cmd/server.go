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

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run the web listener",
	Long:  `run the web listener`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")

		doServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().Int16VarP(&Port, "port", "p", 80, "port on which to listen for web requests")
}

func doServer() {
	var err error

	k8sClient, err = k8s.NewKubeClient()
	if err != nil {
		fmt.Printf("Initializing Kubernetes client failed: %s", err.Error())
		os.Exit(1)
	}

	http.HandleFunc("/", rootPage)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func rootPage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, fmt.Sprintf("Hello, Bug from %s \n", r.RemoteAddr))
	// io.WriteString(w, fmt.Sprintf("Current Number of Pods: %i \n", k8s.GetPods().length()))

	pods, err := k8sClient.GetPods()
	if err != nil {
		fmt.Printf("reading pods failed: %s", err.Error())
		return
	}

	for _, pod := range pods.Items {

		fmt.Printf("%s\n", pod.Name)
	}
}
