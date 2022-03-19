/*
Copyright © 2022 Lutz Behnke <lutz.behnke@gmx.de>
This file is part of the cloud-native demo
*/
package main

import (
	"fmt"
	"os"

	"github.com/cypherfox/cloud-native-demo/cmd/bugsim/cmd"
	"github.com/cypherfox/cloud-native-demo/pkg/version"
)

func main() {
	fmt.Printf("This is bugsim %s build on %s \n", version.BuildVersion, version.BuiltTime)

	err := cmd.Execute()
	if err != nil {
		fmt.Printf("Command failed: %s", err)
		os.Exit(1)
	}
}
