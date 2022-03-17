/*
Copyright © 2022 Lutz Behnke <lutz.behnke@gmx.de>
This file is part of the cloud-native demo
*/
package main

import (
	"fmt"
	"os"

	"github.com/cypherfox/cloud-native-demo/cmd/bugsim/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("Command failed: %s", err)
		os.Exit(1)
	}
}
