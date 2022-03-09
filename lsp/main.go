package main

import (
	"fmt"
	"os"

	"github.com/glennsarti/blog-wasm-vscode/lsp/internal/cmd"
)

func main() {
	exitStatus := -1
	var err error

	progArgs := make([]string, 0)
	if len(os.Args) > 1 {
		progArgs = os.Args[1:]
	}
	if len(progArgs) == 0 {
		runner := cmd.ServeCommandRunner{Version: VersionString()}
		exitStatus, err = runner.Run(progArgs)
		if err != nil {
			fmt.Println("Error ", err)
		}
	} else {
		switch progArgs[0] {
		case "serve":
			runner := cmd.ServeCommandRunner{Version: VersionString()}
			exitStatus, err = runner.Run(progArgs[1:])
			if err != nil {
				fmt.Println("Error ", err)
			}
		}
	}
	os.Exit(exitStatus)
}
