package main

import (
	"log"
	"os"

	"github.com/tesarwijaya/ouroboros/cmd"
)

func main() {
	appCmd := cmd.NewCmd()

	if err := appCmd.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
