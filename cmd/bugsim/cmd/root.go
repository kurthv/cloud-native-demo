/*
Copyright © 2022 Lutz Behnke <lutz.behnke@gmx.de>
This file is part of {{ .appName }}
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bugsim",
	Short: "Provide a web-app to allow users to simulate bugs in a K8s cluster",
	Long: `bugsim allows users of a website to terminate pods in a Kubernetes cluster, in order to
	simulate a fault in the code executed by the pod.
	
	bugsim will ensure that only a certain amount of pods in a deployment are available for
	termination, regardles of the number of concurrent users.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.scratch.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
